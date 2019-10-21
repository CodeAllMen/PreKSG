package sp

import (
	"encoding/json"
	"github.com/MobileCPX/PreBaseLib/splib"
	"github.com/MobileCPX/PreBaseLib/splib/admindata"
	"github.com/MobileCPX/PreBaseLib/splib/common"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/notification"
	"github.com/MobileCPX/PreBaseLib/splib/servicelib"
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

// 订阅通知  sessionId:[AT030159x01x1559807672429] statusNumber:[2] statusText:[Payment authorized] notificationId:[2217664102] command:[deliverSessionState]]
// 扣费通知  command:[recurrentPayment] statusText:[Charged] time:[2019-06-06 09:54:58] amount:[500] trid:[668286775] statusNumber:[2] subscriptionId:[AT030159x01x1559807672429] msisdn:[00436643607604] serviceCode:[AT030159]]

//Post 接收订阅退订续订通知
func (c *NotificationController) Get() {

	body, _ := ioutil.ReadAll(c.Ctx.Request.Body)
	var dnJson sp.DnJson

	err := json.Unmarshal(body, &dnJson)
	if err != nil {
		logs.Error("notification ERROR,接收通知错误", err.Error())
		c.Ctx.WriteString("ERROR")
		return
	}

	logs.Info("notification 通知:", dnJson)
	var reqFormData sp.ChargeNotification
	if dnJson != (sp.DnJson{}) {
		reqFormData.RequestId = dnJson.RequestId
		if dnJson.Transaction != (sp.Transaction{}) {
			reqFormData.TransactionId = dnJson.Transaction.TransactionId
			if dnJson.Transaction.Data != (sp.Data{}) {
				reqFormData.Shortcode = dnJson.Transaction.Data.Shortcode
				reqFormData.ChannelId = dnJson.Transaction.Data.ChannelId
				reqFormData.ApplicationId = dnJson.Transaction.Data.ApplicationId
				reqFormData.Country = dnJson.Transaction.Data.CountryId
				reqFormData.OperatorId = dnJson.Transaction.Data.OperatorId
				reqFormData.Msisdn = dnJson.Transaction.Data.Msisdn
				reqFormData.ActivityTime = dnJson.Transaction.Data.ActivityTime
				reqFormData.SubscriptionEnd = dnJson.Transaction.Data.SubscriptionEnd
				if dnJson.Transaction.Data.Action != (sp.Action{}) {
					reqFormData.Type = dnJson.Transaction.Data.Action.Type
					reqFormData.SubType = dnJson.Transaction.Data.Action.SubType
					reqFormData.Status = dnJson.Transaction.Data.Action.Status
					reqFormData.Rate = dnJson.Transaction.Data.Action.Rate
				}
			}

		}
	}

	track := new(sp.AffTrack)

	var serverConfig sp.ServiceInfo
	//接收通知 订阅成功
	//if reqFormData.SubType == "SUBSCRIBE" && reqFormData.Status == "DELIVERED" { // 订阅、退订通知
	//reqFormData.SubscriptionID = reqFormData.SessionID
	//if reqFormData.StatusNumber == "2" {
	//	// 订阅通知 在用户信息表里通过订阅ID 查询 trackID
	//	userHistory := new(sp.UserReqHistory)
	//	trackID := userHistory.GetTrackIDBySessionID(reqFormData.SubscriptionID)
	// 通过trackID 查询 点击数据
	trackID, _ := strconv.Atoi(reqFormData.TransactionId)
	if trackID != 0 {
		track.TrackID, _ = strconv.ParseInt(reqFormData.TransactionId, 10, 64)
		_ = track.GetOne(tracking.ByTrackID)
		serverConfig = c.getServiceConfig(track.ServiceID)
		//sp.SendMt(serverConfig, reqFormData)
	}

	//}

	notify := new(notification.Notification)

	notify.SubscriptionId = reqFormData.TransactionId
	notify.TransactionID = reqFormData.RequestId
	notify.ServiceID = serverConfig.ServiceID

	//serviceConfig, _ := c.serviceCofig(notify.ServiceID)

	// 先先根据subID 查询mo数据
	moT := new(mo.Mo)
	_, err = moT.GetMoBySubscriptionID(notify.SubscriptionId)

	// 新订阅通知 ，没有找到此订阅信息，需要重新插入mo数据
	notificationType := ""
	if reqFormData.SubType == "SUBSCRIBE" && reqFormData.Status == "Delivered" {

		// 检查subID是否已经存在
		if err == nil && moT.ID != 0 { // 订阅ID 已经存在，重复通知
			logs.Info("订阅已经存在，不能新存入MO信息: ", notify.SubscriptionId)
			c.StringResult("OK")
		}

		var moBase = common.MoBase{}
		moBase.SubscriptionID = notify.SubscriptionId
		moBase.Operator = serverConfig.OperatorId
		moBase.Price = serverConfig.Price
		moBase.Msisdn = reqFormData.Msisdn
		moBase.Track = track.Track
		// 如果是订阅通知
		postbackStatus := true
		if reqFormData.SubType == "SUBSCRIBE" {
			postbackStatus = false
		}

		// 存入MO数据
		moT, notificationType = splib.InsertMO(moBase, false, postbackStatus, serverConfig.ProductName)

		// 订阅成功后注册服务
		go servicelib.AddOrDeleteUserService(serverConfig.UrlPost, moT.Msisdn, moT.SubscriptionID)
		sp.SendMt(serverConfig, &reqFormData)

	}

	// 扣费，退订通知
	if reqFormData.SubType == "Renewal" && reqFormData.Status == "Delivered" { // 成功扣费通知
		notificationType, _ = moT.AddSuccessMTNum(notify.SubscriptionId, notify.TransactionID)
		sp.SendMt(serverConfig, &reqFormData)
	} else if reqFormData.SubType == "Renewal" && reqFormData.Status != "Failed" { // 失败扣费通知
		notificationType, _ = moT.AddFailedMTNum(notify.SubscriptionId, notify.TransactionID)
	} else if reqFormData.SubType == "Unsubscribe" && reqFormData.Status == "Delivered" { // 退订通知
		notificationType, _ = moT.UnsubUpdateMo(notify.SubscriptionId)
	}

	if notificationType != "" {
		notify.NotificationType = notificationType
		notify.Insert()

		nowTime, _ := util.GetNowTime()

		sendNoti := new(admindata.Notification)

		sendNoti.PostbackPrice = moT.PostbackPrice

		sendNoti.OfferID = moT.OfferID
		sendNoti.SubscriptionID = moT.SubscriptionID
		sendNoti.ServiceID = moT.ServiceID
		sendNoti.ClickID = moT.ClickID
		sendNoti.Msisdn = moT.Msisdn
		sendNoti.CampID = sp.ServiceData[moT.ServiceID].CampID
		sendNoti.PubID = moT.PubID
		sendNoti.PostbackStatus = moT.PostbackStatus
		sendNoti.PostbackMessage = moT.PostbackMessage
		sendNoti.TransactionID = notify.TransactionID
		sendNoti.AffName = moT.AffName
		if sendNoti.AffName == "" {
			sendNoti.AffName = "未知"
		}
		sendNoti.Operator = moT.Operator

		sendNoti.Sendtime = nowTime
		sendNoti.NotificationType = notificationType
		sendNoti.SendData(admindata.PROD)
	}

	reqFormData.Insert()

	c.Ctx.WriteString("ok")
}
