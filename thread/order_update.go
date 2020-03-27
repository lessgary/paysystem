package thread

import (
	"paySystem/models"
	"public/redisClient"
	"time"
)

/**
* 支付订单完成
* 接收参数后，对数据做进一步处理
 */
func finished_pay(order_number, pay_ordernumber, status string) (int, string) {
	redis_key := "finished_pay_lock:" + order_number + "_" + pay_ordernumber
	re_status := 100
	re_msg := "请求频繁"
	lock_res := redisClient.Redis.StringWrite(redis_key, order_number, 5)
	if lock_res < 1 {
		return re_status, re_msg
	}
	defer redisClient.Redis.KeyDel(redis_key)

	p_list := models.Gorm.PayListOrder(order_number)
	if len(p_list.Id) < 1 {
		re_msg = "订单号错误"
		return re_status, re_msg
	}

	if p_list.Status != 0 {
		re_msg = "订单状态异常"
		return re_status, re_msg
	}

	p_data := map[string]interface{}{}
	p_data["status"] = status
	p_data["pay_ordernumber"] = pay_ordernumber
	p_data["finish_time"] = time.Now().Unix()
	err := models.Gorm.UpdatesPayList(p_list, p_data)

	re_msg = "订单更新失败"
	//将结果推送到网站去
	if err != nil {
		return re_status, err.Error()
	}
	re_status = 200
	re_msg = "ok"
	go Push(p_list)
	return re_status, re_msg
}

/**
* 提现或代付订单完成
 */
func finished_cash(order_number, pay_ordernumber, status string) (int, string) {
	re_status := 100
	re_msg := "订单不存在"
	c_list := models.Gorm.CashListOrder(order_number)

	if len(c_list.Id) < 1 {
		return re_status, re_msg
	}

	if c_list.Status != 0 {
		re_msg = "订单状态异常"
		return re_status, re_msg
	}

	c_data := map[string]interface{}{}
	c_data["status"] = status
	c_data["pay_number"] = pay_ordernumber
	c_data["finish_time"] = time.Now().Format(format_date)
	err := models.Gorm.UpdatesCashList(c_list, c_data)
	re_msg = "订单更新失败"

	if err != nil {
		return re_status, re_msg
	}
	re_status = 200
	re_msg = "ok"
	c_list.Status = 1
	go PushPayFor(c_list)
	return re_status, re_msg
}

/*
* 更新支付订单
**/
func ThreadUpdatePay(ordernumber, pay_order, status string) (int, string) {
	re_status, re_msg := finished_pay(ordernumber, pay_order, status)
	return re_status, re_msg
}

/*
* 更新提现订单
**/
func ThreadUpdateCash(order_number, pay_order, note, status string) (int, string) {
	redis_key := "finished_cash_lock:" + order_number + "_" + pay_order
	re_status := 100
	re_msg := "请求频繁"
	lock_res := redisClient.Redis.StringWrite(redis_key, order_number, 15)
	if lock_res < 1 {
		return re_status, re_msg
	}
	defer redisClient.Redis.KeyDel(redis_key)
	re_status, re_msg = finished_cash(order_number, pay_order, status)
	return re_status, re_msg
}
