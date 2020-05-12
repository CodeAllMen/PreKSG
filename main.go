package main

import (
	"github.com/MobileCPX/PreBaseLib/splib/click"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"

	_ "github.com/MobileCPX/PreKSG/initial"
	_ "github.com/MobileCPX/PreKSG/routers"

	"github.com/robfig/cron"
)

func init() {
	sp.InitServiceConfig()
}
func main() {
	// models.UpdateDN()

	// cr.AddFunc("0 0 9 * * ?", models.SendMtDaidly)
	// cr.AddFunc("0 0 0 * * 1", models.SetDigiMtSum)
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))

	task()
	beego.Run()
}

// 定时任务
func task() {
	cr := cron.New()

	_, _ = cr.AddFunc("0 20 */1 * * ?", SendClickDataToAdmin) // 一个小时存一次点击数据并且发送到Admin

	cr.Start()
}

func SendClickDataToAdmin() {
	sp.InsertHourClick()

	for _, service := range sp.ServiceData {
		click.SendHourData(service.CampID, click.PROD) // 发送有效点击数据
	}

}
