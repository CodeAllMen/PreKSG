package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/MobileCPX/PreKSG/models/sp"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SubService(severConfig sp.ServiceInfo, track *sp.AffTrack) string {

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	signature_url := "ApiKey=%s&ApiSecret=%s&ApplicationId=%s&CountryId=%s&OperatorId=%s" +
		"&CpId=%s&Timestamp=%s&Lang=%s&ShortCode=%s&Method=%s"
	signature_url = fmt.Sprintf(signature_url, severConfig.ApiKey, severConfig.ApiSecret,
		severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId,
		timestamp, "AR", severConfig.ShortCode, "RedirectToCG")
	//signature_url = models.URLEncodeUpper(signature_url)
	//signature_url = fmt.Sprintf(signature_url, models.URLEncodeUpper(severConfig.ApiKey), URLEncodeUpper(apiSecret),
	//	URLEncodeUpper(applicationId), URLEncodeUpper(countryId), URLEncodeUpper(operatorId), URLEncodeUpper(cpId),
	//	URLEncodeUpper(timestamp), URLEncodeUpper("AR"), URLEncodeUpper(shortCode), URLEncodeUpper("RedirectToCG"))
	fmt.Println("signature_url:  ", signature_url)
	signature := HmacSha256([]byte(signature_url), []byte(severConfig.ApiSecret))

	url := "http://ksg.kncee.com/MSG/v1.1/API/RedirectToCG?MSISDN=&applicationId=%s&countryId=%s&operatorId=%s" +
		"&cpId=%s&requestId=%s&apiKey=%s&signature=%s&timestamp=%s&lang=%s&shortcode=%s" +
		"&ipAddress=%s&lpUrl=%s&rurl=%s"
	url = fmt.Sprintf(url, severConfig.ApplicationId, severConfig.CountryId, severConfig.OperatorId, severConfig.CpId, track.TrackID, severConfig.ApiKey, signature,
		timestamp, "ar", severConfig.ShortCode, track.IP, "http://kg.foxseek.com/lp/"+strings.ToLower(strings.Replace(severConfig.ServiceID, "-", "/", 1)), severConfig.RUrl+strconv.FormatInt(track.TrackID, 10))
	fmt.Println("url: ", url)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)
	reqest.Header.Set("statusMessage", "http://thankyou.com")
	res, err := client.Do(reqest)

	fmt.Println(err)
	body, _ := ioutil.ReadAll(res.Body)

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

func HmacSha256(message, secret []byte) string {
	// secret := []byte("top-secret")
	// message := []byte("start1.99678678471198c6dec3-c5f0-4810-9490-e2b9f2e2d34ahttps://merch.at/cb?x=y")

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)

	// to lowercase hexits
	encode := hex.EncodeToString(hash.Sum(nil))

	return encode
}
