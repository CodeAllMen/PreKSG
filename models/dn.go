package models

import (
	"strconv"

	"github.com/astaxie/beego/orm"
)

func InsertIntoDn(dn *DnStruct) {
	o := orm.NewOrm()
	if dn.Mtid != "" {
		mtid, _ := strconv.Atoi(dn.Mtid)
		var mt_struct MtStruct
		o.QueryTable("mt_struct").Filter("id", mtid).One(&mt_struct)
		dn.Msisdn = mt_struct.Msisdn
		dn.Shortcode = mt_struct.Shortcode
		dn.Keyword = mt_struct.Keyword
		dn.Country = mt_struct.Country
		dn.Price = "10"
		dn.SubId = mt_struct.Subid
		dn.Charge = mt_struct.Charge
		if dn.Status == "success" {
			dn.DnStatus = 1
		} else {
			dn.DnStatus = 0
		}
	}
	o.Insert(dn)
}
