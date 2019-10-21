package routers

import (
	"github.com/MobileCPX/PreKSG/controllers"
	"github.com/MobileCPX/PreKSG/controllers/sp"
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

	beego.Router("/lp/:serviceType/:operator", &sp.LpController{}, "*:LpSub")
	beego.Router("/sub/req/:trackID", &sp.SubFlowController{}, "*:SubReq")
	beego.Router("/thank/:trackID", &sp.SubFlowController{}, "*:Thanks")

	//postback
	beego.Router("/set/postback", &sp.SetPostbackController{})
	//追踪连接
	beego.Router("/aff_click", &sp.TrackingController{}, "*:InsertAffClick")
}
