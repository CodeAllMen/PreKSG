package controllers

import (
	"github.com/MobileCPX/PreKSG/models"
	"github.com/astaxie/beego"
)

type AffController struct {
	beego.Controller
}

//渠道转化
func (this *AffController) Get() {
	start_time := this.GetString("start_time")
	end_time := this.GetString("end_time")
	keyword := this.GetString("keyword")
	operator := this.GetString("operator")
	aff_name := this.GetString("aff_name")
	err, data := models.GetAffdDate(start_time, end_time, keyword, operator, aff_name)
	if err == nil {
		this.Data["json"] =
			map[string]interface{}{
				"code":    "1",
				"data":    data,
				"message": "success",
			}
		this.ServeJSON()
	} else {
		this.Data["json"] =
			map[string]interface{}{
				"code":    0,
				"message": "failed",
			}
		this.ServeJSON()
	}
}
