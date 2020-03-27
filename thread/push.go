package thread

import (
	"fmt"
	"paySystem/hook"
	"paySystem/models"
	"public/common"
	"public/redisClient"
	"strconv"
	"strings"
	"time"
)

func init() {
}

/**
* 订单推送的进程
* @param	access_code	接入商编号
* @param	web_ordernumber	网站订单号
* @param	ordernumber	系统订单号
* @param	amout			订单金额
* @param	status			订单的状态
 */
func Push(p_list models.PayList) {
	access_id := fmt.Sprintf("%d", p_list.Access_id)
	access_info := models.Gorm.AccessListRedis(access_id)
	if access_info.Id < 1 {
		return
	}

	if len(access_info.Push_url) < 1 {
		return
	}

	var result = fmt.Sprintf(`{"web_ordernumber":"%s","pay_id":"%d","pay_class":"%d","order_number":"%s","amount":"%.2f","status":"%d"}`, p_list.Web_ordernumber, p_list.Pay_id, p_list.Pay_class, p_list.Ordernumber, p_list.Amout, p_list.Status)
	//第三部，加密需要的数据
	aes_res := hook.HookAesEncrypt(access_id, result)
	aes_res = strings.Replace(aes_res, "+", "%2B", -1)
	strResult := fmt.Sprintf("merchant_id=%d&result=%s", p_list.Access_id, aes_res)
	//写入redis锁(15秒)，来排除并发的可能
	push_key := "push_" + p_list.Ordernumber
	rd_lock := redisClient.Redis.StringWrite(push_key, p_list.Ordernumber, 30)
	if rd_lock == 0 {
		//重复提单
		return
	} else {
		pushPool(strResult, access_info.Push_url, p_list)
	}
}

/**
* 结果推送，并记录结果
* @param	str	需要推送的内容
* @param	ordernumber	系统订单
 */
func pushPool(str, push_url string, p_list models.PayList) {
	//将推送的结果写入数据库
	p_data := map[string]interface{}{}
	push_status := 0
	msg, status := common.PushHttpPost(push_url, str)
	p_num := p_list.Push_num + 1
	common.LogsWithFileName("push_", "ordernumber->"+p_list.Ordernumber+"\nmsg->"+msg+" \n"+str+"\npush_url->"+push_url)
	if status == 200 && msg == "success" {
		push_status = 1
		p_data["push_status"] = push_status
		p_data["push_num"] = p_num
		models.Gorm.UpdatesPayList(p_list, p_data)
		return
	}

	push_status = 9

	//继续推送,连续10秒,推送12次，如果都失败了，放弃任务，并记录
	for i := 1; i < 10; i++ {
		p_num = p_num + 1
		time.Sleep(30 * time.Second)
		msg, status = common.PushHttpPost(push_url, str)
		common.LogsWithFileName("push_", "ordernumber->"+p_list.Ordernumber+"\ntimes->"+strconv.Itoa(i)+"\nrequest->"+str+"\nmsg->"+msg+"\npush_url->"+push_url)
		if status == 200 && msg == "success" {
			//成功后跳出循环
			push_status = 1
			break
		}
	}
	p_data["push_status"] = push_status
	p_data["push_num"] = p_num
	models.Gorm.UpdatesPayList(p_list, p_data)
}

/**
* 订单推送的进程
* @param	access_code	接入商编号
* @param	web_ordernumber	网站订单号
* @param	ordernumber	系统订单号
* @param	amout			订单金额
* @param	status			订单的状态
 */
func PushPayFor(c_list models.CashList) {
	access_id := fmt.Sprintf("%d", c_list.Access_id)
	access_info := models.Gorm.AccessListRedis(access_id)
	if access_info.Id < 1 {
		return
	}

	if len(access_info.Back_url) < 1 {
		return
	}
	var result = fmt.Sprintf(`{"web_ordernumber":"%s","pay_id":"%d","order_number":"%s","amount":"%.2f","real_amout":"%.2f","note":"%s","status":"%d"}`, c_list.Web_ordernumber, c_list.Pay_id, c_list.Order_number, c_list.Amout, c_list.Real_amout, c_list.Note, c_list.Status)
	//第三部，加密需要的数据
	aes_res := hook.HookAesEncrypt(access_id, result)
	aes_res = strings.Replace(aes_res, "+", "%2B", -1)
	strResult := fmt.Sprintf("merchant_id=%d&result=%s", c_list.Access_id, aes_res)
	//写入redis锁(15秒)，来排除并发的可能
	push_key := "push_" + c_list.Order_number
	rd_lock := redisClient.Redis.StringWrite(push_key, c_list.Order_number, 30)
	if rd_lock == 0 {
		//重复提单
		return
	} else {
		pushCashOrder(strResult, access_info.Back_url, c_list)
	}
}

/**
* 结果推送，并记录结果
* @param	str	需要推送的内容
* @param	ordernumber	系统订单
 */
func pushCashOrder(str, push_url string, c_list models.CashList) {
	//将推送的结果写入数据库
	c_data := map[string]interface{}{}
	push_status := 0
	msg, status := common.PushHttpPost(push_url, str)
	common.LogsWithFileName("push_payfor_", "ordernumber->"+c_list.Order_number+"\n push ok \n"+str+"\npush_url->"+push_url)
	p_num := c_list.Push_num + 1
	if status == 200 && msg == "success" {
		push_status = 1
		c_data["push_status"] = push_status
		c_data["push_num"] = p_num
		models.Gorm.UpdatesCashList(c_list, c_data)
		return
	}

	push_status = 9

	//继续推送,连续10秒,推送12次，如果都失败了，放弃任务，并记录
	for i := 1; i < 10; i++ {
		p_num = p_num + 1
		time.Sleep(30 * time.Second)
		msg, status = common.PushHttpPost(push_url, str)
		common.LogsWithFileName("push_payfor_", "ordernumber->"+c_list.Order_number+"\ntimes->"+strconv.Itoa(i)+"\nrequest->"+str+"\nmsg->"+msg+"\npush_url->"+push_url)
		if status == 200 && msg == "success" {
			//成功后跳出循环
			push_status = 1
			break
		}
	}

	c_data["push_status"] = push_status
	c_data["push_num"] = p_num
	models.Gorm.UpdatesCashList(c_list, c_data)
}
