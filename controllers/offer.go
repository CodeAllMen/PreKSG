package controllers

import (
	"time"

	"github.com/MobileCPX/PreKSG/models"
	"github.com/astaxie/beego"
)

type Offer struct {
	beego.Controller
}

func (this *Offer) Get() {

	track := new(models.Track)
	track.Ip = this.Ctx.Request.Header.Get("X-Real-Ip")
	track.Agent = this.Ctx.Request.Header.Get("Agent")
	track.Time = time.Now().Format("2006-01-02 15:04:05")
	track.CampId = this.GetString("camp")

	var redirect_url string
	switch track.CampId {
	case "10001":
		track.AffName = "cpx"
		track.PubId = "test"
		track.ProId = "test"
		track.ClickId = "test"
		track.ShortCode = "4556066"
		track.Keyword = "FY1"
		track.ProductName = "VDO 06601"
		redirect_url = "http://api.all-apac.com/rl/bH5sbwSN?s1=999"
	default:
		this.Redirect("https://google.com", 302)
		return
	}

	num, _ := models.LoadCap(track.ShortCode, track.Keyword)
	if num >= 500 {
		this.Redirect(redirect_url, 302)
		return
	}
}
