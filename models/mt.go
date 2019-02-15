package models

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var GameUrl = []string{"6121SW", "4361NH", "7785IO", "2351HF", "1268WW", "1290WE", "1253FR", "1179VF", "1259OP", "2627SG", "4554RR", "5678ER", "5315GF", "4532TY", "4431UI", "5633QB", "1123EW", "5755YU", "3324TT", "40214", "3451TY", "4511FB", "1218GG", "7885TB", "40071", "4521KL", "5642BN", "1124KC", "1249PA", "1267LO", "5632TM", "1278PS", "1820KD", "6533LM", "6272FB", "3244TY", "40337"}

func SendMtDaidly() {
	o := orm.NewOrm()
	var mos []MoStruct
	o.Raw("select * from mo_struct where status='1' and left(subtime,10)<?", time.Now().Format("2006-01-02")).QueryRows(&mos)
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
		mt.CampId = v.CampId
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

func ReturnDailyText(mo MoStruct) string {
	newFormat := time.Now().Format("020106")
	var text string
	switch strings.ToUpper(mo.Keyword) {
	case "FY1":
		text = "กดดูที่นี้ " + newFormat + " http://static.ifunnyhub.com/video/526068.mp4"
		text = UnicodeText(text)
	case "GZ", "GY":

		randFormat := time.Now().Format("2006-01-02")
		t, _ := time.Parse("2006-01-02", randFormat)
		seed := rand.New(rand.NewSource(t.UnixNano()))
		rand_num := seed.Intn(37)
		game := "http://static.gogamehub.com/game/" + GameUrl[rand_num] + "/index.html"

		text = "กดดูที่นี้ " + newFormat + " " + game
		text = UnicodeText(text)
	case "WZ":

		randFormat := time.Now().Format("2006-01-02")
		t, _ := time.Parse("2006-01-02", randFormat)
		seed := rand.New(rand.NewSource(t.UnixNano()))
		rand_num := seed.Intn(95)
		rand_num_str := strconv.Itoa(rand_num)
		text = "กดดูที่นี้ " + newFormat + " http://th.mobpre.com/static/wallpaper/" + rand_num_str + ".jpg"
		text = UnicodeText(text)
	case "M1", "M2", "M3":

		newFormat := time.Now().Format("2006-01-02")
		t, _ := time.Parse("2006-01-02", newFormat)
		seed := rand.New(rand.NewSource(t.UnixNano()))
		rand_num := seed.Intn(37)
		game := "http://static.gogamehub.com/game/" + GameUrl[rand_num] + "/index.html"
		text = "RM5.00 Games: " + game + " Help? 0380840028 To cancel send STOP " + strings.ToUpper(mo.Keyword) + " to 33070"
	case "M4", "M5":
		text = "RM5.00 Videos: https://s3-ap-southeast-1.amazonaws.com/ld-res/funny-video/411483.mp4 Help? 0380840028 To cancel send STOP " + strings.ToUpper(mo.Keyword) + " to 33070"
	}
	return text
}

func SetDigiMtSum() {
	o := orm.NewOrm()
	o.Raw("update mo_struct set mt_send = 0 , mt_day = 0").Exec()
}
