package sp

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/MobileCPX/PreKSG/service"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"net/url"
	"strconv"
	"strings"
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
	var err error
	track := new(sp.AffTrack)
	track.TrackID, err = strconv.ParseInt(c.Ctx.Input.Param(":trackID"), 10, 64)
	err = track.GetOne(tracking.ByTrackID)
	//获取AOC连接
	if err != nil {
		c.RedirectURL("http://google.com")
		return
	}
	serviceConf := c.getServiceConfig(track.ServiceID)
	res := service.SubService(serviceConf, track)
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
	URL := "http://kg.foxseek.com" + c.Ctx.Input.URI()
	if strings.Contains(URL, "?") {
		URL = URL + "&status=1"
	} else {

	}

	if c.GetString("status") == "" {
		// 生成随机id
		randomStr, err := httplib.Get("http://offer.globaltraffictracking.com/sub_success/req?url=" +
			url.QueryEscape(URL)).String()
		if err == nil && len(randomStr) > 3 {
			if randomStr[:2] == "AA" {
				//订阅成功记录订阅ID
				c.redirect("http://offer.globaltraffictracking.com/sub_track/" + randomStr + "?sub=" + strconv.Itoa(int(c.trackClickData.TrackID)))
			}
		}
	}
	c.Data["message"] = c.Ctx.Request.Header.Get("statusMessage")
	c.Data["URL"] = c.serviceConf.UrlPost
	c.Data["service"] = c.serviceConf.Service
	c.Data["product"] = c.serviceConf.ProductName
	c.TplName = "uae/thank.html"
}
