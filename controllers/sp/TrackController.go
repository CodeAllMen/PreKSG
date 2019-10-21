package sp

import (
	"github.com/MobileCPX/PreBaseLib/splib/mo"
	"github.com/MobileCPX/PreBaseLib/splib/tracking"
	"github.com/MobileCPX/PreKSG/models/sp"
	"strconv"
)

var DayLimitSubNum = make(map[string]int)

type TrackingController struct {
	BaseController
}

func (c *TrackingController) InsertAffClick() {
	track := new(sp.AffTrack)
	returnStr := ""
	defer func() {
		if returnStr == "false" {
			if returnStr == "false" {
				track.Update()
			}
		}
	}()

	reqTrack := new(tracking.Track)
	reqTrack, err := reqTrack.BodyToTrack(c.Ctx.Request.Body)

	if err != nil {
		c.StringResult("false")
	}

	track.Track = *reqTrack
	serviceConfig, _ := c.serviceCofig(track.ServiceID)
	// 添加判断是否可以订阅条件
	moT := new(mo.Mo)
	SubNum, _ := moT.GetServiceTodaySubNum(track.ServiceID) // GBB 订阅数

	LimitSubNum := serviceConfig.LimitSubNum // GBB 限制订阅数

	trackID, err := track.Insert()
	if err != nil {
		c.StringResult("false")
	}
	if SubNum >= LimitSubNum && track.AffName != "Test" {
		returnStr = "false"
	}

	if returnStr != "false" {
		returnStr = strconv.Itoa(int(trackID)) // 返回自增ID
	}

	c.StringResult(returnStr)

}
