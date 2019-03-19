package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MobileCPX/PreKSG/models"
	"github.com/astaxie/beego"
)

type APIController struct {
	beego.Controller
}

func (this *APIController) Get() {
	mode := this.Ctx.Input.Param(":mode")
	var res string
	switch mode {
	case "sub":
		res = Subscribe(this)
	case "dn":
		res = Notification(this)
	}
	this.Data["json"] = map[string]string{
		"data": res,
	}

	fmt.Println(res)
	this.ServeJSON()
}

func (this *APIController) Post() {
	var res string
	res = Notification(this)
	this.Ctx.WriteString("1")
	fmt.Println(res)
}

func Subscribe(this *APIController) string {
	ptxid := this.GetString("ptxid")
	operator := this.GetString("op")
	res := models.Subscribe(ptxid, operator)
	return res
}

func Notification(this *APIController) string {
	body, _ := ioutil.ReadAll(this.Ctx.Request.Body)
	var dnJson models.DnJson
	json.Unmarshal([]byte(body), &dnJson)
	fmt.Println(string(body))
	fmt.Println(dnJson)
	models.InsertIntoDn(dnJson)
	return ""
}
