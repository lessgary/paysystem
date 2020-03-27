package controllers

import (
	"paySystem/hook"
	"paySystem/thread"

	"github.com/astaxie/beego"
)

type PayController struct {
	beego.Controller
}

/**
* 支付的请求
 */
func (this *PayController) PayCreate() {
	//定义需要输出的结果
	c_status := 100
	c_msg := "请求完成"
	paramsMap := map[string]string{}

	api_jump := 0
	form_param := map[string]string{}
	tpl_param := map[string]string{}
	//接收值
	merchant_id := this.GetString("merchant_id")
	pay_data := this.GetString("pay_data")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, paramsMap = hook.ChkInputAndMap(merchant_id, pay_data)

	if c_status == 200 {
		c_status, c_msg, tpl_param, api_jump, form_param = thread.PayCreate(merchant_id, paramsMap)
	}
	//1:扫码,2:调转,3:post提交
	if c_status == 200 {
		if api_jump == 1 {
			this.Data["img_url"] = tpl_param["img_url"]
			this.Data["amout"] = tpl_param["amout"]
			if tpl_param["pay_class"] == "2" && tpl_param["is_mobile"] == "1" {
				this.TplName = "mobile/alipay.tpl"
			} else if tpl_param["pay_class"] == "2" && tpl_param["is_mobile"] == "0" {
				this.TplName = "pc/alipay.tpl"
			} else if tpl_param["pay_class"] == "3" && tpl_param["is_mobile"] == "1" {
				this.TplName = "mobile/tenpay.tpl"
			} else if tpl_param["pay_class"] == "3" && tpl_param["is_mobile"] == "0" {
				this.TplName = "pc/tenpay.tpl"
			} else if tpl_param["pay_class"] == "4" && tpl_param["is_mobile"] == "1" {
				this.TplName = "mobile/wechatpay.tpl"
			} else if tpl_param["pay_class"] == "4" && tpl_param["is_mobile"] == "0" {
				this.TplName = "pc/wechatpay.tpl"
			} else if tpl_param["pay_class"] == "6" && tpl_param["is_mobile"] == "1" {
				this.TplName = "mobile/unionpay.tpl"
			} else if tpl_param["pay_class"] == "6" && tpl_param["is_mobile"] == "0" {
				this.TplName = "pc/unionpay.tpl"
			}
		} else if api_jump == 2 {
			this.Data["js_go_url"] = tpl_param["img_url"]
			this.TplName = "jump.tpl"
		} else if api_jump == 3 {
			//默认使用表单提交
			this.Data["Api_create_method"] = tpl_param["api_method"]
			this.Data["Api_create_url"] = tpl_param["pay_url"]
			this.Data["Api_form_param"] = form_param
			this.TplName = "put.tpl"
		} else if api_jump == 4 {
			this.Data["jump_url"] = tpl_param["img_url"]
			this.TplName = "auto_jump.tpl"
		} else if api_jump == 5 {
			this.Data["js_go_url"] = tpl_param["img_url"]
			this.TplName = "new_jump.tpl"
		}
	} else {
		this.Data["msg"] = c_msg
		this.TplName = "error.tpl"
	}
}

/**
* 代付接口
 */
func (this *PayController) PayFor() {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	paramsMap := map[string]string{}

	//接收值
	merchant_id := this.GetString("merchant_id")
	pay_data := this.GetString("pay_data")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, paramsMap = hook.ChkInputAndMap(merchant_id, pay_data)

	c_status, c_msg, d["result"] = thread.PayFor(merchant_id, paramsMap)

	//将数据装载到json返回值
	res := JsonOut{c_status, c_msg, d}

	//输出json数据
	this.Data["json"] = &res
	this.ServeJSON()
}

/**
* 支付支持的银行
 */
func (this *PayController) PayBank() {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"

	//接收值
	is_mobile := this.GetString("is_mobile")
	pay_id := this.GetString("pay_id")

	d["is_mobile"] = is_mobile
	d["pay_id"] = pay_id

	c_status, c_msg, d["bank_list"] = thread.PayBank(is_mobile, pay_id)

	//将数据装载到json返回值
	res := JsonOut{c_status, c_msg, d}

	//输出json数据
	this.Data["json"] = &res
	this.ServeJSON()
}
