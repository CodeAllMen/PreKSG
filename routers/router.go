package routers

import (
	"github.com/MobileCPX/PreKSG/controllers"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/of", &controllers.Offer{})
	beego.Router("/thank/:kw", &controllers.UAEThank{})

	//泰国
	beego.Router("/:mode/:kw", &controllers.UAELP{})
	beego.Router("/api/:mode", &controllers.APIController{})

	//数据查询
	beego.Router("/aff_data", &controllers.AffController{}) //获取渠道订阅信息
}
