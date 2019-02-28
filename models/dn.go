package models

import (
	"strconv"

	"github.com/astaxie/beego/orm"
)

type Action struct {
	Type    string `json:"type"`
	SubType string `json:"subType"`
	Status  string `json:"status"`
	Rate    string `json:"rate"`
}

type Data struct {
	Shortcode       string `json:"shortcode"`
	ChannelId       string `json:"channelId"`
	ApplicationId   string `json:"applicationId"`
	CountryId       string `json:"countryId"`
	OperatorId      string `json:"operatorId"`
	Msisdn          string `json:"msisdn"`
	Action          Action `json:"action"`
	ActivityTime    string `json:"activityTime"`
	SubscriptionEnd string `json:"subscriptionEnd"`
}

type Transaction struct {
	TransactionId string `json:"transactionId"`
	Data          Data   `json:"data"`
}

type DnJson struct {
	RequestId   string      `json:"requestId"`
	Transaction Transaction `json:"transaction"`
}

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
