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
	var op string
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
		op = "et"
	case "21102":
		track.AffName = "superADS"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21103":
		track.AffName = "adorca"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21104":
		track.AffName = "bitterstrawberry"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21105":
		track.AffName = "mobisummer"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21106":
		track.AffName = "dgmax"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21107":
		track.AffName = "dgmax"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21108":
		track.AffName = "Vene"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"

	case "5601":
		track.AffName = "mobvista"
		track.ProId = "jq-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5602":
		track.AffName = "olimob"
		track.ProId = "jq-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5603":
		track.AffName = "mobipium"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5604":
		track.AffName = "clickmob"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5605":
		track.AffName = "adorca"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5606":
		track.AffName = "hyperclick"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5607":
		track.AffName = "hyperclick"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5608":
		track.AffName = "funnymobi"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5609":
		track.AffName = "funnymobi"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5610":
		track.AffName = "mobisummer"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5611":
		track.AffName = "mobisummer"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5612":
		track.AffName = "addiliate"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5613":
		track.AffName = "addiliate"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5614":
		track.AffName = "mobrider"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5615":
		track.AffName = "mobrider"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5616":
		track.AffName = "flex"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5617":
		track.AffName = "flex"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5618":
		track.AffName = "funsmobi"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5619":
		track.AffName = "funsmobi"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "5620":
		track.AffName = "mobisense"
		track.ProId = "zyy-uae-et"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "5621":
		track.AffName = "mobisense"
		track.ProId = "zyy-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"

	case "5622":
		track.AffName = "gadmobe"
		track.ProId = "yyz-uae-du"          //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"

	case "21109":
		track.AffName = "olimob"
		track.ProId = "jq-uae-du"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "21110":
		track.AffName = "clickmob"
		track.ProId = "jq-uae-du"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "21111":
		track.AffName = "olimob"
		track.ProId = "jq-uae-du"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "du"
	case "21112":
		track.AffName = "olimob"
		track.ProId = "jq-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	case "21113":
		track.AffName = "Avatar"
		track.ProId = "yyz-uae-et"           //产品名
		track.PubId = this.GetString("aid") //子渠道
		track.ClickId = this.GetString("cid")
		op = "et"
	default:
		this.Ctx.WriteString("400")
		return
	}

	randNum := rand.New(rand.NewSource(time.Now().Unix())).Intn(1)
	// 选择产品
	productID, _ := this.GetInt("p")
	if productID != 0 && productID < 4 {
		randNum = productID
	}

	var shortCode, keyword, productName, service string
	switch randNum {
	case 0:
		shortCode = "1111"
		keyword = "BB"
		productName = "Bodybuild"
		service = "build"
	// case 0:
	// 	shortCode = "1111"
	// 	keyword = "GF"
	// 	productName = "Gold Finger"
	// 	service = "game"
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
	this.TplName = "uae/" + op + ".html"

}
