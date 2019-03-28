package controllers

import (
	"math/rand"
	"strconv"
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
	switch track.CampId {
	case "10000":
		track.AffName = ""
		track.ProId = "" //产品名
		track.PubId = "" //子渠道
		track.ClickId = ""

	case "21101":
		track.AffName = "godbee"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21102":
		track.AffName = "superADS"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21103":
		track.AffName = "adorca"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21104":
		track.AffName = "bitterstrawberry"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21105":
		track.AffName = "mobisummer"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21106":
		track.AffName = "dgmax"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "21107":
		track.AffName = "dgmax"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")

	case "5601":
		track.AffName = "mobvista"
		track.ProId = "jq-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	case "5602":
		track.AffName = "olimob"
		track.ProId = "jq-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
	default:
		this.Ctx.WriteString("400")
		return
	}

	randNum := rand.New(rand.NewSource(time.Now().Unix())).Intn(1)

	var shortCode, keyword, productName, service string
	switch randNum {
	case 0:
		shortCode = "1111"
		keyword = "GF"
		productName = "Gold Finger"
		service = "game"
	case 1:
		shortCode = "1111"
		keyword = "MYA"
		productName = "My Anime"
		service = "anime"
	case 2:
		shortCode = "1111"
		keyword = "POM"
		productName = "Poi Movie"
		service = "movie"
	case 3:
		shortCode = "1111"
		keyword = "BB"
		productName = "Bodybuild"
		service = "build"
	}
	track.ShortCode = shortCode
	track.Keyword = keyword
	track.ProductName = productName
	id, _ := models.InsertTrack(track)
	id_str := strconv.FormatInt(id, 10)

	this.Data["pro"] = productName
	this.Data["code"] = shortCode
	this.Data["key"] = keyword
	this.Data["service"] = service
	this.Data["ptxid"] = id_str
	this.TplName = "uae/offer.html"

}
