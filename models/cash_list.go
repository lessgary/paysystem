package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) CreateCashList(cash_list CashList) error {
	res := d_o.GDb.Create(cash_list)
	return res.Error
}

func (d_o *DbOrm) CashListOrder(order_number string) CashList {
	var c_list CashList
	err := d_o.GDb.Model(&CashList{}).Where("order_number = ?", order_number).First(&c_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&CashList{}).Where("order_number = ?", order_number).First(&c_list)
	}

	return c_list
}

/**
*  根据订单号查询订单
 */
func (d_o *DbOrm) CashList(order_number string) CashList {
	var c_list CashList
	err := d_o.GDb.Model(&CashList{}).Where("web_ordernumber = ?", order_number).First(&c_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&CashList{}).Where("web_ordernumber = ?", order_number).First(&c_list)
	}

	return c_list
}

/**
*
 */
func (d_o *DbOrm) CashListRedis(order_number string) CashList {
	redis_key := "cash_list:" + order_number
	res := redisClient.Redis.StringRead(redis_key)
	var c_list CashList
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &c_list)
		return c_list
	}
	c_list = d_o.CashList(order_number)

	//判断是否存在用户名
	if len(c_list.Id) > 1 {
		res = common.Interface2Json(c_list)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_long_outime)
		}
	}
	return c_list
}

/**
* 分页查询条代付订单
 */
func (d_o *DbOrm) PageCashList(pageSize, offset int, s_time, e_time string, c_where map[string]interface{}) []CashList {
	var cash_lists []CashList

	d_o.GDb.Where("create_time>=? and create_time<=?", s_time, e_time).Where(c_where).Order("create_time").Limit(pageSize).Offset(offset).Find(&cash_lists)

	return cash_lists
}

func (d_o *DbOrm) CountCashList(s_time, e_time string, c_where map[string]interface{}) int {
	var c_count int

	d_o.GDb.Model(&CashList{}).Where("create_time>=? and create_time<=?", s_time, e_time).Where(c_where).Count(&c_count)
	return c_count
}

func (d_o *DbOrm) UpdatesCashList(cash_list CashList, c_data map[string]interface{}) error {
	res := d_o.GDb.Model(&cash_list).UpdateColumns(c_data)
	return res.Error
}

func (d_o *DbOrm) PushCashList(create_time string) []CashList {
	var cash_lists []CashList

	d_o.GDb.Where("push_status=9 and status=1 and push_num<15 and create_time>?", create_time).Order("create_time").Limit(100).Find(&cash_lists)

	return cash_lists
}
