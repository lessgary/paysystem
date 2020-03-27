package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) CreatePayList(pay_list PayList) error {
	res := d_o.GDb.Create(pay_list)
	return res.Error
}

/**
* 外部查询订单
 */
func (d_o *DbOrm) PayList(order_number string) PayList {
	var p_list PayList
	err := d_o.GDb.Model(&PayList{}).Where("web_ordernumber = ?", order_number).First(&p_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayList{}).Where("web_ordernumber = ?", order_number).First(&p_list)
	}

	return p_list
}

/**
* 外部查询订单
 */
func (d_o *DbOrm) PayListRedis(order_number string) PayList {
	redis_key := "pay_list:" + order_number
	res := redisClient.Redis.StringRead(redis_key)
	var p_list PayList
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &p_list)
		return p_list
	}
	p_list = d_o.PayList(order_number)

	//判断是否存在用户名
	if len(p_list.Id) > 1 {
		res = common.Interface2Json(p_list)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_long_outime)
		}
	}

	return p_list
}

/**
* 系统订单号查询
 */
func (d_o *DbOrm) PayListOrder(order_number string) PayList {
	var p_list PayList
	err := d_o.GDb.Model(&PayList{}).Where("ordernumber = ?", order_number).First(&p_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayList{}).Where("ordernumber = ?", order_number).First(&p_list)
	}

	return p_list
}

func (d_o *DbOrm) PushPayList(ctime int64) []PayList {
	var pay_lists []PayList

	d_o.GDb.Where("push_status=9 and status=1 and push_num<15 and ctime>?", ctime).Order("ctime").Limit(100).Find(&pay_lists)

	return pay_lists
}

/**
* 分页查询条支付订单
 */
func (d_o *DbOrm) PagePayList(pageSize, offset int, s_time, e_time int64, p_where map[string]interface{}) []PayList {
	var pay_lists []PayList

	d_o.GDb.Where("ctime>=? and ctime<=?", s_time, e_time).Where(p_where).Order("ctime").Limit(pageSize).Offset(offset).Find(&pay_lists)

	return pay_lists
}

func (d_o *DbOrm) CountPayList(s_time, e_time int64, p_where map[string]interface{}) int {
	var p_count int

	d_o.GDb.Model(&PayList{}).Where("ctime>=? and ctime<=?", s_time, e_time).Where(p_where).Count(&p_count)
	return p_count
}

func (d_o *DbOrm) UpdatePayList(p_list PayList, field, val string) error {
	res := d_o.GDb.Model(&p_list).UpdateColumn(field, val)
	return res.Error
}

func (d_o *DbOrm) UpdatesPayList(p_list PayList, p_data map[string]interface{}) error {
	res := d_o.GDb.Model(&p_list).UpdateColumns(p_data)
	return res.Error
}
