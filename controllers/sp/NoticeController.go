package sp

import (
	"encoding/json"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreBaseLib/splib/common"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"strconv"
)

// 接收订阅退订续订通知
type NotificationController struct {
	BaseController
}

type PostParseForm struct {
	Command         string `form:"command"`
	SessionIdID     string `form:"sessionId"`
	StatusNumber    string `form:"statusNumber"`
	StatusText      string `form:"statusText"`
	TransactionID   string `form:"trid"`
	TransactionTime string `form:"time"`
	Amount          string `form:"amount"`
	NotificationId  string `form:"notificationId"`
	SubscriptionID  string `form:"subscriptionId"`
	Msisdn          string `form:"msisdn"`
	ServiceCode     string `form:"serviceCode"`
}

type callUrl struct {
	Msisdn          string `xml:"msisdn"`
	PackageId       int    `xml:"package_id"`
	TransactionType string `xml:"TransactionType"`
	Amount          string `xml:"Amount"`
	Channel         string `xml:"channel"`
	TransactionId1  string `xml:"transaction_id1"`
	TransactionId2  string `xml:"transaction_id2"`
	TransactionId   string `xml:"transaction_id"`
	Keyword         string `xml:"keyword"`
}

type CallBackResult struct {
	CallUrl callUrl `xml:"Call_url"`
}

// 订阅通知  sessionId:[AT030159x01x1559807672429] statusNumber:[2] statusText:[Payment authorized] notificationId:[2217664102] command:[deliverSessionState]]
// 扣费通知  command:[recurrentPayment] statusText:[Charged] time:[2019-06-06 09:54:58] amount:[500] trid:[668286775] statusNumber:[2] subscriptionId:[AT030159x01x1559807672429] msisdn:[00436643607604] serviceCode:[AT030159]]

// Post 接收订阅退订续订通知
func (c *NotificationController) Post() {

	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	// var dnJson sp.DnJson
	var callBackResult CallBackResult

	// err := json.Unmarshal(body, &dnJson)
	err := json.Unmarshal(body, &callBackResult)

	fmt.Println("MW UAE Data: ", string(body))

	if err != nil {
		logs.Error("notification ERROR,接收通知错误", err.Error())
		c.Ctx.WriteString("ERROR")
		return
	}

	logs.Info("notification 通知:", callBackResult)
	var reqFormData sp.ChargeNotification

	if callBackResult != (CallBackResult{}) {

		reqFormData.SubType = callBackResult.CallUrl.TransactionType
		reqFormData.Msisdn = callBackResult.CallUrl.Msisdn
		reqFormData.ChannelId = callBackResult.CallUrl.Channel
		reqFormData.Charge = callBackResult.CallUrl.Amount
		reqFormData.Keyword = callBackResult.CallUrl.Keyword

		if reqFormData.SubType == "SUB" || reqFormData.SubType == "UNSUBG" {
			reqFormData.RequestId = callBackResult.CallUrl.TransactionId2
			reqFormData.TransactionId = callBackResult.CallUrl.TransactionId1
		} else if reqFormData.SubType == "UNSUB" {

		}
	}

	track := new(sp.AffTrack)

	var serverConfig sp.ServiceInfo
	// 接收通知 订阅成功
	// if reqFormData.SubType == "SUBSCRIBE" && reqFormData.Status == "DELIVERED" { // 订阅、退订通知
	// reqFormData.SubscriptionID = reqFormData.SessionID
	// if reqFormData.StatusNumber == "2" {
	//	// 订阅通知 在用户信息表里通过订阅ID 查询 trackID
	//	userHistory := new(sp.UserReqHistory)
	//	trackID := userHistory.GetTrackIDBySessionID(reqFormData.SubscriptionID)
	// 通过trackID 查询 点击数据
	fmt.Println("trackID1: ", reqFormData.TransactionId)
	trackID, _ := strconv.Atoi(reqFormData.TransactionId)

	if trackID != 0 {
		_ = track.GetOneByMsisdn(reqFormData.Msisdn)

		fmt.Println("handle track id: ", track.TrackID)

		// 如果没有msisdn，则到mo表进行查询
		if track.TrackID == 0 {
			// 此处应该查询mo表，关于此msisdn 的第一次订阅时间，如果是 6.1 之前的就回传到老系统，否则回传到新系统
			// 直接获取一个6.1号之前的track记录
			checkMo := new(mo.Mo)
			if _, err = checkMo.GetMoByMsisdnAndSubType(reqFormData.Msisdn); err != nil {
				fmt.Println("根据msisdn获取mo数据失败")
			} else {
				track.TrackID = checkMo.TrackID
				if err = track.GetOne(tracking.ByTrackID); err != nil {
					fmt.Println("根据msisdn trackid获取track数据失败")
				}
			}
		}

		serverConfig, err = c.getServiceConfigNotification(track.ServiceID)

		// 如果serverConfig获取出错，则进行数据存储
		// 后面可以为了同步后台数据，随机筛选一个配置进行传输
		if err != nil {
			reqFormData.Insert()
			c.Ctx.WriteString("ok")
			c.StopRun()
		}
		// sp.SendMt(serverConfig, reqFormData)
	} else {
		// 通过电话号码 msisdn 进行获取 track_id

	}

	fmt.Println("config: ", serverConfig)

	// c.Ctx.WriteString("ok")
	// return

	// }

	notify := new(notification.Notification)

	notify.SubscriptionId = reqFormData.TransactionId
	notify.TransactionID = reqFormData.RequestId
	notify.ServiceID = track.ServiceID

	// serviceConfig, _ := c.serviceCofig(notify.ServiceID)

	// 先先根据subID 查询mo数据
	moT := new(mo.Mo)
	_, err = moT.GetMoBySubscriptionID(notify.SubscriptionId)

	// 新订阅通知 ，没有找到此订阅信息，需要重新插入mo数据
	notificationType := ""
	if reqFormData.SubType == "SUB" {

		var moBase = common.MoBase{}
		moBase.SubscriptionID = notify.SubscriptionId
		moBase.Operator = serverConfig.OperatorId
		moBase.Price = fmt.Sprintf("%v", track.PostbackPrice)
		moBase.Msisdn = reqFormData.Msisdn
		moBase.Track = track.Track
		moBase.OfferID = track.OfferID
		moBase.TrackID = track.TrackID
		moBase.ServiceID = track.ServiceID
		moBase.PromoterID = track.PromoterID
		// 如果是订阅通知
		postbackStatus := true
		// if reqFormData.SubType == "SUBSCRIBE" {
		// 	postbackStatus = false
		// }

		fmt.Println("准备进入发送短信")
		fmt.Println("ApiSecret: ", serverConfig.ApiSecret)
		fmt.Println("OperatorId: ", serverConfig.OperatorId)
		sp.SendMt(serverConfig, &reqFormData)
		fmt.Println("发送短信完成")

		// 检查subID是否已经存在
		if err == nil && moT.ID != 0 { // 订阅ID 已经存在，重复通知
			logs.Info("订阅已经存在，不能新存入MO信息: ", notify.SubscriptionId)
			c.StringResult("OK")
		}

		// 存入MO数据
		moT, notificationType = splib.InsertMO(moBase, false, postbackStatus, serverConfig.ProductName)

		// 订阅成功后注册服务
		go sp.AddUser(serverConfig.UrlPost+serverConfig.ReqUrl, moT.Msisdn, moT.SubscriptionID)

	}

	// 扣费，退订通知
	if reqFormData.SubType == "REN" { // 成功扣费通知
		notificationType, _ = moT.AddSuccessMTNum(notify.SubscriptionId, notify.TransactionID)
		// sp.SendMt(serverConfig, &reqFormData)
		// } else if reqFormData.SubType == "RENEWAL" && reqFormData.Status != "Failed" { // 失败扣费通知
		// notificationType, _ = moT.AddFailedMTNum(notify.SubscriptionId, notify.TransactionID)
	} else if reqFormData.SubType == "UNSUB" || reqFormData.SubType == "UNSUBG" { // 退订通知
		notificationType, _ = moT.UnsubUpdateMo(notify.SubscriptionId)
	}

	fmt.Println("notificationType: ", notificationType)
	fmt.Println("发送订阅通知给后台")
	if notificationType != "" {
		fmt.Println("开始发送订阅通知给后台")
		notify.NotificationType = notificationType
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.PostbackPrice = track.PostbackPrice

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID
		sendNoti.ServiceID = moT.ServiceID
		sendNoti.ClickID = moT.ClickID
		sendNoti.Msisdn = moT.Msisdn
		sendNoti.CampID = track.CampID
		sendNoti.PubID = moT.PubID
		sendNoti.PostbackStatus = moT.PostbackStatus
		sendNoti.PostbackMessage = moT.PostbackMessage
		sendNoti.TransactionID = notify.TransactionID
		sendNoti.PromoterID = moT.PromoterID

		sendNoti.Keyword = moT.Keyword
		sendNoti.ShortCode = moT.ShortCode

		sendNoti.AffName = moT.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}
		sendNoti.Operator = moT.Operator

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		fmt.Println("service id ====== ", sendNoti.ServiceID)
		if sendNoti.ServiceID == "BB-NEW-ET" {
			sendNoti.SendData(admindata.SEC)
		} else {
			sendNoti.SendData(admindata.PROD)
		}
		fmt.Println("发送订阅通知给后台 完成")
	}

	reqFormData.Insert()

	c.Ctx.WriteString("ok")
}
