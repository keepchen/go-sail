package constants

import "fmt"

// 错误码
const (
	ErrNone                      = CodeType(0)      //没有错误，占位
	ErrRequestParamsInvalid      = CodeType(100000) //请求参数有误
	ErrAuthorizationTokenInvalid = CodeType(100001) //令牌已失效
	ErrInternalSeverError        = CodeType(999999) //服务器内部错误
)

// 错误码信息表
//
// READONLY for concurrency safety
var initErrorCodeMsgMap = map[CodeType]string{
	ErrNone:                      "SUCCESS",
	ErrRequestParamsInvalid:      "Bad request parameters",
	ErrAuthorizationTokenInvalid: "Token invalid",
	ErrInternalSeverError:        "Internal server error",
}

// String 获取错误信息字符
func (ct CodeType) String() string {
	if msg, ok := ctm.maps[ct]; ok {
		return msg
	}

	return fmt.Sprintf("[Warn] ErrorCode {%d} not defined!", ct)
}

// Int 获取错误码
func (ct CodeType) Int() int {
	return int(ct)
}
