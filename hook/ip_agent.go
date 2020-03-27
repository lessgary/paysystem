/**
* 钩子程序，用于路由加载之前的各种权限判断
* @author:james
* @date:2015-09-08
 */
package hook

import (
	md "paySystem/models"
)

/**
* 判断IP白名单和接入商账户是否正确
* @author:james
* @date:2015-09-07
* @return int
	接入商不存在=-1
	接入商被锁定=0
	正常=1
	IP没有白名单=401
*/
func CheckIPAndAgent(ip, merchant_id string) (int, string) {
	if ip == "127.0.0.1" {
		return 200, "正常"
	}

	//第一步骤，检查接入商是否存在
	access_info := md.Gorm.AccessList(merchant_id)
	if access_info.Id < 1 {
		return 406, "商户不存在"
	}
	if access_info.Status != 1 {
		//接入商被锁定
		return 403, "商户被锁定"
	}

	chkIP := md.Gorm.AccessIPRedis(merchant_id, ip)
	if len(chkIP.Ip) > 0 {
		return 200, "正常"
	} else {
		return 403, "您的IP:" + ip + "不在访问白名单中"
	}
}
