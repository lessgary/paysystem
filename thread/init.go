package thread

import (
	"paySystem/payapi"
)

//创建映射结构体
type PAYAPI interface {
	//初始化,加载配置
	Init(string)
	//创建支付订单
	CreatePay(*payapi.PayData) (int, string, string, string, string, map[string]string)
	//查询支付订单
	QueryPayOrder(map[string]string) (int, int, string)
	//回调验证
	CallBackPay(string, string) int
	//代付
	PayFor(map[string]string) (int, string, string)
	//查询代付订单
	QueryPayFor(map[string]string) (int, string, string)
	//商户余额查询
	PayBalance(map[string]string) (int, float64, string)
}

var format_date string = "2006-01-02 15:04:05"
var format_day string = "2006-01-02"
