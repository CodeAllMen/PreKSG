package sp

import (
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"strconv"
	"strings"
)

type LpController struct {
	BaseController
}

func (c *LpController) LpSub() {
	page := c.Ctx.Input.Param(":serviceType")
	operator := c.Ctx.Input.Param(":operator")
	trackID := ""
	serviceID := strings.ToUpper(page) + "-" + strings.ToUpper(operator)
	serviceConfig := c.getServiceConfig(serviceID)
	if c.GetString("tid") != "" {
		trackID = c.GetString("tid")
	} else {
		// LP 页面存入此次点击信息，获取aff_track 表自增ID
		trackID = tracking.LpPageTracking(c.Ctx.Request, "http://kg.argameloft.com/aff/click", serviceConfig.ServiceID)
		// 将trackID转为int类型，判断trackID是否为数字类型
		_, err := strconv.Atoi(trackID)

		if err != nil { // 说明trackID不是int类型，不能订阅服务
			c.Ctx.ResponseWriter.ResponseWriter.WriteHeader(404)
			c.StopRun()
		}
	}

	c.Data["pro"] = serviceConfig.ProductName
	c.Data["code"] = serviceConfig.ShortCode

	if serviceConfig.KeyWord == "MA" {
		if serviceConfig.ShortCode == "1111" {
			c.Data["key"] = "MYA"
		} else {
			c.Data["key"] = serviceConfig.KeyWord
		}
	} else {
		c.Data["key"] = serviceConfig.KeyWord
	}

	c.Data["service"] = serviceConfig.Service
	// c.Data["ptxid"] = id_str
	c.Data["description"] = serviceConfig.Description
	c.Data["descriptionAr"] = serviceConfig.DescriptionAr
	c.Data["content"] = serviceConfig.Content
	c.Data["contentAr"] = serviceConfig.DescriptionAr
	c.Data["UrlPost"] = serviceConfig.UrlPost
	c.Data["Price"] = serviceConfig.Price

	// ET的
	if serviceConfig.ShortCode == "1111" {
		c.Data["URL"] = "/api/sub_sms/" + trackID + "/" + operator
		c.TplName = "uae/" + operator + ".html"
	} else {
		// DU的
		c.Data["URL"] = "/api/sub/" + trackID + "/" + operator
		c.TplName = "uae_sms/" + operator + ".html"
	}

}
