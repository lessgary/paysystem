package main

import (
	_ "paySystem/routers"
	"paySystem/thread"

	"github.com/astaxie/beego"
)

/**
* 有主动查询的api接口如下
1，leyin，乐盈支付
2,koudai,口袋支付
3,huanxun,环讯支付
*/
func main() {
	//防止丢单，循环执行
	go thread.ThreadPush()
	beego.Run()
}
