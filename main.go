package main

import (
	"github.com/astaxie/beego"

	_ "github.com/MobileCPX/PreKSG/initial"
	"github.com/MobileCPX/PreKSG/models"
	_ "github.com/MobileCPX/PreKSG/routers"

	"github.com/robfig/cron"
)

func main() {
	models.Open("127.0.0.1", 6379, "mlbj")
	// models.UpdateDnTable()
	cr := cron.New()
	cr.AddFunc("0 0 0 * * ?", models.SetCap)
	cr.AddFunc("0 0 0 * * ?", models.SetPostback)
	// cr.AddFunc("0 0 9 * * ?", models.SendMtDaidly)
	// cr.AddFunc("0 0 0 * * 1", models.SetDigiMtSum)
	cr.Start()
	beego.Run()
}
