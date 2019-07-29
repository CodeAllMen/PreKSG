package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func InsertIntoDn(dnJson DnJson) {
	o := orm.NewOrm()
	var dn DnStruct
	dn.RequestId = dnJson.RequestId
	dn.TransactionId = dnJson.Transaction.TransactionId
	dn.Shortcode = dnJson.Transaction.Data.Shortcode
	dn.ChannelId = dnJson.Transaction.Data.ChannelId
	dn.ApplicationId = dnJson.Transaction.Data.ApplicationId
	dn.Country = dnJson.Transaction.Data.CountryId
	dn.OperatorId = dnJson.Transaction.Data.OperatorId
	dn.Msisdn = dnJson.Transaction.Data.Msisdn
	dn.Type = dnJson.Transaction.Data.Action.Type
	dn.SubType = dnJson.Transaction.Data.Action.SubType
	dn.Status = dnJson.Transaction.Data.Action.Status
	dn.Rate = dnJson.Transaction.Data.Action.Rate
	dn.ActivityTime = dnJson.Transaction.Data.ActivityTime
	dn.SubscriptionEnd = dnJson.Transaction.Data.SubscriptionEnd
	dn.Time = time.Now().Format("2006-01-02 15:04:05")
	var mo MoStruct
	o.QueryTable("mo_struct").Filter("track_id", dn.TransactionId).One(&mo)
	dn.SubId = mo.Id
	o.Insert(&dn)

	if dn.SubType == "SUBSCRIBE" && dn.Status == "DELIVERED" {
		SendSubMt(dn)
		InsertIntoMo(dn)
	}

	if dn.SubType == "UNSUBSCRIBE" {
		UnsubMo(dn)
	}
}

func SendSubMt(dn DnStruct) {

	dnIdStr := strconv.FormatInt(dn.Id, 10)

	url := "http://ksg.kncee.com/MSG/v1.1/API/SendSMS?applicationId=%s&countryId=%s&operatorId=%s&MSISDN=%s" +
		"&cpId=%s&requestId=%s&apiKey=%s&signature=%s&timestamp=%s&lang=%s&shortcode=%s&msgText=%s"

	var shortCode, countryId, applicationId, cpId, apiKey, apiSecret, urlPost string
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	var signature, operatorId, msgText string
	pass := RandUpString(8)

	switch dn.ApplicationId {
	case "12":
		shortCode = "1111"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "ivvT4azDdWN3UTgMPAOelOnIsscSGSKJ"
		apiSecret = "vZkXOxb70S9Os6DfYZKyay+60brtDRZZVHFYNBayA7E5dnpBf2Xsu5drtNtBty1D"
		msgText = "Thank you for subscribing to Gold Finger service. You can visit the portal on http://ku.g0finger.com/. Username: " + dn.Msisdn + ".Password: " + pass
		applicationId = "12"
		urlPost = "http://ku.g0finger.com/"
	case "13":
		shortCode = "1111"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "5MC0F8sB2INDoYujroXAKhBml1wkpWBp"
		apiSecret = "7erVdrMdoavtY1MPQy/gJn7L63B/tj2+nHr+ccwyOYQDSCoa5b3EQOUcI4F0sHLh"
		msgText = "Thank you for subscribing to Gold Finger service. You can visit the portal on http://ar.abanime.com/. Username: " + dn.Msisdn + ".Password: " + pass
		applicationId = "13"
		urlPost = "http://ar.abanime.com/"
	case "14":
		shortCode = "1111"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "Znvg0aF42RLalt5nFTnsUGbc4Fc5h2Sf"
		apiSecret = "tqhkRFEbpXhpk31xkCQjSmao9dlrsXOk3wZYSaJYnWROlVxJVUgAr+wQ/Lqiyj1x"
		msgText = "Thank you for subscribing to Gold Finger service. You can visit the portal on http://ar.poimovie.com/. Username: " + dn.Msisdn + ".Password: " + pass
		applicationId = "14"
		urlPost = "http://ar.poimovie.com/"
	case "15":
		shortCode = "1111"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "kLJ6ToymFc5yGHP6N6jYM0fq9qJdAIat"
		apiSecret = "diy3QXB6J5Ekp7BBXxvnv0ZEhuGLMAdgTJoy1zq7FOBvXviLG8RM8/IZZf8f0r4E"
		msgText = "Thank you for subscribing to Gold Finger service. You can visit the portal on http://ar.fit8tube.com/. Username: " + dn.Msisdn + ".Password: " + pass
		applicationId = "15"
		urlPost = "http://ar.fit8tube.com/"
	}

	AddUser(urlPost, dn.Msisdn, pass)

	signature_url := "ApiKey=%s&ApiSecret=%s&ApplicationId=%s&CountryId=%s&OperatorId=%s" +
		"&CpId=%s&MSISDN=%s&Timestamp=%s&Lang=%s&ShortCode=%s&MsgText=%s&Method=%s"
	signature_url = fmt.Sprintf(signature_url, URLEncodeUpper(apiKey), URLEncodeUpper(apiSecret),
		URLEncodeUpper(applicationId), URLEncodeUpper(countryId), URLEncodeUpper(operatorId), URLEncodeUpper(cpId),
		URLEncodeUpper(dn.Msisdn), URLEncodeUpper(timestamp), URLEncodeUpper("AR"), URLEncodeUpper(shortCode), strings.ToUpper(URLEncodeUpper(msgText)), URLEncodeUpper("SendSMS"))
	fmt.Println("signature_url:  ", signature_url)
	signature = HmacSha256([]byte(signature_url), []byte(apiSecret))

	url = fmt.Sprintf(url, applicationId, countryId, operatorId, dn.Msisdn, cpId, dnIdStr, apiKey, signature,
		timestamp, "ar", shortCode, URLEncodeUpper(msgText))
	fmt.Println("url: ", url)
	client := &http.Client{}
	res, err := client.Get(url)
	fmt.Println(err)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

}

func AddUser(url, name, pass string) {

	client := &http.Client{}
	urlPost := url + "user/add?user=" + name + "&pass=" + pass + "&sign=ksg"
	reqjson, _ := http.NewRequest("POST", urlPost, nil)
	res, _ := client.Do(reqjson)
	defer res.Body.Close()
}

func RandUpString(l int) string {
	var result bytes.Buffer
	var temp byte
	for i := 0; i < l; {
		if RandInt(48, 57) != temp {
			temp = RandInt(48, 57)
			result.WriteByte(temp)
			i++
		}
	}
	return result.String()
}

func RandInt(min int, max int) byte {
	rand.Seed(time.Now().UnixNano())
	return byte(min + rand.Intn(max-min))
}
