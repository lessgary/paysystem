package payapi

/**
* 定义需要传递到api的数据结构
 */
type PayData struct {
	Amount      string //订单金额
	OrderNumber string //订单编号
	PayClass    string //账户
	BankCode    string //选择的银行或第三方支付平台
	IsMobile    string //是否手机版 0是网页  1是手机
}

var date_format string = "2006-01-02 15:04:05"

var day_format string = "2006-01-02"

var month_format string = "2006-01"

var s_format string = "2006-01-02 00:00:00"

func init() {

}
