package models

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

type AffMoMtCilck struct {
	AffName  string
	Aff_data []PubData
}

type PubData struct {
	Pubname  string
	Ser_list []ProData
}

type ProData struct {
	Servername    string
	Total_num     int
	Active_num    int
	Unsub_num     int
	Click_num     int
	SuccessMT_Num int
	FailtMT_Num   int
	PostNum       int
	Churn_rate    string
}

type MoMtClickData struct {
	AffName     string
	PubId       string
	ProId       string
	SubNum      int
	SuccessMt   int
	MtFailed    int
	UnsubNum    int
	PostbackNum int
	ClickNum    int
}

//渠道转化
func GetAffdDate(startTime, endTime, keyword, operator, aff_name string) (error, []AffMoMtCilck) {
	o := orm.NewOrm()
	o.Using("default")
	var data []AffMoMtCilck
	var affData []MoMtClickData
	var total MoMtClickData

	filter_sql_mo := ""
	if operator != "All" {
		filter_sql_mo = filter_sql_mo + fmt.Sprintf(" and operator = '%s'", operator)
	}
	if aff_name != "All" {
		filter_sql_mo = filter_sql_mo + fmt.Sprintf(" and aff_name = '%s'", aff_name)
	}
	// filter_sql_dn := strings.Replace(filter_sql_mo, "text", "keyword", -1)

	mtGroupByAffSql := fmt.Sprintf("select n.aff_name,n.pub_id, count(case when n.dn_status = 0 then 1 else null end) as "+
		"MtFailed,count(case when n.dn_status = 1 then 1 else null end) as SuccessMt from dn_struct n where"+
		" n.sub_id in (select id from mo_struct where left(subtime,10)>='%s' and left(subtime,10)<'%s' %s) "+
		"group by n.aff_name,n.pub_id", startTime, endTime, filter_sql_mo)

	clickGroupByAffSql := fmt.Sprintf("select aff_name,pub_id,count(1) as click_num from track where left(time,10)>='%s'"+
		" and left(time,10)<'%s' group by aff_name,pub_id", startTime, endTime)

	moGroupByAffSql := fmt.Sprintf("select count(1) as SubNum,sum(postback_status) as PostbackNum, aff_name, pub_id, count(case "+
		"when left(unsubtime,10) < '%s' and unsubtime<>'' then 1 else null end) as UnsubNum from mo_struct where left(subtime,10)>='%s'"+
		" and left(subtime,10)<'%s' %s group by aff_name,pub_id", endTime, startTime, endTime, filter_sql_mo)

	totalSql := fmt.Sprintf("select mo.aff_name as Aff_name ,mo.pub_id,mo.SubNum as sub_num,mo.PostbackNum as Postback_num,"+
		"mo.UnsubNum as unsub_num ,mt.MtFailed as Mt_failed,mt.SuccessMt as Success_mt,"+
		"click.click_num from (%s) as mo left join (%s) as mt on mo.aff_name=mt.aff_name and mo.pub_id=mt.pub_id, (%s) as click where "+
		"mo.aff_name=click.aff_name and mo.pub_id=click.pub_id order by aff_name,pub_id", moGroupByAffSql, mtGroupByAffSql, clickGroupByAffSql)

	// totalSql := fmt.Sprintf("select mo.aff_name as Aff_name ,mo.pub_id,mo.SubNum as sub_num,mo.PostbackNum as Postback_num,"+
	// 	"mo.UnsubNum as unsub_num ,"+
	// 	"click.click_num from (%s) as mo, (%s) as click where "+
	// 	"mo.aff_name=click.aff_name and mo.pub_id=click.pub_id order by aff_name,pub_id", moGroupByAffSql, clickGroupByAffSql)

	_, err := o.Raw(totalSql).QueryRows(&affData)

	total.AffName = "Total"
	total.PubId = "Total"
	for _, subCharge := range affData {
		total.ClickNum += subCharge.ClickNum
		total.SubNum += subCharge.SubNum
		total.UnsubNum += subCharge.UnsubNum
		total.PostbackNum += subCharge.PostbackNum
		total.SuccessMt += subCharge.SuccessMt
		total.MtFailed += subCharge.MtFailed
	}
	affData = append(affData, total)
	copy(affData[1:], affData[0:len(affData)-1])
	affData[0] = total
	var affName, PubName string
	for i, subData := range affData {
		var oneData AffMoMtCilck
		var pubData PubData
		var serviceData ProData
		if affData[i].AffName != affName {
			affName = affData[i].AffName
			PubName = affData[i].PubId
			oneData.AffName = affName
			pubData.Pubname = affData[i].PubId
			serviceData.Servername = affData[i].ProId
			pubData.Ser_list = append(pubData.Ser_list, GetserviceDataList(subData))
			oneData.Aff_data = append(oneData.Aff_data, pubData)
			data = append(data, oneData)
		} else {
			if affData[i].PubId != PubName {
				PubName = affData[i].PubId
				oneData.AffName = affName
				pubData.Pubname = affData[i].PubId
				serviceData.Servername = affData[i].ProId
				pubData.Ser_list = append(pubData.Ser_list, GetserviceDataList(subData))
				for i, _ := range data {
					if data[i].AffName == affName {
						data[i].Aff_data = append(data[i].Aff_data, pubData)
						break
					}
				}
			} else {
				PubName = affData[i].PubId
				oneData.AffName = affName
				pubData.Pubname = affData[i].PubId
				serviceData.Servername = affData[i].ProId
				pubData.Ser_list = append(pubData.Ser_list, GetserviceDataList(subData))
				for i, _ := range data {
					if data[i].AffName == affName {
						for j, _ := range data[i].Aff_data {
							if data[i].Aff_data[j].Pubname == PubName {
								data[i].Aff_data[j].Ser_list = append(data[i].Aff_data[j].Ser_list, GetserviceDataList(subData))
								break
							}
						}
						break
					}
				}
			}
		}
	}
	return err, data
}

func GetserviceDataList(subData MoMtClickData) ProData {
	var service ProData
	service.Servername = subData.ProId
	service.PostNum = subData.PostbackNum
	service.Click_num = subData.ClickNum
	service.Total_num = subData.SubNum
	service.FailtMT_Num = subData.MtFailed
	service.SuccessMT_Num = subData.SuccessMt
	service.Unsub_num = subData.UnsubNum
	churn_rate := float32(service.Unsub_num) / float32(service.Total_num) * 100
	service.Churn_rate = fmt.Sprintf("%.2f", churn_rate) + "%"
	return service
}
