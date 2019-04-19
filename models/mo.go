package models

import (
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func InsertIntoMo(dn DnStruct) (int, string) {
	o := orm.NewOrm()
	var mo_old MoStruct
	o.QueryTable("mo_struct").Filter("msisdn", dn.Msisdn).Filter("shortcode", dn.Shortcode).Filter("keyword", dn.Keyword).One(&mo_old)
	if mo_old.Id != 0 {
		logs.Debug("sub error msisdn exist: ", dn.Msisdn)
		return 201, "msisdn exist"
	}
	var mo MoStruct
	mo.Msisdn = dn.Msisdn
	mo.Operator = dn.OperatorId
	mo.Country = dn.Country
	mo.Shortcode = dn.Shortcode
	mo.Keyword = dn.Keyword
	mo.Subtime = dn.Time
	mo.Status = 1
	mo.Unsubtime = ""

	track, err := SearchTrackById(dn.TransactionId)
	if err == nil {
		mo.AffName = track.AffName
		mo.PubId = track.PubId
		mo.ProId = track.ProId
		mo.ClickId = track.ClickId
		mo.CampId = track.CampId
		mo.ProductName = track.ProductName
	}

	o.Insert(&mo)

	var redis_postback uint64
	redis_postback, _ = LoadPostback(mo.CampId)

	IncrCap(mo.Shortcode, strings.ToUpper(mo.Keyword))

	_, postback := Get_postback_url(mo.CampId)
	if int(redis_postback) <= postback.Cap {
		IfIsPostback := PostbackRate(&mo, postback.Rate) //回传比例  7表示回传百分之70的量  YES 表示回传postback  NO表示不回传
		if IfIsPostback {
			mo.PostbackCode, mo.PostbackMessage = PostbackRequest(&mo, postback)
			mo.PostbackTime = mo.Subtime
			mo.PostbackStatus = 1
			var post = new(PostbackRecord)
			post.CampId = mo.CampId
			post.AffName = mo.AffName
			post.Proid = mo.ProId
			post.Pubid = mo.PubId
			post.Clickid = mo.ClickId
			post.Time = mo.PostbackTime
			post.CampId = mo.CampId
			IncrPostback(mo.CampId)
			o.Insert(post)
		}
		o.Update(mo)
	}

	return 200, "success"
}

func UnsubMo(dn DnStruct) (int, string) {

	o := orm.NewOrm()
	var mo_old MoStruct
	o.QueryTable("mo_struct").Filter("msisdn", dn.Msisdn).Filter("shortcode", dn.Shortcode).Filter("keyword", dn.Keyword).Filter("status", 1).One(&mo_old)
	if mo_old.Id == 0 {
		logs.Debug("unsub error no msisdn: ", dn.Msisdn)
		return 202, "no msisdn"
	}
	mo_old.Status = 0
	mo_old.Unsubtime = time.Now().Format("2006-01-02 15:04:05")
	o.Update(&mo_old)
	return 200, "success"
}
