package thread

import (
	"fmt"
	"paySystem/models"
	"public/common"
	"public/redisClient"
	"strconv"
	"time"
)

/**
* 支付
 */
func PayCreate(merchant_id string, p_map map[string]string) (int, string, map[string]string, int, map[string]string) {
	re_status := 100
	re_msg := "商户号错误"
	//1:扫码,2:调转,3:post提交
	api_jump := 3
	img_url := ""
	pay_url := ""
	api_method := ""
	form_param := map[string]string{}

	tpl_param := map[string]string{}
	//获取参数
	is_mobile := p_map["is_mobile"]
	amout := p_map["amount"]
	pay_id := p_map["pay_id"]
	web_ordernumber := p_map["order_number"]
	pay_class := p_map["pay_class"]
	bank_code := p_map["bank_code"]
	tpl_param["pay_class"] = pay_class
	tpl_param["is_mobile"] = is_mobile
	tpl_param["amout"] = amout
	tpl_param["img_url"] = ""

	if pay_class != "1" {
		bank_code = ""
	}

	redis_key := "PayCreate:" + web_ordernumber

	lock_res := redisClient.Redis.StringWrite(redis_key, web_ordernumber, 5)
	if lock_res < 1 {
		re_msg = "请求频繁"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	defer redisClient.Redis.KeyDel(redis_key)

	is_exist := models.Gorm.PayListRedis(web_ordernumber)
	if len(is_exist.Id) > 1 {
		re_msg = "订单已存在"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	acc_list := models.Gorm.AccessListRedis(merchant_id)
	if acc_list.Id < 1 {
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	if pay_class == "1" && len(bank_code) == 0 {
		re_msg = "bank_code不能为空"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	//转账金额必须大于0
	amoutInt, err := strconv.ParseFloat(amout, 64)

	if err != nil {
		re_msg = "支付金额格式错误"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	pay_conf := models.Gorm.PayConfigRedis(pay_id)
	if pay_conf.Id < 1 {
		re_msg = "支付渠道不存在"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	if pay_conf.Status != 1 {
		re_msg = "支付渠道维护中"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	pay_bank := models.Gorm.PayBankList(is_mobile, pay_conf.Cash_type_code, pay_class)
	if len(pay_bank.Bank_code) < 1 {
		re_msg = "该渠道未开通此支付方式"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}

	if len(bank_code) == 0 {
		bank_code = pay_bank.Bank_code
	}

	api_jump = pay_bank.Jump_status

	pay_id_int, _ := strconv.Atoi(pay_id)
	access_id, _ := strconv.Atoi(merchant_id)
	pay_class_int, _ := strconv.Atoi(pay_class)
	is_mobile_int, _ := strconv.Atoi(is_mobile)

	real_amout := amoutInt
	//生成存款订单
	list_id := models.GetKeyId()
	ctime := time.Now().Unix()
	ordernumber := pay_conf.Cash_type_code + list_id
	p_list := models.PayList{Id: list_id, Cash_type_code: pay_conf.Cash_type_code, Pay_id: pay_id_int, Access_id: access_id, Access_code: acc_list.Code, Amout: amoutInt, Real_amout: real_amout, Ctime: ctime, Ordernumber: ordernumber, Web_ordernumber: web_ordernumber, Pay_class: pay_class_int, Bank_code: bank_code, Web_type: is_mobile_int}
	c_res := models.Gorm.CreatePayList(p_list)
	if c_res != nil {
		re_msg = "订单插入失败，订单号重复"
		return re_status, re_msg, tpl_param, api_jump, form_param
	}
	re_msg = "支付请求失败"
	re_status, re_msg, api_method, pay_url, img_url, form_param = apiPayCreate(pay_id, pay_conf.Cash_type_code, amout, ordernumber, bank_code, is_mobile)
	tpl_param["img_url"] = img_url
	tpl_param["pay_url"] = pay_url
	tpl_param["api_method"] = api_method
	return re_status, re_msg, tpl_param, api_jump, form_param
}

/**
* 支付回调
 */
func PayCallBack(ordernumber, amout, sign, sign_str string, is_cent int) (int, string) {
	re_status := 100
	re_msg := "订单请求频繁"
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redisClient.Redis.StringWrite(ordernumber, ordernumber, 15)
	if order_lock == 0 {
		return re_status, re_msg
	}

	re_msg = "订单号不存在"
	p_list := models.Gorm.PayListOrder(ordernumber)
	if len(p_list.Id) < 1 {
		return re_status, re_msg
	}

	if p_list.Status != 0 {
		re_msg = "订单状态错误"
		return re_status, re_msg
	}

	back_amout, _ := strconv.ParseFloat(amout, 64)
	if is_cent == 1 {
		back_amout = back_amout / 100.00
	}
	//有的支付额度是不定的
	back_amout = back_amout + 1.00
	if back_amout < p_list.Amout {
		re_status = 101
		re_msg = "订单金额错误"
		return re_status, re_msg
	}
	pay_id := strconv.Itoa(p_list.Pay_id)
	re_status = apiPayCallBack(pay_id, p_list.Cash_type_code, sign, sign_str)
	if re_status != 200 {
		re_msg = "验签失败"
	}
	return re_status, re_msg
}

/**
* 代付
 */
func PayFor(merchart_id string, p_map map[string]string) (int, string, string) {
	re_status := 100

	re_msg := "订单号为空"
	re_result := "fail"

	web_orderunmber := p_map["order_number"]
	if len(web_orderunmber) < 1 {
		return re_status, re_msg, re_result
	}

	redis_key := "PayFor:" + web_orderunmber

	lock_res := redisClient.Redis.StringWrite(redis_key, web_orderunmber, 5)
	if lock_res < 1 {
		re_msg = "请求频繁"
		return re_status, re_msg, re_result
	}

	defer redisClient.Redis.KeyDel(redis_key)

	is_exist := models.Gorm.CashListRedis(web_orderunmber)
	if len(is_exist.Id) > 1 {
		re_msg = "订单已存在"
		return re_status, re_msg, re_result
	}

	acc_list := models.Gorm.AccessListRedis(merchart_id)
	if acc_list.Id < 1 {
		re_msg = "商户不存在"
		return re_status, re_msg, re_result
	}

	pay_id := p_map["pay_id"]

	pay_data := map[string]string{}

	p_conf := models.Gorm.PayConfigRedis(pay_id)
	if p_conf.Id < 1 {
		re_msg = "支付渠道不存在"
		return re_status, re_msg, re_result
	}

	amout_f, err := strconv.ParseFloat(p_map["amount"], 64)
	if err != nil {
		re_msg = "代付金额错误"
		return re_status, re_msg, re_result
	}
	real_amout := amout_f

	pay_cash := models.Gorm.PayCashType(p_conf.Cash_type_code)
	if len(pay_cash.Id) < 1 {
		re_msg = "支付不存在"
		return re_status, re_msg, re_result
	}

	fee_amout := 0.00

	bank_code := p_map["bank_code"]
	card_name := p_map["card_name"]
	card_number := p_map["card_number"]
	bank_branch := p_map["bank_branch"]
	bank_ext := p_map["bank_ext"]
	bank_phone := p_map["bank_phone"]
	bank_province := p_map["bank_province"]
	bank_city := p_map["bank_city"]
	bank_area := p_map["bank_area"]
	bank_cnaps := p_map["bank_cnaps"]
	bank_title := ""
	s_bank := models.Gorm.SysBankRedis(bank_code)
	if s_bank.Id < 1 {
		re_msg = "银行编码错误"
		return re_status, re_msg, re_result
	}
	bank_title = s_bank.Bank_title

	//创建提现订单
	c_id := models.GetKeyId()
	create_time := time.Now()
	pay_time := create_time.Format(format_date)
	finish_time := create_time.AddDate(0, 0, -1).Format(format_date)
	order_number := p_conf.Cash_type_code + c_id

	c_list := models.CashList{Id: c_id, Bank_code: bank_code, Bank_title: bank_title, Card_number: card_number, Card_name: card_name, Account: acc_list.Code, Create_time: pay_time, Web_ordernumber: web_orderunmber, Order_number: order_number, Cash_type: 1, Fee_amout: fee_amout, Finish_time: finish_time, Access_id: acc_list.Id, Pay_id: p_conf.Id, Amout: amout_f, Real_amout: real_amout}

	err = models.Gorm.CreateCashList(c_list)
	if err != nil {
		re_msg = "生成订单失败"
		return re_status, re_msg, re_result
	}

	if bank_ext == "" {
		bank_ext = common.GetRadomString(18, "number")
	}
	pay_data["amout"] = fmt.Sprintf("%.2f", real_amout) //代付金额
	pay_data["bank_code"] = bank_code                   //银行编码
	pay_data["card_name"] = card_name                   //持卡人姓名
	pay_data["card_number"] = card_number               //卡号
	pay_data["bank_branch"] = bank_branch               //支行信息
	pay_data["ordernumber"] = order_number              //代付单号
	pay_data["bank_ext"] = bank_ext                     //身份证
	pay_data["bank_phone"] = bank_phone                 //手机号码
	pay_data["bank_province"] = bank_province           //省份
	pay_data["bank_city"] = bank_city                   //城市
	pay_data["bank_area"] = bank_area                   //区
	pay_data["bank_cnaps"] = bank_cnaps                 //联行号
	pay_data["bank_title"] = bank_title                 //银行名称
	pay_data["pay_time"] = pay_time

	re_status, re_msg, re_result = apiPayFor(pay_id, p_conf.Cash_type_code, pay_data)

	if re_status != 200 || re_result == "processing" {
		return re_status, re_msg, re_result
	}

	status := "9"
	pay_order := order_number
	if len(re_msg) > 13 {
		pay_order = re_msg
	}
	if re_status == 200 && re_result == "success" {
		status = "1"
	}
	ThreadUpdateCash(order_number, pay_order, re_msg, status)
	return re_status, re_msg, re_result
}

/**
* 代付回调
 */
func PayForCallBack(ordernumber, amout, sign, sign_str string, is_cent int) (int, string) {
	re_status := 100
	re_msg := "订单请求频繁"
	//写入redis锁(10秒)，来排除重复处理的可能
	order_lock := redisClient.Redis.StringWrite(ordernumber, ordernumber, 15)
	if order_lock == 0 {
		return re_status, re_msg
	}

	re_msg = "订单号不存在"
	c_list := models.Gorm.CashListOrder(ordernumber)
	if len(c_list.Id) < 1 {
		return re_status, re_msg
	}

	if c_list.Status != 0 {
		re_msg = "订单状态错误"
		return re_status, re_msg
	}

	back_amout, _ := strconv.ParseFloat(amout, 64)
	if is_cent == 1 {
		back_amout = back_amout / 100.00
	}
	//有的支付额度是不定的
	back_amout = back_amout + 1.00

	pay_id := strconv.Itoa(c_list.Pay_id)

	p_conf := models.Gorm.PayConfigRedis(pay_id)
	if len(p_conf.Cash_type_code) < 0 {
		re_msg = "订单异常"
		return re_status, re_msg
	}
	re_msg = "验签失败"
	re_status = apiPayCallBack(pay_id, p_conf.Cash_type_code, sign, sign_str)
	return re_status, re_msg
}

/**
* 查询网银支付支持的银行列表
 */
func PayBank(is_mobile, pay_id string) (int, string, []map[string]string) {
	re_status := 100
	re_msg := "支付渠道ID错误"
	p_map := []map[string]string{}
	//根据pay_id查询cash_type_code
	p_conf := models.Gorm.PayConfigRedis(pay_id)
	if len(p_conf.Cash_type_code) < 0 {
		return re_status, re_msg, p_map
	}
	b_list := models.Gorm.BankList(is_mobile, p_conf.Cash_type_code)
	if len(b_list) < 1 {
		re_msg = "该支付渠道不支持网银"
		return re_status, re_msg, p_map
	}
	re_status = 200
	re_msg = "ok"
	for _, bank_info := range b_list {
		bank_map := map[string]string{}
		bank_map["bank_code"] = bank_info.Bank_code
		bank_map["bank_title"] = bank_info.Bank_title
		bank_map["style_code"] = bank_info.Style_code
		p_map = append(p_map[0:], bank_map)
	}
	return re_status, re_msg, p_map
}
