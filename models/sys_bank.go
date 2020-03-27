package models

import (
	"encoding/json"
	"public/common"
	"public/redisClient"
)

func (d_o *DbOrm) SysBank(bank_code string) SysBank {
	var s_bank SysBank
	err := d_o.GDb.Model(&SysBank{}).Where("bank_code = ?", bank_code).First(&s_bank)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&SysBank{}).Where("bank_code = ?", bank_code).First(&s_bank)
	}

	return s_bank
}

func (d_o *DbOrm) SysBankRedis(bank_code string) SysBank {
	redis_key := "sys_bank:" + bank_code
	res := redisClient.Redis.StringRead(redis_key)
	var s_bank SysBank
	if len(res) > 1 {
		json.Unmarshal([]byte(res), &s_bank)
		return s_bank
	}
	s_bank = d_o.SysBank(bank_code)

	//判断是否存在用户名
	if s_bank.Id > 1 {
		res = common.Interface2Json(s_bank)
		if len(res) > 0 {
			redisClient.Redis.StringWrite(redis_key, res, M_redis_time)
		}
	}
	return s_bank
}
