/**
  create by yy on 2020/5/17
*/

package sp

import (
	"fmt"
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
		err        error
		list       []*sp.ChargeNotification
		total      = 0.0
		fee        float64
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
