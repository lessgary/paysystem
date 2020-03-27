<!doctype html>
<html>
  	<head>
    	<title>第三方支付中间件</title>
    	<meta http-equiv="content-type" content="text/html; charset=utf-8">
		
	</head>
<style>
table,tr,td{
	border:1px #000000 solid;
}
</style>
<body>
<p><b>第三方支付中间件，api调用说明</b></p>
<p>
ajax返回值说明:
</p>
<p>
{
  "Status": 200,
  "Msg": "请求完成",
  "Data": {
    "account": "rrrrr"
  }
}
<p>

<p>加密校验，验证您的加密结果是否更接口的一样<br />
/test/test_encode.do</p>
<table>
<form method="post" action="/test/test_encode.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>private_key</td>
<td><input name="private_key" type="text"></td>
<td>密钥</td>
</tr>
<tr>
<td>params</td>
<td><textarea name="params"></textarea></td>
<td>需要加密的参数</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>测试支付<br />
/test/pay.do</p>
<table>
<form method="post" action="/test/pay.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>merchant_id</td>
<td><input name="merchant_id" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><textarea name="params"></textarea></td>
<td>
	<p>需要加密的参数，格式和post表头提交的一样</p>
</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>测试支付订单查询<br />
/test/pay_query.do</p>
<table>
<form method="post" action="/test/pay_query.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>merchant_id</td>
<td><input name="merchant_id" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>ordernumber</td>
<td><textarea name="ordernumber"></textarea></td>
<td>
	<p>订单号</p>
</td>
</tr>
<tr>
<td>返回值</td>
<td colspan="2">
</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>测试代付<br />
/test/pay_for.do</p>
<table>
<form method="post" action="/test/pay_for.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>merchant_id</td>
<td><input name="merchant_id" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>params</td>
<td><textarea name="params"></textarea></td>
<td>
	<p>加密后的字符串</p>
</td>
</tr>
<tr>
<td>返回值</td>
<td colspan="2">
</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>测试代付订单查询<br />
/test/pay_for_query.do</p>
<table>
<form method="post" action="/test/pay_for_query.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>merchant_id</td>
<td><input name="merchant_id" type="text"></td>
<td>接入商英文编号</td>
</tr>
<tr>
<td>ordernumber</td>
<td><textarea name="ordernumber"></textarea></td>
<td>
	<p>订单号</p>
</td>
</tr>
<tr>
<td>返回值</td>
<td colspan="2">
</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

<p>测试支付支持的银行<br />
/test/pay_bank.do</p>
<table>
<form method="post" action="/test/pay_bank.do">
<tr>
<td>参数名</td>
<td>值</td>
<td>描述</td>
</tr>
<tr>
<td>is_mobile</td>
<td><input name="is_mobile" type="text"></td>
<td>是否手机版</td>
</tr>
<tr>
<td>pay_id</td>
<td><input name="pay_id" type="text"></td>
<td>
	<p>支付渠道ID</p>
</td>
</tr>
<tr>
<td>返回值</td>
<td colspan="2">
</td>
</tr>
<tr><td><input type="submit" value="ok"></td></tr>
</form>
</table>

</body>
</html>