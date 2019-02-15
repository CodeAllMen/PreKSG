package controllers

import (
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

func Subscribe(this *APIController) string {
	ptxid := this.GetString("ptxid")
	operator := this.GetString("operator")
	res := models.Subscribe(ptxid, operator)
	return res
}

func Notification(this *APIController) string {
	body, _ := ioutil.ReadAll(this.Ctx.Request.Body)
	fmt.Println(string(body))
	return ""
}
