package hook

import (
	"encoding/json"
	"net/url"
	"public/common"
	"regexp"
	"strings"
)

/**
* sql注入过滤判断
* @param	str	需要做判断的字符串
* return	bool	如果包含注入关键词输出true，否则false
 */
func ChkDanger(str string) (int, string) {
	h_status := 200
	h_msg := "参数合法"
	str = strings.ToLower(str)
	danger_keys := `(?:')|(?:--)|(\b(select|update|delete|insert|trancate|char|into|substr|ascii|declare|exec|count|master|into|drop|execute|script)\b)`
	re, _ := regexp.Compile(danger_keys)
	if re.MatchString(str) {
		h_status = 401
		h_msg = "字符非法"
	}
	return h_status, h_msg
}

/**
* 解密出真实参数到公用变量paramsMap
* @param params	string	传递的参数值
* @return
	int		状态码   完成=200    输入错误=400
	string	文字描述
	map	解密后的参数结果集
*/
func ChkInputAndMap(merchant_id, params string) (int, string, map[string]string) {
	status := 200
	msg := "请求完成"
	res := map[string]string{}

	if len(params) < 1 {
		status = 400
		msg = "params参数不得为空"
		common.LogsWithFileName("pay_request_", "merchant_id->"+merchant_id+"\nparams->"+params)
		return status, msg, res
	}

	//数据解密

	en_param, err := url.QueryUnescape(params)
	en_param = strings.Replace(en_param, " ", "+", -1)
	if err != nil {
		status = 100
		msg = "URL转码失败"
		common.LogsWithFileName("pay_request_", "merchant_id->"+merchant_id+"\nparams->"+params)
		return status, msg, res
	}

	aes_str := HookAesDecrypt(merchant_id, en_param)

	common.LogsWithFileName("pay_request_", "merchant_id->"+merchant_id+"\nparams->"+params+"\r\nen_param->"+en_param+"\r\nparamsMapStr->"+aes_str)
	//防注入分析
	status, msg = ChkDanger(aes_str)
	if status != 200 {
		return status, msg, res
	}

	//将解密后的数据切割，并解析到公用参数map里面去
	err = json.Unmarshal([]byte(aes_str), &res)

	if err != nil {
		status = 100
		msg = "json解析失败"
		return status, msg, res
	}
	return status, msg, res
}
