package sp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/MobileCPX/PreKSG/libs"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
	ID          int64 `orm:"pk;auto;column(id)"`
	SendTime    string
	RequestId   string      `json:"requestId"`
	Transaction Transaction `json:"transaction"`
}

type ChargeNotification struct {
	Id              int64  `orm:"pk;auto"`
	Keyword         string `orm:"size(20);index"`
	Price           string `orm:"size(10);;index"`
	Time            string `orm:"size(30);index"`
	SubId           int64
	AffName         string `orm:"size(50);index"`
	PubId           string `orm:"size(100);index"`
	Charge          string
	DnStatus        int    `orm:"index"`
	RequestId       string `orm:"size(100)"`
	TransactionId   string `orm:"size(100)"`
	Shortcode       string `orm:"size(100)"`
	ChannelId       string `orm:"size(100)"`
	ApplicationId   string `orm:"size(100)"`
	Country         string `orm:"size(100)"`
	OperatorId      string `orm:"size(100)"`
	Msisdn          string `orm:"size(100)"`
	Mtid            string `orm:"size(100)"`
	ActivityTime    string `orm:"size(100)"`
	SubscriptionEnd string `orm:"size(100)"`
	Type            string `orm:"size(100)"`
	SubType         string `orm:"size(100)"`
	Status          string `orm:"size(100)"`
	Rate            string `orm:"size(100)"`
	SystemMark      int    `json:"system_mark"` // 1 为(http://offer.foxseeksp.com/offer/index), 2为(http://sp.foxseek.com/offer/index)
	SendTime        string
}

func (charge *ChargeNotification) Insert() {
	o := orm.NewOrm()
	charge.SendTime, _ = util.GetNowTime()
	_, err := o.Insert(charge)
	if err != nil {
		logs.Error("ChargeNotification Insert 存入通知数据失败 ERROR,", err.Error(), *charge)
	}
}

func (charge *ChargeNotification) GetChargeList(startTime, endTime string) (list []*ChargeNotification, err error) {
	db := orm.NewOrm()

	if _, err = db.Raw("select * from charge_notification where sub_type='RENEWAL' and status='DELIVERED' and send_time>=? and send_time<=?", startTime, endTime).QueryRows(&list); err != nil {
		err = libs.NewReportError(err)
	}

	return
}

func (charge *ChargeNotification) GetList() (list []*ChargeNotification, err error) {
	db := orm.NewOrm()

	if _, err = db.Raw("select * from charge_notification where sub_type='RENEWAL' and status='DELIVERED' and system_mark=0").QueryRows(&list); err != nil {
		err = libs.NewReportError(err)
	}

	return
}

func (charge *ChargeNotification) GetChargeListSub(startTime, endTime, startTime2, endTime2, systemMark string) (list []*ChargeNotification, err error) {
	db := orm.NewOrm()

	if _, err = db.Raw("select * from charge_notification c "+
		"where c.send_time>=? and "+
		"c.send_time<=? and c.sub_type='RENEWAL' and c.status='DELIVERED' and c.system_mark=? and  "+
		"c.msisdn in(select msisdn from charge_notification d "+
		"where d.sub_type='SUBSCRIBE' and d.status='DELIVERED' and d.send_time>=? and d.send_time<=?);", startTime, endTime, systemMark, startTime2, endTime2).QueryRows(&list); err != nil {
		err = libs.NewReportError(err)
	}

	return
}

func (charge *ChargeNotification) UpdateSystemMark() (err error) {

	o := orm.NewOrm()

	_, err = o.Update(charge)

	return
}

func SendMt(severConfig ServiceInfo, notification *ChargeNotification) {

	var urlPost string
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	var signature string
	pass := RandUpString(8)

	AddUser(urlPost, notification.Msisdn, pass)
	msgText := "Thank you for subscribing to " + severConfig.ProductName + " service. You can visit the portal on" + severConfig.UrlPost + ". Username: " + notification.Msisdn + ".Password: " + pass
	signature_url := "ApiKey=%s&ApiSecret=%v&ApplicationId=%s&CountryId=%s&OperatorId=%s" +
		"&CpId=%s&MSISDN=%s&Timestamp=%s&Lang=%s&ShortCode=%s&MsgText=%s&Method=%s"
	signature_url = fmt.Sprintf(signature_url, severConfig.ApiKey, libs.EscapeQueryParam(severConfig.ApiSecret),
		severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId,
		notification.Msisdn, URLEncodeUpper(timestamp), URLEncodeUpper("AR"), severConfig.ShortCode, strings.ToUpper(URLEncodeUpper(msgText)), URLEncodeUpper("SendSMS"))
	fmt.Println("signature_url:  ", signature_url)
	signature = HmacSha256([]byte(signature_url), []byte(severConfig.ApiSecret))

	urlOrigin := "http://ksg.kncee.com/MSG/v1.1/API/SendSMS?"
	urlParams := "applicationId=%s&countryId=%s&operatorId=%s&MSISDN=%s" +
		"&cpId=%s&requestId=%s&apiKey=%s&signature=%s&timestamp=%s&lang=%s&shortcode=%s&msgText=%s"
	urlOrigin = urlOrigin + fmt.Sprintf(urlParams, severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, notification.Msisdn, severConfig.CpId, notification.RequestId, severConfig.ApiKey, signature,
		timestamp, "ar", severConfig.ShortCode, URLEncodeUpper(msgText))
	fmt.Println("url: ", urlOrigin)
	client := &http.Client{}
	res, err := client.Get(urlOrigin)
	fmt.Println(err)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func URLEncodeUpper(str string) string {
	url_encode := url.QueryEscape(str)
	url_encode = strings.Replace(url_encode, "%2B", "%2b", -1)
	url_encode = strings.Replace(url_encode, "%2F", "%2f", -1)
	url_encode = strings.Replace(url_encode, "%3A", "%3a", -1)
	return url_encode
}

func HmacSha256(message, secret []byte) string {
	// secret := []byte("top-secret")
	// message := []byte("start1.99678678471198c6dec3-c5f0-4810-9490-e2b9f2e2d34ahttps://merch.at/cb?x=y")

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)

	// to lowercase hexits
	encode := hex.EncodeToString(hash.Sum(nil))

	return encode
}

func AddUser(reqURL, msisdn, pass string) {

	reqURL = strings.Replace(reqURL, "{msisdn}", msisdn, -1)
	reqURL = strings.Replace(reqURL, "{pass}", pass, -1)
	result, err := httplib.Get(reqURL).String()
	logs.Info("RequestService", result)
	if err != nil {
		logs.Info("添加用户或者删除用户失败，", reqURL)
	}
	// client := &http.Client{}
	// urlPost := url + "user/add?user=" + name + "&pass=" + pass + "&sign=ksg"
	// reqjson, _ := http.NewRequest("POST", urlPost, nil)
	// res, _ := client.Do(reqjson)
	// defer res.Body.Close()
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
