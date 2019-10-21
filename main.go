package main

import (
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"

	_ "github.com/MobileCPX/PreKSG/initial"
	"github.com/MobileCPX/PreKSG/models"
	_ "github.com/MobileCPX/PreKSG/routers"

	"github.com/robfig/cron"
)

func init() {
	sp.InitServiceConfig()
}
func main() {
	models.Open("127.0.0.1", 6379, "mlbj")
	// models.UpdateDN()
	cr := cron.New()
	cr.AddFunc("0 0 0 * * ?", models.SetCap)
	cr.AddFunc("0 0 0 * * ?", models.SetPostback)
	// cr.AddFunc("0 0 9 * * ?", models.SendMtDaidly)
	// cr.AddFunc("0 0 0 * * 1", models.SetDigiMtSum)
	cr.Start()
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
	}))
	beego.Run()
}
