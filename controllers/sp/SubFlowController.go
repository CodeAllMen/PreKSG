package sp

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreKSG/libs"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/MobileCPX/PreKSG/service"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
)

type SubFlowController struct {
	BaseController

	// 追踪点击数据
	trackClickData *sp.AffTrack

	// serviceConf
	serviceConf sp.ServiceInfo
}

func (c *SubFlowController) Prepare() {
	// 获取track 数据
	c.trackClickData = c.getTrackData()

	// 配置信息
	c.serviceConf = c.getServiceConfig(c.trackClickData.Track.ServiceID)
}

func (c *SubFlowController) SubReq() {
	logs.Info("SubReq: ", c.Ctx.Input.URI())
	var (
		err error
		res string
	)
	track := new(sp.AffTrack)
	track.TrackID, err = strconv.ParseInt(c.Ctx.Input.Param(":trackID"), 10, 64)
	err = track.GetOne(tracking.ByTrackID)
	// 获取AOC连接
	if err != nil {
		c.RedirectURL(c.Ctx.Input.URI() + "/404")
		return
	}
	serviceConf := c.getServiceConfig(track.ServiceID)

	res = service.SubService(serviceConf, track)

	// c.Data["json"] = map[string]string{
	//	"data": res,
	// }
	if res == "" {
		c.RedirectURL(c.Ctx.Input.URI() + "/404")
	} else {
		c.RedirectURL(res)
	}
	fmt.Println(res)
}

// 第二次修改 流程，现在的流程是  通过短信进行 订阅
func (c *SubFlowController) SubReqSMS() {
	logs.Info("SubReqSMS: ", c.Ctx.Input.URI())
	var err error
	track := new(sp.AffTrack)
	track.TrackID, err = strconv.ParseInt(c.Ctx.Input.Param(":trackID"), 10, 64)
	phoneNumber := c.Ctx.Input.Param(":phoneNumber")
	err = track.GetOne(tracking.ByTrackID)
	// 获取AOC连接
	if err != nil {
		fmt.Println(err)
		c.RedirectURL(c.Ctx.Input.URI() + "/404")
		return
	}

	serviceConf := c.getServiceConfig(track.ServiceID)
	res := service.SubServiceSMS(serviceConf, track, phoneNumber)

	// 这里如果处理出错，也就是请求pin码出错，直接跳到某个页面，比如404或谷歌，否则 显示 验证pin码页面
	if res != "" {
		fmt.Println(res)
		c.RedirectURL(c.Ctx.Input.URI() + "/404")
		return
	}

	if serviceConf.KeyWord == "MA" {
		c.Data["KeyWord"] = "MYA"
	} else {
		c.Data["KeyWord"] = serviceConf.KeyWord
	}

	c.Data["ProductName"] = serviceConf.ProductName

	c.Data["URL"] = "/api/validate_sms/" + c.Ctx.Input.Param(":trackID") + "/" + phoneNumber
	c.TplName = "uae/pin.html"
}

func (c *SubFlowController) ValidateSMS() {
	type returnData struct {
		Code int    `json:"code"`
		Err  string `json:"err"`
		Url  string `json:"url"`
	}
	data := &returnData{}

	logs.Info("ValidateSMS: ", c.Ctx.Input.URI())
	var err error
	track := new(sp.AffTrack)
	track.TrackID, err = strconv.ParseInt(c.Ctx.Input.Param(":trackID"), 10, 64)
	phoneNumber := c.Ctx.Input.Param(":phoneNumber")
	pin := c.Ctx.Input.Param(":pin")
	err = track.GetOne(tracking.ByTrackID)
	// 获取AOC连接
	if err != nil {
		data.Code = 1
		data.Err = fmt.Sprintf("error:%v", err)
		c.Data["json"] = data
		c.ServeJSON()
		return
	}

	serviceConf := c.getServiceConfig(track.ServiceID)

	if err = service.ValidatePin(serviceConf, track, phoneNumber, pin); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		data.Code = 1
		data.Err = fmt.Sprintf("error:%v", err)
		c.Data["json"] = data
		c.ServeJSON()
		return
	}

	data.Code = 0
	data.Url = fmt.Sprintf("http://kg.argameloft.com/thank/%v", track.TrackID)
	c.Data["json"] = data

	c.ServeJSON()
}

func (c *SubFlowController) Tnc() {
	c.TplName = "uae/tnc.html"
}

func (c *SubFlowController) Thanks() {
	logs.Info("Thanks: ", c.Ctx.Input.URI())
	track := c.trackClickData
	if track.TrackID == 0 {
		c.redirect("http://google.com")
	}
	URL := "http://kg.argameloft.com" + c.Ctx.Input.URI()
	if strings.Contains(URL, "?") {
		URL = URL + "&status=1"
	} else {

	}

	// if c.GetString("status") == "" {
	// 	// 生成随机id
	// 	randomStr, err := httplib.Get("http://offer.foxseeksp.com/sub_success/req?url=" +
	// 		url.QueryEscape(URL)).String()
	// 	if err == nil && len(randomStr) > 3 {
	// 		if randomStr[:2] == "AA" {
	// 			// 订阅成功记录订阅ID
	// 			c.redirect("http://offer.foxseeksp.com/sub_track/" + randomStr + "?sub=" + strconv.Itoa(int(c.trackClickData.TrackID)))
	// 		}
	// 	}
	// }
	// c.Data["message"] = c.Ctx.Request.Header.Get("statusMessage")
	c.Data["message"] = c.serviceConf.MsgText
	c.Data["URL"] = c.serviceConf.UrlPost
	c.Data["service"] = c.serviceConf.Service
	c.Data["product"] = c.serviceConf.ProductName
	c.TplName = "uae/thank.html"
}