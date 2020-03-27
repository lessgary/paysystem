package controllers

import (
	"github.com/astaxie/beego"
)

func init() {
}

type MainController struct {
	beego.Controller
}

type CallController struct {
	beego.Controller
}

/*
* 定义一个json的返回值类型(首字母必须大写)
 */
type JsonOut struct {
	Status int
	Msg    string
	Data   map[string]interface{}
}
