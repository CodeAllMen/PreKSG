package sp

import (
	"errors"
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/postback"
	"github.com/MobileCPX/PreBaseLib/util"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

// Postback 网盟信息
type Postback struct {
	ID           int     `orm:"pk;auto;column(id)" json:"-"`                            // 自增ID
	CreateTime   string  `orm:"column(create_time)" json:"-"`                           // 添加时间
	UpdateTime   string  `orm:"column(update_time)" json:"-"`                           // 更新时间
	DayCap       int     `orm:"column(day_cap)" json:"day_cap"`                         // 更新时间
	AffName      string  `orm:"column(aff_name);size(30)" json:"aff_name"`              // 网盟名称
	PostbackURL  string  `orm:"column(postback_url);size(300)" json:"postback_url"`     // postback URL
	PostbackRate int     `orm:"column(postback_rate);default(50)" json:"postback_rate"` // 回传概率
	Payout       float32 `orm:"column(Payout)" json:"payout"`                           // 转化单价
	PromoterName string  `orm:"column(promoter_name)" json:"promoter_name"`             // 外放人
	PromoterID   int     `orm:"column(promoter_id)" json:"promoter_id"`                 // 外放人
	CampID       int     `orm:"column(camp_id)" json:"camp_id"`                         // CampID
	OfferID      int     `orm:"column(offer_id)" json:"offer_id"`                       // offer_id
	postback.Postback
}

func (postback *Postback) TableName() string {
	return "postback"
}

// Insert  插入Postback
func (postBack *Postback) Insert() error {
	o := orm.NewOrm()
	postBack.CreateTime, _ = util.GetNowTime()
	_, err := o.Insert(postBack)
	if err != nil {
		logs.Error("Postback Insert 插入postback数据失败，", err.Error())
	}
	return err
}

// Update  更新Postback
func (postBack *Postback) Update() error {
	o := orm.NewOrm()
	postBack.UpdateTime, _ = util.GetNowTime()
	_, err := o.Update(postBack)
	if err != nil {
		logs.Error("Postback Update 插入postback数据成功")
	}
	return err
}

// CheckOfferIDIsExist 检查offer_id是否已经存在
func (postback *Postback) CheckOfferIDIsExist(offerID int) error {
	o := orm.NewOrm()
	err := o.QueryTable(postback.TableName()).Filter("offer_id", offerID).One(postback)
	if err != nil {
		logs.Error("Postback CheckOfferIDIsExist  ERROR, 检查offer_id是否已经存在 失败")
	}

	return err
}

// GetAffNameByOfferID
func GetAffNameByOfferID(offerID int) (*Postback, error) {
	o := orm.NewOrm()
	postback := new(Postback)
	err := o.QueryTable("postback").Filter("offer_id", offerID).One(postback)
	if err != nil {
		logs.Error("GetAffNameByOfferID 错误，offerID：", offerID, " ERROR: ", err.Error())
	}
	return postback, err
}

// GetPostbackInfoByOfferID 获取Postback信息 通过 OfferID
func (postback *Postback) GetPostbackInfoByOfferID(offerID int, affName, serviceName string) (*Postback, error) {
	o := orm.NewOrm()
	if offerID != 0 {
		err := o.QueryTable(postback.TableName()).Filter("offer_id", offerID).One(postback)
		if err != nil {
			logs.Error("用户订阅成功，但是没有找到此网盟 ", affName, "OfferID", offerID)
			util.BeegoEmail(serviceName, "没有找到此 "+"OfferID"+strconv.Itoa(offerID), "aff_name: "+affName+
				" postback回传失败", []string{})
		}
		return postback, err
	}
	return postback, errors.New("offerID为空")
}

// GetPostbackInfoByOfferID 获取Postback信息 通过 网盟名称
func (postback *Postback) GetPostbackInfoByAffName(affName, serviceName string) (*Postback, error) {
	o := orm.NewOrm()
	if affName != "" {
		err := o.QueryTable(postback.TableName()).Filter("aff_name", affName).OrderBy("-id").One(postback)
		if err != nil {
			logs.Error("用户订阅成功，但是没有找到此网盟 ", affName)
			util.BeegoEmail(serviceName, "没有找到此 "+affName+"信息", affName+" postback回传失败", []string{})
		}
		return postback, err
	}
	return postback, errors.New("网盟为空")
}

// CheckTodayPostbackStatus 检查今日订阅数和回传数，判断是否符合回传
func (postback *Postback) CheckTodayPostbackStatus(todaySubNum, todayPostbackNum int) (isPostback bool) {
	defer logs.Info("postbakck 状态 ", isPostback)
	if todaySubNum == 0 {
		isPostback = true
		return
	}
	currentRate := float32(todayPostbackNum) / float32(todaySubNum)
	if currentRate > float32(postback.PostbackRate)/float32(100) {
		isPostback = false
	} else {
		isPostback = true
	}
	return
}

func (postback *Postback) PostbackRequest(mo *mo.Mo) (isSuccess bool, code string) {
	postbackURL := postback.PostbackURL
	timestamp := time.Now().Unix()
	postbackURL = strings.Replace(postbackURL, "{click_id}", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "##clickid##", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "{pro_id}", mo.ProID, -1)
	postbackURL = strings.Replace(postbackURL, "{other}", mo.ProID, -1)
	postbackURL = strings.Replace(postbackURL, "{pub_id}", mo.PubID, -1)
	postbackURL = strings.Replace(postbackURL, "##pub_id##", mo.ClickID, -1)
	postbackURL = strings.Replace(postbackURL, "{operator}", mo.Operator, -1)
	postbackURL = strings.Replace(postbackURL, "{auto}", strconv.Itoa(int(timestamp)), -1)
	postbackURL = strings.Replace(postbackURL, "##auto_id##", strconv.Itoa(int(timestamp)), -1)
	postbackURL = strings.Replace(postbackURL, "{payout}", fmt.Sprintf("%f", postback.Payout), -1)

	postResult, err := httplib.Get(postbackURL).String()
	if err == nil {
		// postback 成功
		isSuccess = true
		logs.Info("postback URL: ", postbackURL, " CODE: ", code)
	} else {
		logs.Info("postback URL: ", postbackURL, " CODE: ", code)
		logs.Error("postback ERROR , msisdn : " + mo.Msisdn + " aff_name : " + mo.AffName + " error " + err.Error())
	}
	code = fmt.Sprintf("%.200s", postResult) // postback信息最长200
	return
}

func (postback *Postback) CheckOfferID(offerID int64) error {
	o := orm.NewOrm()
	return o.QueryTable(postback.TableName()).Filter("offer_id", offerID).One(postback)

}

func (postback *Postback) InsertPostback() error {
	o := orm.NewOrm()
	postback.CreateTime, _ = util.GetNowTime()
	_, err := o.Insert(postback)
	if err != nil {
		logs.Error("Postback InsertPostback ERROR:", err.Error(), postback)
	}
	return err
}
