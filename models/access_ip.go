package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) AccessIP(merchant_id, ip string) AccessIp {
	var access_ip AccessIp
	err := d_o.GDb.Model(&AccessIp{}).Where("access_id = ? and ip=?", merchant_id, ip).First(&access_ip)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&AccessIp{}).Where("access_id = ? and ip=?", merchant_id, ip).First(&access_ip)
	}

	return access_ip
}

func (d_o *DbOrm) AccessIPRedis(merchant_id, ip string) AccessIp {
	redis_key := "access_ip:" + merchant_id + "_" + ip
	res := redisClient.Redis.StringRead(redis_key)
	var access_ip AccessIp
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &access_ip)
		return access_ip
	}
	access_ip = d_o.AccessIP(merchant_id, ip)

	//判断是否存在用户名
	if len(access_ip.Ip) > 0 {
		res = common.Interface2Json(access_ip)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_long_outime)
		}
	}
	return access_ip
}
