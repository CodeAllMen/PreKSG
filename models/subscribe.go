package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

func Subscribe(ptxid, operator string) string {

	track, err := SearchTrackById(ptxid)
	if err != nil || track.Id == 0 {
		return "201"
	}

	track.Operator = operator

	defer func() {
		o := orm.NewOrm()
		o.Update(track)
	}()

	url := "http://ksg.kncee.com/MSG/v1.1/API/RedirectToCG?MSISDN=&applicationId=%s&countryId=%s&operatorId=%s" +
		"&cpId=%s&requestId=%s&apiKey=%s&signature=%s&timestamp=%s&lang=%s&shortcode=%s" +
		"&ipAddress=%s&lpUrl=%s&rurl=%s"

	var shortCode, countryId, applicationId, cpId, apiKey, contentPage, apiSecret string
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	var signature, operatorId string

	switch track.Keyword {
	case "GF":
		shortCode = "1111"
		applicationId = "12"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "ivvT4azDdWN3UTgMPAOelOnIsscSGSKJ"
		apiSecret = "vZkXOxb70S9Os6DfYZKyay+60brtDRZZVHFYNBayA7E5dnpBf2Xsu5drtNtBty1D"
		contentPage = "http://ku.g0finger.com/"
	case "MYA":
		shortCode = "1111"
		applicationId = "13"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "5MC0F8sB2INDoYujroXAKhBml1wkpWBp"
		apiSecret = "7erVdrMdoavtY1MPQy/gJn7L63B/tj2+nHr+ccwyOYQDSCoa5b3EQOUcI4F0sHLh"
		contentPage = "http://ar.abanime.com/"
	case "POM":
		shortCode = "1111"
		applicationId = "14"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "Znvg0aF42RLalt5nFTnsUGbc4Fc5h2Sf"
		apiSecret = "tqhkRFEbpXhpk31xkCQjSmao9dlrsXOk3wZYSaJYnWROlVxJVUgAr+wQ/Lqiyj1x"
		contentPage = "http://ar.poimovie.com/"
	case "BB":
		shortCode = "1111"
		applicationId = "15"
		countryId = "247"
		operatorId = "28"
		cpId = "9"
		apiKey = "kLJ6ToymFc5yGHP6N6jYM0fq9qJdAIat"
		apiSecret = "diy3QXB6J5Ekp7BBXxvnv0ZEhuGLMAdgTJoy1zq7FOBvXviLG8RM8/IZZf8f0r4E"
		contentPage = "http://ar.fit8tube.com/"
	}

	signature_url := "ApiKey=%s&ApiSecret=%s&ApplicationId=%s&CountryId=%s&OperatorId=%s" +
		"&CpId=%s&Timestamp=%s&Lang=%s&ShortCode=%s&Method=%s"
	signature_url = fmt.Sprintf(signature_url, URLEncodeUpper(apiKey), URLEncodeUpper(apiSecret),
		URLEncodeUpper(applicationId), URLEncodeUpper(countryId), URLEncodeUpper(operatorId), URLEncodeUpper(cpId),
		URLEncodeUpper(timestamp), URLEncodeUpper("AR"), URLEncodeUpper(shortCode), URLEncodeUpper("RedirectToCG"))
	fmt.Println("signature_url:  ", signature_url)
	signature = HmacSha256([]byte(signature_url), []byte(apiSecret))

	url = fmt.Sprintf(url, applicationId, countryId, operatorId, cpId, ptxid, apiKey, signature,
		timestamp, "ar", shortCode, track.Ip, "http://kg.foxseek.com/op/"+strings.ToLower(track.Keyword), contentPage)
	fmt.Println("url: ", url)
	client := &http.Client{}
	res, err := client.Get(url)
	fmt.Println(err)
	body, _ := ioutil.ReadAll(res.Body)

	if !strings.Contains(string(body), "CGWUrl") {
		track.AocError = string(body)
		return ""
	}

	var resJson map[string]interface{}
	json.Unmarshal(body, &resJson)

	for k, v := range resJson {
		if k == "CGWUrl" {
			track.AocUrl = fmt.Sprint(v)
		}
	}
	fmt.Println(track.AocUrl)
	return track.AocUrl
}
