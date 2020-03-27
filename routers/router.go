package routers

import (
	"paySystem/controllers"
	"paySystem/hook"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

/*
* 定义一个json的返回值类型(首字母必须大写)
 */
type JsonOut struct {
	Status int
	Msg    string
	Data   map[string]interface{}
}

/**
* 在路由前面挂一个钩子，去判断指定规则的请求是否合法
* @author:james
* @date:2015-09-08
 */
var FilterUser = func(ctx *context.Context) {
	merchant_id := strings.Trim(ctx.Request.PostFormValue("merchant_id"), " ")

	c_ip := ctx.Request.Header.Get("X-Forwarded-For")
	if len(c_ip) < 7 {
		c_ip = ctx.Input.IP()
	}

	//判断参数是否合法
	chk, msg := hook.ChkDanger(merchant_id)
	if chk == 200 {
		//判断接入时和ip白名单是否可用
		chk, msg = hook.CheckIPAndAgent(c_ip, merchant_id)
	}

	if chk != 200 {
		//将数据装载到json返回值
		d := map[string]interface{}{}
		res := JsonOut{chk, msg, d}
		ctx.Output.JSON(res, true, true)
	}
}

func init() {
	beego.Router("/", &controllers.MainController{})

	//###################支付接口######################/
	//交易统计查询接口

	//##################支付回调########################/

	//##################代付回调########################

	//###################对外接口######################/

	beego.InsertFilter("/pay/*", beego.BeforeRouter, FilterUser)

	beego.Router("/pay/create.do", &controllers.PayController{}, "post:PayCreate")

	beego.Router("/pay/pay_for.do", &controllers.PayController{}, "post:PayFor")

	beego.Router("/pay/pay_bank.do", &controllers.PayController{}, "post:PayBank")

	//付款成功的页面
	beego.Router("/back/pay_success.do", &controllers.BackController{}, "*:PaySuccess")

	//###################静态文件路由######################/
	beego.SetStaticPath("/static", "./views")
	beego.SetStaticPath("/download", "./download")

	//####################内部使用的接口####################/
	//生成二维码
	beego.Router("/public/create_qr_code.do", &controllers.PublicController{}, "*:CreateQrCode")
	//###################测试路由######################/
	beego.Router("/test/test.do", &controllers.TestController{}, "*:Test")
	beego.Router("/test/test_encode.do", &controllers.TestController{}, "*:TestEncode")
	beego.Router("/test/pay.do", &controllers.TestController{}, "*:TestPay")
	beego.Router("/test/pay_query.do", &controllers.TestController{}, "*:TestPayQuery")
	beego.Router("/test/pay_for.do", &controllers.TestController{}, "*:TestPayFor")
	beego.Router("/test/pay_for_query.do", &controllers.TestController{}, "*:TestPayForQuery")
	beego.Router("/test/pay_bank.do", &controllers.TestController{}, "*:TestPayBank")
	beego.Router("/test/create_qr_code.do", &controllers.TestController{}, "*:TestCreateQrCode")
	beego.Router("/test/qr_code.do", &controllers.TestController{}, "*:TestQrCode")
	beego.Router("/test/alipay_js.do", &controllers.TestController{}, "*:TestAliJs")
	beego.Router("/test/test_js.do", &controllers.TestController{}, "*:TestJs")
	beego.Router("/test/decode_qr_code.do", &controllers.TestController{}, "*:TestEncodeCreateQrCode")
}
