/**
  create by yy on 2020/5/17
*/

package sp

import (
	"fmt"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreKSG/libs"
	"github.com/MobileCPX/PreKSG/models/sp"
	"github.com/astaxie/beego"
	"strconv"
)

type CountController struct {
	beego.Controller
}

func (c *CountController) Count() {

	var (
		err   error
		list  []*sp.ChargeNotification
		total = 0.0
		fee   float64
	)

	startTime := c.GetString("start_time")
	endTime := c.GetString("end_time")

	chargeModel := new(sp.ChargeNotification)

	if list, err = chargeModel.GetChargeList(startTime, endTime); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	// 遍历 数据，并且进行累加
	for _, data := range list {
		// 扣费成功的才进行计算
		if fee, err = strconv.ParseFloat(data.Rate, 64); err != nil {
			err = libs.NewReportError(err)
			fmt.Println(err)
		} else {
			total = total + fee
		}
	}

	result := "%v  到  %v 的总扣费为：%v, 扣费成功数为：%v"
	result = fmt.Sprintf(result, startTime, endTime, total, len(list))

	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CountController) CountSub() {

	var (
		err   error
		list  []*sp.ChargeNotification
		total = 0.0
		fee   float64
	)

	startTime := c.GetString("start_time")
	endTime := c.GetString("end_time")

	startTime2 := c.GetString("start_time_2")
	endTime2 := c.GetString("end_time_2")

	systemMark := c.GetString("sm")

	chargeModel := new(sp.ChargeNotification)

	if list, err = chargeModel.GetChargeListSub(startTime, endTime, startTime2, endTime2, systemMark); err != nil {
		err = libs.NewReportError(err)
		fmt.Println(err)
	}

	// 遍历 数据，并且进行累加
	for _, data := range list {
		// 扣费成功的才进行计算
		if fee, err = strconv.ParseFloat(data.Rate, 64); err != nil {
			err = libs.NewReportError(err)
			fmt.Println(err)
		} else {
			total = total + fee
		}
	}

	result := "%v  到  %v 的总扣费为：%v, 扣费成功数为：%v"
	result = fmt.Sprintf(result, startTime, endTime, total, len(list))

	c.Data["json"] = result
	c.ServeJSON()

}

func (c *CountController) DivideSystem() {
	// 首先 获取 所有数据
	// 然后 进行遍历
	// 通过 charge_notification表的 transaction_id 跟aff_track 表关联获取 service_id

	var (
		err        error
		chargeList []*sp.ChargeNotification
		affTrack   *sp.AffTrack
		result     string
	)

	chargeNotification := new(sp.ChargeNotification)

	if chargeList, err = chargeNotification.GetList(); err != nil {
		err = libs.NewReportError(err)
		result = fmt.Sprintf("%v", err)
		c.Data["json"] = result
		c.ServeJSON()
	}

	for _, charge := range chargeList {

		affTrack = new(sp.AffTrack)

		if charge.TransactionId == "" {
			continue
		}

		if affTrack.TrackID, err = strconv.ParseInt(charge.TransactionId, 10, 64); err != nil {
			err = libs.NewReportError(err)
			result = fmt.Sprintf("%v", err)
			continue
		}

		if err = affTrack.GetOne(tracking.ByTrackID); err != nil {
			continue
		}

		if affTrack.ServiceID == "BB-NEW-ET" ||
			affTrack.ServiceID == "GF-NEW-ET" ||
			affTrack.ServiceID == "EB-NEW-ET" ||
			affTrack.ServiceID == "MA-NEW-ET" ||
			affTrack.ServiceID == "POM-NEW-ET" {
			charge.SystemMark = 2
		} else {
			charge.SystemMark = 1
		}

		if err = charge.UpdateSystemMark(); err != nil {
			err =libs.NewReportError(err)
			result = fmt.Sprintf("%v", err)
			break
		}

	}

	c.Data["json"] = result

	c.ServeJSON()
}
