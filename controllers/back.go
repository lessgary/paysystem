package controllers

import (
	"github.com/astaxie/beego"
)

type BackController struct {
	beego.Controller
}

func (this *BackController) PaySuccess() {
	this.TplName = "success.tpl"
}
