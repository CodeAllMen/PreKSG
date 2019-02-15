package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"github.com/MobileCPX/PreKSG/models"
)

type UAELP struct {
	beego.Controller
}

func (this *UAELP) Get() {
	track := new(models.Track)
	track.AffName = this.GetString("affName")
	track.ProId = this.GetString("proId") //产品名
	track.PubId = this.GetString("pubId") //子渠道
	track.ClickId = this.GetString("clickId")
	track.Ip = this.Ctx.Request.Header.Get("X-Real-Ip")
	track.Agent = this.Ctx.Request.Header.Get("Agent")
	track.Time = time.Now().Format("2006-01-02 15:04:05")
	track.ClickStatus = "0"

	var shortCode, keyword, productName, service string
	switch strings.ToUpper(this.Ctx.Input.Param(":kw")) {
	case "GF":
		shortCode = "1111"
		keyword = "GF"
		productName = "Gold Finger"
		service = "game"
	case "MYA":
		shortCode = "1111"
		keyword = "MYA"
		productName = "My Anime"
		service = "anime"
	case "POM":
		shortCode = "1111"
		keyword = "POM"
		productName = "Poi Movie"
		service = "movie"
	case "BB":
		shortCode = "1111"
		keyword = "BB"
		productName = "Bodybuild"
		service = "build"
	default:
		this.Ctx.WriteString("400")
		return
	}
	track.ShortCode = shortCode
	track.Keyword = keyword
	track.ProductName = productName
	// code, id := models.StartSession(track, "th-ais")
	id, _ := models.InsertTrack(track)
	id_str := strconv.FormatInt(id, 10)

	this.Data["pro"] = productName
	this.Data["code"] = shortCode
	this.Data["key"] = keyword
	this.Data["service"] = service
	this.Data["ptxid"] = id_str
	this.TplName = "uae/" + this.Ctx.Input.Param(":mode") + ".html"
}
