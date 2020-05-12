package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MobileCPX/PreKSG/libs"
	"github.com/MobileCPX/PreKSG/models/sp"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type requestPin struct {
	I3 string `json:"I3"`
	I0 string `json:"I0"`
}

func SubService(severConfig sp.ServiceInfo, track *sp.AffTrack) string {

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	signature_url := "ApiKey=%v&ApiSecret=%s&ApplicationId=%v&CountryId=%v&OperatorId=%v" +
		"&CpId=%v&Timestamp=%v&Lang=%v&ShortCode=%v&Method=%v"
	signature_url = fmt.Sprintf(signature_url, severConfig.ApiKey, libs.EscapeQueryParam(severConfig.ApiSecret),
		severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId,
		timestamp, "AR", severConfig.ShortCode, "RedirectToCG")
	// signature_url = models.URLEncodeUpper(signature_url)
	// signature_url = fmt.Sprintf(signature_url, models.URLEncodeUpper(severConfig.ApiKey), URLEncodeUpper(apiSecret),
	//	URLEncodeUpper(applicationId), URLEncodeUpper(countryId), URLEncodeUpper(operatorId), URLEncodeUpper(cpId),
	//	URLEncodeUpper(timestamp), URLEncodeUpper("AR"), URLEncodeUpper(shortCode), URLEncodeUpper("RedirectToCG"))
	fmt.Println("signature_url:  ", signature_url)
	signature := HmacSha256([]byte(signature_url), []byte(severConfig.ApiSecret))

	urlOrigin := "http://ksg.intech-mena.com/MSG/v1.1/API/RedirectToCG?MSISDN=&applicationId=%v&countryId=%v&operatorId=%v" +
		"&cpId=%v&requestId=%v&apiKey=%v&signature=%v&timestamp=%v&lang=%v&shortcode=%v" +
		"&ipAddress=%v&lpUrl=%v&rurl=%v"
	urlOrigin = fmt.Sprintf(urlOrigin,
		severConfig.ApplicationId,
		severConfig.CountryId,
		severConfig.OperatorId,
		severConfig.CpId,
		track.TrackID,
		severConfig.ApiKey,
		signature,
		timestamp, "ar",
		severConfig.ShortCode,
		track.IP,
		"http://kg.argameloft.com/lp/"+strings.ToLower(strings.Replace(severConfig.ServiceID, "-", "/", 1)),
		severConfig.RUrl+strconv.FormatInt(track.TrackID, 10))
	fmt.Println("url: ", urlOrigin)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", urlOrigin, nil)
	reqest.Header.Set("statusMessage", "http://thankyou.com")
	res, err := client.Do(reqest)

	fmt.Println(err)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	if !strings.Contains(string(body), "CGWUrl") {
		return ""
	}

	var resJson map[string]interface{}
	json.Unmarshal(body, &resJson)
	aocUrl := ""
	for k, v := range resJson {
		if k == "CGWUrl" {
			aocUrl = fmt.Sprint(v)
		}
	}
	fmt.Println(aocUrl)
	return aocUrl
}

func SubServiceSMS(severConfig sp.ServiceInfo, track *sp.AffTrack, phoneNumber string) string {
	// 构造参数， 请求 request pin接口
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	urlOrigin := "http://ksg.kncee.com/MSG/v1.1/API/RequestPinCode?"
	signatureUrl := "ApiKey=%v&ApiSecret=%s&ApplicationId=%v&CountryId=%v&OperatorId=%v" +
		"&CpId=%v&MSISDN=%v&Timestamp=%v&Lang=%v&ShortCode=%v&Method=%v"
	// urlParams := url.Values{}
	// urlParams.Add("ApiKey", severConfig.ApiKey)
	// urlParams.Add("ApiSecret", severConfig.ApiSecret)
	// urlParams.Add("ApplicationId", severConfig.ApplicationId)
	// urlParams.Add("CountryId", severConfig.CountryId)
	// urlParams.Add("OperatorId", severConfig.OperatorId)
	// urlParams.Add("CpId", severConfig.CpId)
	// urlParams.Add("MSISDN", phoneNumber)
	// urlParams.Add("Timestamp", timestamp)
	// urlParams.Add("Lang", "AR")
	// urlParams.Add("ShortCode", severConfig.ShortCode)
	// urlParams.Add("Method", "RequestPinCode")
	// apiSecret := libs.EncodeGetUrlParamValues(urlParams)

	signatureUrl = fmt.Sprintf(signatureUrl, severConfig.ApiKey, libs.EscapeQueryParam(severConfig.ApiSecret),
		severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId, phoneNumber, timestamp, "AR", severConfig.ShortCode, "RequestPinCode")
	// signatureUrl := libs.EncodeGetUrlParamValues(urlParams)
	fmt.Println("signature_url:  ", signatureUrl)
	signature := HmacSha256([]byte(signatureUrl), []byte(severConfig.ApiSecret))

	urlOrigin = urlOrigin + "MSISDN=%v&applicationId=%v&countryId=%v&operatorId=%v" +
		"&cpId=%v&requestId=%v&apiKey=%v&signature=%v&timestamp=%v&lang=%v&shortcode=%v" +
		"&ipAddress=%v&lpUrl=%v"
	urlOrigin = fmt.Sprintf(urlOrigin, phoneNumber, severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId, track.TrackID, severConfig.ApiKey, signature,
		timestamp, "ar", severConfig.ShortCode, track.IP, "http://kg.argameloft.com/lp/"+strings.ToLower(strings.Replace(severConfig.ServiceID, "-", "/", 1)))

	fmt.Println("request RequestPin url is:" + urlOrigin)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", urlOrigin, nil)
	reqest.Header.Set("statusMessage", "http://thankyou.com")
	res, err := client.Do(reqest)

	if err != nil {
		err = libs.NewReportError(err)
		return libs.GetErrorString(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("request http://ksg.kncee.com/MSG/v1.1/API/RequestPinCode, result: " + string(body))

	var requestResult requestPin

	if err = json.Unmarshal(body, &requestResult); err != nil {
		err = libs.NewReportError(err)
		return libs.GetErrorString(err)
	}

	if requestResult.I3 != "" {
		return libs.GetErrorString(libs.NewReportError(errors.New(requestResult.I3)))
	}

	if requestResult.I0 != "" {
		return libs.GetErrorString(libs.NewReportError(errors.New(requestResult.I0)))
	}

	return ""
}

func ValidatePin(severConfig sp.ServiceInfo, track *sp.AffTrack, phoneNumber, pin string) (err error) {
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	urlOrigin := "http://ksg.kncee.com/MSG/v1.1/API/ValidatePinCode?"
	signatureUrl := "ApiKey=%v&ApiSecret=%v&ApplicationId=%v&CountryId=%v&OperatorId=%v" +
		"&CpId=%v&MSISDN=%v&Timestamp=%v&Lang=%v&ShortCode=%v&Code=%v&Method=%v"
	// urlParams := url.Values{}
	// urlParams.Add("ApiSecret", severConfig.ApiSecret)
	// apiSecret := libs.EncodeGetUrlParamValues(urlParams)

	signatureUrl = fmt.Sprintf(signatureUrl, severConfig.ApiKey, libs.EscapeQueryParam(severConfig.ApiSecret),
		severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId,
		phoneNumber, timestamp, "AR", severConfig.ShortCode, pin, "ValidatePinCode")
	fmt.Println("signature_url:  ", signatureUrl)
	signature := HmacSha256([]byte(signatureUrl), []byte(severConfig.ApiSecret))

	urlOrigin = urlOrigin + "MSISDN=%v&applicationId=%v&countryId=%v&operatorId=%v" +
		"&cpId=%v&requestId=%v&apiKey=%v&signature=%v&timestamp=%v&lang=%v&shortcode=%v&code=%v"
	urlOrigin = fmt.Sprintf(urlOrigin, phoneNumber, severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId, track.TrackID, severConfig.ApiKey, signature,
		timestamp, "ar", severConfig.ShortCode, pin)

	fmt.Println("request ValidatePin url is:" + urlOrigin)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", urlOrigin, nil)
	reqest.Header.Set("statusMessage", "http://thankyou.com")
	res, err := client.Do(reqest)

	if err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("request http://ksg.kncee.com/MSG/v1.1/API/ValidatePinCode, result: " + string(body))

	var requestResult requestPin

	if err = json.Unmarshal(body, &requestResult); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
		return
	}

	if requestResult.I3 != "" {
		return libs.NewReportError(errors.New(requestResult.I3))
	}

	if requestResult.I0 != "" {
		return libs.NewReportError(errors.New(requestResult.I0))
	}

	return
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
