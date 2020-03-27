package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) PayConfig(pay_id string) PayThirdpayConfig {
	var pay_conf PayThirdpayConfig
	err := d_o.GDb.Model(&PayThirdpayConfig{}).Where("id = ?", pay_id).First(&pay_conf)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayThirdpayConfig{}).Where("id = ?", pay_id).First(&pay_conf)
	}

	return pay_conf
}

func (d_o *DbOrm) PayConfigRedis(pay_id string) PayThirdpayConfig {
	redis_key := "pay_thirdpay_config:" + pay_id
	res := redisClient.Redis.StringRead(redis_key)
	var pay_conf PayThirdpayConfig
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &pay_conf)
		return pay_conf
	}
	pay_conf = d_o.PayConfig(pay_id)

	//判断是否存在用户名
	if pay_conf.Id > 0 {
		res = common.Interface2Json(pay_conf)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_time)
		}
	}
	return pay_conf
}
