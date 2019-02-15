package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func BuFaMt() {
	o := orm.NewOrm()
	var mos []MoStruct
	o.Raw("select * from mo_struct where id in (23,24,25)").QueryRows(&mos)
	for _, v := range mos {
		if v.Operator == "DIGI" ||
			v.Operator == "UMOBILE" ||
			v.Operator == "TUNETALK" ||
			v.Operator == "XOX" {
			if v.MtSend == 2 {
				continue
			}
			if v.MtDay > 0 {
				v.MtDay = v.MtDay - 1
				o.Update(&v)
				continue
			}
		}

		var mt MtStruct
		mt.Operator = v.Operator
		mt.Shortcode = v.Shortcode
		mt.Time = time.Now().Format("2006-01-02 15:04:05")
		mt.CampId = v.CampId
		mt.Country = v.Country
		mt.Keyword = v.Keyword
		mt.Moid = v.Moid
		mt.Msisdn = v.Msisdn
		mt.Charge = "1"
		o := orm.NewOrm()
		mtid, _ := o.Insert(&mt)
		mtid_str := strconv.FormatInt(mtid, 10)

		var url, username, password string
		switch v.Operator {
		case "AIS":
			url = "http://th.gw2cloud.com/api/sender.php"
			username = "cpxth"
			password = "VLm20AZ7eKcWWqr5c"
		case "DIGI", "TUNETALK", "UMOBILE", "XOX":
			url = "http://my.gw2cloud.com/api/sender.php"
			username = "cpxmy"
			password = "O9bn30691hT5u038"
		}

		client := &http.Client{}
		reqest, err := http.NewRequest("HEAD", url, nil)
		reqest.Header.Add("x-request-mt-username", username)
		reqest.Header.Add("x-request-mt-password", password)
		reqest.Header.Add("x-request-mt-sc", v.Shortcode)
		reqest.Header.Add("x-request-mt-msisdn", v.Msisdn)
		reqest.Header.Add("x-request-mt-dlr", "1")
		reqest.Header.Add("x-request-mt-charge", "1")

		reqest.Header.Add("x-request-mt-text", ReturnDailyText(v))
		reqest.Header.Add("x-request-mt-service", strings.ToUpper(v.Keyword))
		reqest.Header.Add("x-request-mt-id", v.Moid)
		reqest.Header.Add("x-request-mt-op", v.Operator)
		reqest.Header.Add("x-request-mt-carryover", mtid_str)
		if v.Country == "TH" && v.Operator == "AIS" {
			reqest.Header.Add("x-request-mt-enc-type", "ucs2")
		}

		fmt.Println("Send MT Request:-----", v.Id)
		for k, v := range reqest.Header {
			fmt.Println(k, ":", v)
		}

		response, err := client.Do(reqest)
		if err != nil {
			logs.Debug("Send Mt Error: ", err.Error())
		}

		fmt.Println("Send MT Response:-----", v.Id)
		for k, v := range response.Header {
			fmt.Println(k, ":", v)
			switch k {
			case "X-Response-Status":
				mt.ResponseStatus = v[0]
			case "X-Response-Messageid":
				mt.ResponseMessageid = v[0]
			case "X-Response-Errorcode":
				mt.ResponseErrorcode = v[0]
			}
		}

		o.Update(&mt)
		v.MtDay = 2
		v.MtSend = v.MtSend + 1
		o.Update(&v)
	}
}

//更新mt表
func UpdateMtTable() {
	o := orm.NewOrm()
	var mts []MtStruct
	o.QueryTable("mt_struct").All(&mts)
	for i := range mts {
		var mo MoStruct
		o.QueryTable("mo_struct").Filter("moid", mts[i].Moid).One(&mo)
		mts[i].Subid = mo.Id
		mts[i].CampId = mo.CampId
		o.Update(&mts[i])
	}
}

//更新dn表
func UpdateDnTable() {
	o := orm.NewOrm()
	var dns []DnStruct
	o.QueryTable("dn_struct").All(&dns)
	for i := range dns {
		var mt MtStruct
		mtid, _ := strconv.Atoi(dns[i].Mtid)
		o.QueryTable("mt_struct").Filter("id", mtid).One(&mt)
		dns[i].Msisdn = mt.Msisdn
		dns[i].Shortcode = mt.Shortcode
		dns[i].Keyword = mt.Keyword
		dns[i].Country = mt.Country
		dns[i].Price = "10"
		dns[i].SubId = mt.Subid
		dns[i].Charge = mt.Charge
		o.Update(&dns[i])
	}
}
