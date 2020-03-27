package thread

import (
	"paySystem/payapi"
)

/**
* 访问支付接口，返回结果
* @param	cash_type_code	//
* @param	partner_id		//平台商家id
* @param	partnet_key	//平台商家key
* @param	pay_loginname	//平台商家操作员
* @param	pay_url			//平台商家支付接口
* @param	pay_domain			//平台对应的域名
* @param	pay_amout		//支付金额
* @param	pay_ordernumber	//本系统入的订单号
* @param	pay_account	//网站的会员
* @param	bank_code		//选择的银行编码
* @param	ext		//扩展字段


* return	int		//返回状态,200=完成
* retutn	string	//返回的是get还是post
* return	string	//需要get提交的参数
* return	string	//需要form提交的参数
* return	string	//需要提交的url，或者二维码的地址
 */
func apiPayCreate(pay_id, cash_type_code, pay_amout, pay_ordernumber, bank_code, is_mobile string) (int, string, string, string, string, map[string]string) {

	//定义初始值
	api_status := 100
	re_msg := ""
	form_param := map[string]string{}
	pay_url := ""
	img_url := ""
	api_method := ""

	//对数据赋值
	create_PayData := &payapi.PayData{Amount: pay_amout, OrderNumber: pay_ordernumber, BankCode: bank_code, IsMobile: is_mobile}

	api_Pay := apiPayInit(pay_id, cash_type_code)

	//初始化api的配置
	if api_Pay != nil {
		api_status, re_msg, api_method, pay_url, img_url, form_param = api_Pay.CreatePay(create_PayData)
	}

	return api_status, re_msg, api_method, pay_url, img_url, form_param
}

/**
* 支付查询
 */
func apiQueryPay(pay_id, cash_code string, pay_data map[string]string) (int, int, string) {
	//定义初始值
	api_status := 100
	pay_status := 100
	api_msg := "查询失败"

	api_Pay := apiPayInit(pay_id, cash_code)

	//初始化api的配置
	if api_Pay != nil {
		api_status, pay_status, api_msg = api_Pay.QueryPayOrder(pay_data)
	}
	return api_status, pay_status, api_msg
}

/**
* 判断回调，返回签名是否正确
* @param	data	map[string]string	接口回调参数

* return	int		//返回状态,200=完成
 */
func apiPayCallBack(pay_id, cash_code, sign, sign_str string) int {
	//定义初始值
	api_status := 100

	api_Pay := apiPayInit(pay_id, cash_code)

	//初始化api的配置
	if api_Pay != nil {
		api_status = api_Pay.CallBackPay(sign, sign_str)
	}
	return api_status
}

/**
* 代付
 */
func apiPayFor(pay_id, cash_code string, pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_result := "fail"
	api_msg := "代付失败"

	api_Pay := apiPayInit(pay_id, cash_code)

	//初始化api的配置
	if api_Pay != nil {
		api_status, api_msg, pay_result = api_Pay.PayFor(pay_data)
	}
	return api_status, api_msg, pay_result
}

/**
* 代付查询
 */
func apiQueryPayFor(pay_id, cash_code string, pay_data map[string]string) (int, string, string) {
	//定义初始值
	api_status := 100
	pay_status := "fail"
	api_msg := "渠道编码错误"

	api_Pay := apiPayInit(pay_id, cash_code)

	//初始化api的配置
	if api_Pay != nil {
		api_status, pay_status, api_msg = api_Pay.QueryPayFor(pay_data)
	}
	return api_status, pay_status, api_msg
}

/**
* 初始化支付API
 */
func apiPayInit(pay_id, cash_code string) PAYAPI {
	//对结构体实例化
	var api_Pay PAYAPI

	switch cash_code {
	// case "cfpay":
	// 	api_Pay = new(payapi.CFPAY)
	// 	api_Pay.Init(pay_id)

	}
	return api_Pay
}
