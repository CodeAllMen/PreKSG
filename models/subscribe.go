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

	fmt.Println(operator)

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

	var shortCode, countryId, applicationId, cpId, apiKey, apiSecret, rurl string
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	var signature, operatorId string

	switch strings.ToUpper(track.Keyword) {
	case "GF":
		switch operator {
		//1 Etisalat 2 DU
		case "1":
			shortCode = "1111"
			operatorId = "28"
			apiKey = "ivvT4azDdWN3UTgMPAOelOnIsscSGSKJ"
			apiSecret = "vZkXOxb70S9Os6DfYZKyay+60brtDRZZVHFYNBayA7E5dnpBf2Xsu5drtNtBty1D"
			rurl = "http://ku.g0finger.com/"
		case "2":
			shortCode = "3246"
			operatorId = "32"
			apiKey = "7xlO9uuUktUka8CFlHPfJsDzzlaS51vD"
			apiSecret = "/JACZwinCd+WVSd+mNxJFAMvzWNWKYXmia1oA+6vPrUj7c7eMI29Yf1h+b6zd+Un"
			rurl = "http://kg.foxseek.com/thank/gf"
		}
		applicationId = "12"
		countryId = "247"
		cpId = "9"
	case "MYA":
		switch operator {
		//1 Etisalat 2 DU
		case "1":
			shortCode = "1111"
			operatorId = "28"
			apiKey = "5MC0F8sB2INDoYujroXAKhBml1wkpWBp"
			apiSecret = "7erVdrMdoavtY1MPQy/gJn7L63B/tj2+nHr+ccwyOYQDSCoa5b3EQOUcI4F0sHLh"
			rurl = "http://ar.abanime.com/"
		case "2":
			shortCode = "3246"
			operatorId = "32"
			apiKey = "xxKHAnaleglcp6He9KEk85pABMJSXmwK"
			apiSecret = "ncVAQM4VE8SrRqZAKW4MWtia/66PBs37c8PtoQgUReU02Mjwe06QSl+P1OIwJTx0"
			rurl = "http://kg.foxseek.com/thank/mya"
		}
		applicationId = "13"
		countryId = "247"
		cpId = "9"
	case "POM":
		switch operator {
		//1 Etisalat 2 DU
		case "1":
			shortCode = "1111"
			operatorId = "28"
			apiKey = "Znvg0aF42RLalt5nFTnsUGbc4Fc5h2Sf"
			apiSecret = "tqhkRFEbpXhpk31xkCQjSmao9dlrsXOk3wZYSaJYnWROlVxJVUgAr+wQ/Lqiyj1x"
			rurl = "http://ar.poimovie.com/"
		case "2":
			shortCode = "3246"
			operatorId = "32"
			apiKey = "czoFxiJ3HfS6EZwZjHhIe41v3J1AizyS"
			apiSecret = "bB9csgdZ7wkt9ryFQXDx8y7/ozlv9gmgCkcQIBZKOLDoOSqiVd5ri9Pf0N8SdgDw"
			rurl = "http://kg.foxseek.com/thank/pom"
		}
		applicationId = "14"
		countryId = "247"
		cpId = "9"
	case "BB":
		switch operator {
		//1 Etisalat 2 DU
		case "1":
			shortCode = "1111"
			operatorId = "28"
			apiKey = "kLJ6ToymFc5yGHP6N6jYM0fq9qJdAIat"
			apiSecret = "diy3QXB6J5Ekp7BBXxvnv0ZEhuGLMAdgTJoy1zq7FOBvXviLG8RM8/IZZf8f0r4E"
			rurl = "http://at.fitnessnice.com/"
		case "2":
			shortCode = "3246"
			operatorId = "32"
			apiKey = "19cN7sCClgx2OLjlkHwHcTKcCdqxdY6G"
			apiSecret = "O0yYy5tPQ7noIfksc3YA0nJc3ch+AEoI6GbWcdb1vpMnWfV1YFR8LILj4ooHjwKn"
			rurl = "http://kg.foxseek.com/thank/bb"
		}
		applicationId = "15"
		countryId = "247"
		cpId = "9"
	}

	signature_url := "ApiKey=%s&ApiSecret=%s&ApplicationId=%s&CountryId=%s&OperatorId=%s" +
		"&CpId=%s&Timestamp=%s&Lang=%s&ShortCode=%s&Method=%s"
	signature_url = fmt.Sprintf(signature_url, URLEncodeUpper(apiKey), URLEncodeUpper(apiSecret),
		URLEncodeUpper(applicationId), URLEncodeUpper(countryId), URLEncodeUpper(operatorId), URLEncodeUpper(cpId),
		URLEncodeUpper(timestamp), URLEncodeUpper("AR"), URLEncodeUpper(shortCode), URLEncodeUpper("RedirectToCG"))
	fmt.Println("signature_url:  ", signature_url)
	signature = HmacSha256([]byte(signature_url), []byte(apiSecret))

	url = fmt.Sprintf(url, applicationId, countryId, operatorId, cpId, ptxid, apiKey, signature,
		timestamp, "ar", shortCode, track.Ip, "http://kg.foxseek.com/op/"+strings.ToLower(track.Keyword), rurl)
	fmt.Println("url: ", url)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("statusMessage", "http://thankyou.com")
	res, err := client.Do(reqest)

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
