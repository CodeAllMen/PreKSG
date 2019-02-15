package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"time"
)

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
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

func URLEncodeUpper(str string) string {
	url_encode := url.QueryEscape(str)
	url_encode = strings.Replace(url_encode, "%2B", "%2b", -1)
	url_encode = strings.Replace(url_encode, "%2F", "%2f", -1)
	return url_encode
}
