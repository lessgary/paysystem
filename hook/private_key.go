package hook

import (
	"paySystem/models"
	"public/common"
)

var HookAES *common.AES

func getPrivateKey(merchant_id string) (string, string) {
	//第一步骤，读取密钥
	a_list := models.Gorm.AccessListRedis(merchant_id)
	//设置密钥
	return a_list.Private_key, a_list.Private_iv
}

func HookAesDecrypt(merchant_id, decode_str string) string {
	private_key, private_iv := getPrivateKey(merchant_id)
	mer_aes := common.SetAES(private_key, private_iv, "", 32)
	aes_res := mer_aes.AesDecryptString(decode_str)
	return aes_res
}

func HookAesEncrypt(merchant_id, decode_str string) string {
	private_key, private_iv := getPrivateKey(merchant_id)
	mer_aes := common.SetAES(private_key, private_iv, "", 32)
	aes_res := mer_aes.AesEncryptString(decode_str)
	return aes_res
}
