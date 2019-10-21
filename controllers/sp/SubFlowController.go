package sp

import (
	"fmt"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/MobileCPX/PreKSG/service"
	"github.com/astaxie/beego/logs"
)

type SubFlowController struct {
	BaseController

	// 追踪点击数据
	trackClickData *sp.AffTrack

	//serviceConf
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
	//获取AOC连接
	//track :=c.trackClickData
	res := service.SubService(c.serviceConf, c.trackClickData)
	//c.Data["json"] = map[string]string{
	//	"data": res,
	//}
	if res == "" {
		c.RedirectURL("http://google.com")
	} else {
		c.RedirectURL(res)
	}
	fmt.Println(res)
}

func (c *SubFlowController) Thanks() {
	logs.Info("Thanks: ", c.Ctx.Input.URI())
	track := c.trackClickData
	if track.TrackID == 0 {
		c.redirect("http://google.com")
	}
	c.Data["message"] = c.Ctx.Request.Header.Get("statusMessage")

	c.Data["service"] = c.serviceConf.Service
	c.Data["product"] = c.serviceConf.ProductName
	c.TplName = "uae/thank.html"
}
