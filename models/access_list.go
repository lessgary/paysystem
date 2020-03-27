package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) AccessList(access_id string) AccessList {
	var access_list AccessList
	err := d_o.GDb.Model(&AccessList{}).Where("id = ?", access_id).First(&access_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&AccessList{}).Where("id = ?", access_id).First(&access_list)
	}

	return access_list
}

func (d_o *DbOrm) AccessListRedis(access_id string) AccessList {
	redis_key := "access_list:" + access_id
	res := redisClient.Redis.StringRead(redis_key)
	var access_list AccessList
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &access_list)
		return access_list
	}
	access_list = d_o.AccessList(access_id)

	//判断是否存在用户名
	if access_list.Id > 1 {
		res = common.Interface2Json(access_list)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_long_outime)
		}
	}
	return access_list
}
