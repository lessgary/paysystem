package thread

import (
	"paySystem/models"
)

/**
* 接收参数后，对数据做进一步处理
 */
func ThreadAdminCall(ordernumber string) (int, string) {
	re_status := 100
	re_msg := "订单不存在"
	p_list := models.Gorm.PayListOrder(ordernumber)
	if len(p_list.Id) < 1 {
		return re_status, re_msg
	}

	re_status = 200
	re_msg = "ok"
	Push(p_list)

	return re_status, re_msg
}

/**
* 接收参数后，对数据做进一步处理
 */
func ThreadAdminCallPayFor(ordernumber string) (int, string) {
	re_status := 100
	re_msg := "订单不存在"
	c_list := models.Gorm.CashListOrder(ordernumber)
	if len(c_list.Id) < 1 {
		return re_status, re_msg
	}

	re_status = 200
	re_msg = "ok"
	PushPayFor(c_list)

	return re_status, re_msg
}
