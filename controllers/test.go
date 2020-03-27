package controllers

import (
	"image/png"
	"paySystem/hook"
	"paySystem/thread"
	"public/common"

	"github.com/astaxie/beego"
)

type TestController struct {
	beego.Controller
}

var TestAES *common.AES

func (c *TestController) Test() {
	c.TplName = "index.tpl"
}

/**
* 测试加密
 */
func (this *TestController) TestEncode() {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	private_key := this.GetString("private_key")
	params := this.GetString("params")

	TestAES = common.SetAES(private_key, "", "pkcs5", 32)
	d["result"] = TestAES.AesEncryptString(params)
	//将数据装载到json返回值
	res := JsonOut{c_status, c_msg, d}

	//输出json数据
	this.Data["json"] = &res
	this.ServeJSON()
}

/**
* 测试支付
 */
func (this *TestController) TestPay() {
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	paramsMap := map[string]string{}

	//接收值
	merchant_id := this.GetString("merchant_id")
	params := this.GetString("params")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, paramsMap = hook.ChkInputAndMap(merchant_id, params)

	if c_status == 200 {
		c_status, c_msg, d["tpl_param"], d["api_jump"], d["form_param"] = thread.PayCreate(merchant_id, paramsMap)
	}

	//将数据装载到json返回值
	res := JsonOut{c_status, c_msg, d}

	//输出json数据
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *TestController) TestPayFor() {
	//定义需要输出的结果
	d := map[string]interface{}{}
	c_status := 100
	c_msg := "请求完成"
	paramsMap := map[string]string{}

	//接收值
	merchant_id := this.GetString("merchant_id")
	params := this.GetString("params")
	/////////////////获得输入的值/////////////////
	c_status, c_msg, paramsMap = hook.ChkInputAndMap(merchant_id, params)

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
func (this *TestController) TestPayBank() {
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

/**
* 支付支持的银行
 */
func (this *TestController) TestCreateQrCode() {

	//接收值
	qr_code := this.GetString("qr_code")

	c_img, _ := common.CreateQrCode(qr_code)

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
	png.Encode(this.Ctx.ResponseWriter, c_img)
}

func (this *TestController) TestEncodeCreateQrCode() {
	//接收值
	qr_code := this.GetString("qr_code")
	qr_key := this.GetString("qr_key")

	TestAES = common.SetAES(qr_key, "", "", 32)

	de_qr := TestAES.AesDecryptString(qr_code)

	c_img, _ := common.CreateQrCode(de_qr)

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
	png.Encode(this.Ctx.ResponseWriter, c_img)
}

func (this *TestController) TestQrCode() {
	//接收值
	img_url := this.GetString("qr_code")

	this.Data["img_url"] = img_url
	this.TplName = "test.tpl"
}

func (this *TestController) TestAliJs() {
	money := this.GetString("money")
	account := this.GetString("account")
	order_id := this.GetString("order_id")
	this.Data["money"] = money
	this.Data["account"] = account
	this.Data["order_id"] = order_id
	this.TplName = "alipay_js.tpl"
}

func (this *TestController) TestJs() {
	money := this.GetString("money")
	account := this.GetString("account")
	order_id := this.GetString("order_id")
	this.Data["money"] = money
	this.Data["account"] = account
	this.Data["order_id"] = order_id
	this.TplName = "test.tpl"
}
