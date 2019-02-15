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

func SubMo(mo *MoStruct) {

	switch mo.Operator {
	case "AIS", "TMH", "DTAC":
		mo.Country = "TH"
	case "DIGI", "TUNETALK", "UMOBILE", "XOX":
		mo.Country = "MY"
	}

	code, _ := InsertIntoMo(mo)
	if code != 200 {
		return
	}

	if mo.Country == "TH" && mo.Operator == "AIS" {
		go func() {
			time.Sleep(10 * time.Second)
			SendMt(mo, "0")
		}()
	}
	go func() {
		time.Sleep(20 * time.Second)
		SendMt(mo, "1")
	}()
}

func InsertIntoMo(mo *MoStruct) (int, string) {
	o := orm.NewOrm()
	var mo_old MoStruct
	o.QueryTable("mo_struct").Filter("msisdn", mo.Msisdn).Filter("shortcode", mo.Shortcode).Filter("keyword__iexact", mo.Keyword).Filter("status", 1).One(&mo_old)
	if mo_old.Id != 0 {
		logs.Debug("sub error msisdn exist: ", mo.Msisdn)
		return 201, "msisdn exist"
	}
	mo.Status = 1
	mo.Subtime = time.Now().Format("2006-01-02 15:04:05")
	mo.Unsubtime = ""

	track, err := SearchTrackById(mo.TrackId)
	if err == nil {
		mo.AffName = track.AffName
		mo.PubId = track.PubId
		mo.ProId = track.ProId
		mo.ClickId = track.ClickId
		mo.CampId = track.CampId
		mo.ProductName = track.ProductName
	}
	mo.MtSend = 0
	mo.MtDay = 2
	o.Insert(mo)

	var redis_postback uint64
	redis_postback, _ = LoadPostback(mo.CampId)

	IncrCap(mo.Shortcode, strings.ToUpper(mo.Keyword))

	_, postback := Get_postback_url(mo.CampId)
	if int(redis_postback) <= postback.Cap {
		IfIsPostback := PostbackRate(mo, postback.Rate) //回传比例  7表示回传百分之70的量  YES 表示回传postback  NO表示不回传
		if IfIsPostback {
			mo.PostbackCode, mo.PostbackMessage = PostbackRequest(mo, postback)
			mo.PostbackTime = mo.Subtime
			mo.PostbackStatus = 1
			var post = new(PostbackRecord)
			post.CampId = mo.CampId
			post.AffName = mo.AffName
			post.Proid = mo.ProId
			post.Pubid = mo.PubId
			post.Clickid = mo.ClickId
			post.Time = mo.PostbackTime
			IncrPostback(mo.CampId)
			o.Insert(post)
		}
		o.Update(mo)
	}

	return 200, "success"
}

func UnsubMo(mo *MoStruct) (int, string) {

	if strings.ToUpper(mo.Keyword) == "ALL" {
		UnsubAll(mo.Msisdn)
		return 200, "success"
	}
	o := orm.NewOrm()
	var mo_old MoStruct
	o.QueryTable("mo_struct").Filter("msisdn", mo.Msisdn).Filter("shortcode", mo.Shortcode).Filter("keyword", mo.Keyword).Filter("status", 1).One(&mo_old)
	if mo_old.Id == 0 {
		logs.Debug("unsub error no msisdn: ", mo.Msisdn)
		return 202, "no msisdn"
	}
	mo_old.Status = 0
	mo_old.Unsubtime = time.Now().Format("2006-01-02 15:04:05")
	o.Update(&mo_old)
	return 200, "success"
}

func UnsubAll(msisdn string) (int, string) {
	o := orm.NewOrm()
	o.Raw("update mo_struct set status='0' where msisdn = ? and unsubtime = ? ", msisdn, time.Now().Format("2006-01-02 15:04:05")).Exec()
	return 200, "success"
}

func SendMt(mo *MoStruct, typeOfmt string) {

	var mt MtStruct
	mt.Operator = mo.Operator
	mt.Shortcode = mo.Shortcode
	mt.Time = time.Now().Format("2006-01-02 15:04:05")
	mt.CampId = mo.CampId
	mt.Country = mo.Country
	mt.Keyword = mo.Keyword
	mt.Moid = mo.Moid
	mt.Msisdn = mo.Msisdn
	mt.Charge = typeOfmt
	mt.CampId = mo.CampId
	mt.Subid = mo.Id
	o := orm.NewOrm()
	mtid, _ := o.Insert(&mt)
	mtid_str := strconv.FormatInt(mtid, 10)

	var url, username, password string
	switch mo.Operator {
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
	reqest.Header.Add("x-request-mt-sc", mo.Shortcode)
	reqest.Header.Add("x-request-mt-msisdn", mo.Msisdn)
	reqest.Header.Add("x-request-mt-dlr", "1")

	reqest.Header.Add("x-request-mt-text", ReturnTHText(mo, typeOfmt))
	if typeOfmt == "0" {
		reqest.Header.Add("x-request-mt-charge", "0")
		reqest.Header.Add("x-request-mt-additional", "%7B%22action%22%3A%22WELCOME%22%7D")
	} else {
		reqest.Header.Add("x-request-mt-charge", "1")
	}
	reqest.Header.Add("x-request-mt-service", strings.ToUpper(mo.Keyword))
	reqest.Header.Add("x-request-mt-id", mo.Moid)
	reqest.Header.Add("x-request-mt-op", mo.Operator)
	reqest.Header.Add("x-request-mt-carryover", mtid_str)
	if mo.Country == "TH" && mo.Operator == "AIS" {
		reqest.Header.Add("x-request-mt-enc-type", "ucs2")
	}

	fmt.Println("Send MT Request:-----", mo.Id)
	for k, v := range reqest.Header {
		fmt.Println(k, ":", v)
	}

	response, err := client.Do(reqest)
	if err != nil {
		logs.Debug("Send Mt Error: ", err.Error())
		mt.ResponseErrorcode = "Http Error"
		o.Update(&mt)
		return
	}

	fmt.Println("Send MT Response:-----", mo.Id)
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

	mo.MtSend = 1
	o.Update(mo)
}

func ReturnTHText(mo *MoStruct, typeOfmt string) string {
	newFormat := time.Now().Format("020106")
	var text string
	if typeOfmt == "0" {

		switch strings.ToUpper(mo.Keyword) {
		case "FY1":
			text = "ขอบคุณที่สมัครบริการ VDO 10บาท/วัน. ยกเลิกกด *137 โทรออก สอบถามโทร022079258"
		case "GZ", "GY":
			text = "ขอบคุณที่สมัครบริการ Games 10บาท/วัน. ยกเลิกกด *137 โทรออก สอบถามโทร022079258"
		case "WZ":
			text = "ขอบคุณที่สมัครบริการ Wallpaper 10บาท/วัน. ยกเลิกกด *137 โทรออก สอบถามโทร022079258"
		}
		text = UnicodeText(text)
	} else {
		switch strings.ToUpper(mo.Keyword) {
		case "FY1":
			text = "กดดูที่นี้ " + newFormat + " http://th.mobpre.com/mt/video/aa/aa.mp4"
			text = UnicodeText(text)
		case "GZ", "GY":
			text = "กดดูที่นี้ " + newFormat + " http://static.gogamehub.com/game/1273AD/index.html"
			text = UnicodeText(text)
		case "WZ":
			text = "กดดูที่นี้ " + newFormat + " http://th.mobpre.com/mt/wallpaper/aa/aa.png"
			text = UnicodeText(text)
		case "M1", "M2", "M3":
			text = "RM5.00 Games: http://static.gogamehub.com/game/1273AD/index.html Help? 0380840028 To cancel send STOP " + strings.ToUpper(mo.Keyword) + " to 33070"
		case "M4", "M5":
			text = "RM5.00 Videos: https://s3-ap-southeast-1.amazonaws.com/ld-res/funny-video/411483.mp4 Help? 0380840028 To cancel send STOP " + strings.ToUpper(mo.Keyword) + " to 33070"
		}
	}
	return text
}

func UnicodeText(text string) string {
	ucs2HexArray := []rune(text)
	s := fmt.Sprintf("%U", ucs2HexArray)
	a := strings.Replace(s, "U+", "", -1)
	b := strings.Replace(a, "[", "", -1)
	c := strings.Replace(b, "]", "", -1)
	d := strings.Replace(c, " ", "", -1)
	return d
}
