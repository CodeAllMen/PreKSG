package models

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func InsertTrack(track *Track) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(track)
	return id, err
}

func SearchTrackById(id_str string) (*Track, error) {
	o := orm.NewOrm()
	id, err := strconv.Atoi(id_str)
	if err != nil {
		return nil, err
	}
	var track = new(Track)
	o.QueryTable("track").Filter("id", id).One(track)
	return track, nil
}

func UpdateTrackById(track *Track) {
	track.ClickStatus = "1"
	track.ClickTime = time.Now().Format("2006-01-02 15:04:05")
	o := orm.NewOrm()
	o.Update(track)
}

func Get_postback_url(camp_id string) (error, *Old_Postback) {
	var postback Old_Postback
	o := orm.NewOrm()
	err := o.QueryTable("postback").Filter("camp_id", camp_id).One(&postback)
	if err != nil {
		logs.Error("get postback url error: ", err)
	}
	return err, &postback
}

func PostbackRate(mo *MoStruct, rate int) bool {
	var status bool
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := int(source.Int63n(100))
	if num < rate {
		status = true
	} else {
		status = false
	}
	return status
}

func PostbackRequest(mo *MoStruct, postback *Old_Postback) (string, string) { // postback请求
	var urls, code string
	code = "400"
	url_model := postback.Url
	if url_model != "" {
		urls = strings.Replace(url_model, "##clickid##", mo.ClickId, -1)
		urls = strings.Replace(urls, "##proid##", mo.ProId, -1)
		urls = strings.Replace(urls, "##pubid##", mo.PubId, -1)
	}
	urls = strings.Replace(urls, "payout=0.35", "payout=0", -1)
	fmt.Println("===urls===:", urls)
	post_result, err := http.Get(urls)
	message, _ := ioutil.ReadAll(post_result.Body)
	post_result.Body.Close()
	if err == nil {
		code = strconv.Itoa(post_result.StatusCode)
	} else {
		code = err.Error()
	}
	return code, string(message)
}
