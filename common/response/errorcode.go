package response

import (
	"strconv"
)

var ErrorMap map[string]string

type Error struct {
	Code    int
	Message string
}

type Errors struct {
	Error map[string]Error
}

func init() {
	ErrorMap = map[string]string{
		//系统级别错误
		"10000": "无效的API",
		"10001": "未定义错误",
		"10002": "已知错误",
		"10003": "系统错误",
		"10004": "数据错误",
		"10011": "配置信息异常",


	}
}

func GetError(code int) string {
	key := strconv.Itoa(code)

	if msg, ok := ErrorMap[key]; ok {
		return msg
	}

	return "未知错误"
}
