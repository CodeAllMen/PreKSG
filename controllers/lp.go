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
	track.ProId = this.GetString("proId") // 产品名
	track.PubId = this.GetString("pubId") // 子渠道
	track.ClickId = this.GetString("clickId")
	track.Ip = this.Ctx.Request.Header.Get("X-Real-Ip")
	track.Agent = this.Ctx.Request.Header.Get("Agent")
	track.Time = time.Now().Format("2006-01-02 15:04:05")
	track.ClickStatus = "0"

	var shortCode, keyword, productName, service, description, descriptionAr, content, shortCodeDU, contentAr string
	switch strings.ToUpper(this.Ctx.Input.Param(":kw")) {
	case "GF":
		shortCode = "1111"
		shortCodeDU = "3246"
		keyword = "GF"
		productName = "Gold Finger"
		description = "amazing Mobile Games"
		descriptionAr = " ألعاب مدهشة "
		service = "game"
		content = "games"
		contentAr = " ألعاب جديدة تضاف"
	case "MYA":
		shortCode = "1111"
		shortCodeDU = "3246"
		keyword = "MYA"
		productName = "My Anime"
		description = "amazing Animated Videos"
		descriptionAr = " مقاطع فيديو متحركة مذهلة  "
		service = "anime"
		content = "videos"
		contentAr = " تتم إضافة مقاطع فيديو جديدة "
	case "POM":
		shortCode = "1111"
		shortCodeDU = "3246"
		keyword = "POM"
		productName = "Poi Movie"
		description = "latest movies"
		descriptionAr = " أحدث الأفلام "
		service = "movie"
		content = "movies"
		contentAr = " تمت إضافة أفلام جديدة "
	case "BB":
		shortCode = "1111"
		shortCodeDU = "3246"
		keyword = "BB"
		productName = "Bodybuild"
		description = "weight loss workout plan and muscle building tricks"
		descriptionAr = " خطة تمارين لفقدان الوزن وحيل بناء العضلات "
		service = "build"
		content = "content"
		contentAr = " محتوى جديد يضاف "
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
	this.Data["shortCodeDU"] = shortCodeDU
	this.Data["key"] = keyword
	this.Data["service"] = service
	this.Data["ptxid"] = id_str
	this.Data["description"] = description
	this.Data["descriptionAr"] = descriptionAr
	this.Data["content"] = content
	this.Data["contentAr"] = contentAr
	this.TplName = "uae/" + this.Ctx.Input.Param(":mode") + ".html"
}

type UAEThank struct {
	beego.Controller
}

func (this *UAEThank) Get() {
	message := this.Ctx.Request.Header.Get("statusMessage")

	var service, productName string
	this.Data["message"] = message
	switch strings.ToUpper(this.Ctx.Input.Param(":kw")) {
	case "GF":
		service = "game"
		productName = "Gold Finger"
	case "MYA":
		service = "anime"
		productName = "My Anime"
	case "POM":
		service = "movie"
		productName = "Poi Movie"
	case "BB":
		service = "build"
		productName = "Bodybuild"
	case "RC":
		service = "recipe"
		message = "Thank you for subscribing to Recipe service. You can visit the portal on http://en.recipenice.com/. "
		productName = "Recipe"
	default:
		this.Ctx.WriteString("400")
		return
	}
	this.Data["service"] = service
	this.Data["product"] = productName
	this.TplName = "uae/thank.html"
}
