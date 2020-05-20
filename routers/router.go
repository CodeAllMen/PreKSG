package routers

import (
	"github.com/MobileCPX/PreBaseLib/splib/postback"
	"github.com/MobileCPX/PreKSG/controllers"
	"github.com/MobileCPX/PreKSG/controllers/sp"
	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/of", &controllers.Offer{})
	// beego.Router("/thank/:kw", &controllers.UAEThank{})

	// 泰国
	beego.Router("/:mode/:kw", &controllers.UAELP{})
	beego.Router("/api/:mode", &controllers.APIController{})

	// 数据查询  上面的router都是已经没用了的，不用管
	beego.Router("/aff_data", &controllers.AffController{}) // 获取渠道订阅信息

	beego.Router("/lp/:serviceType/:operator", &sp.LpController{}, "*:LpSub")
	beego.Router("/api/sub/:trackID/:op", &sp.SubFlowController{}, "*:SubReq")
	beego.Router("/api/sub_sms/:trackID/:op/:phoneNumber", &sp.SubFlowController{}, "*:SubReqSMS")
	beego.Router("/api/validate_sms/:trackID/:phoneNumber/:pin", &sp.SubFlowController{}, "*:ValidateSMS")
	beego.Router("/api/dn", &sp.NotificationController{})

	beego.Router("/thank/:trackID", &sp.SubFlowController{}, "GET:Thanks")

	beego.Router("/tnc/:trackID", &sp.SubFlowController{}, "GET:Tnc")

	// postback
	beego.Router("/set/postback", &postback.SetPostbackController{})
	// 追踪连接
	beego.Router("/aff/click", &sp.TrackingController{}, "*:InsertAffClick")

	beego.Router("/count", &sp.CountController{}, "Get:Count")

	beego.Router("/count_sub", &sp.CountController{}, "Get:CountSub")
}
