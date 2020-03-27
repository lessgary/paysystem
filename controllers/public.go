package controllers

import (
	"image/png"
	"public/common"

	"github.com/astaxie/beego"
)

type PublicController struct {
	beego.Controller
}

/**
* 生成二维码
 */
func (this *PublicController) CreateQrCode() {
	//接收值
	qr_code := this.GetString("qr_code")

	c_img, _ := common.CreateQrCode(qr_code)

	this.Ctx.ResponseWriter.Header().Set("Content-Type", "image/png")
	png.Encode(this.Ctx.ResponseWriter, c_img)
}
