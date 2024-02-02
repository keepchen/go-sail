package constants

import "fmt"

// 错误码
const (
	ErrNone                      = CodeType(0)      //没有错误，占位
	ErrRequestParamsInvalid      = CodeType(100000) //请求参数有误
	ErrAuthorizationTokenInvalid = CodeType(100001) //令牌已失效
	ErrInternalServerError       = CodeType(999999) //服务器内部错误
)

// 错误码信息表
//
// READONLY for concurrency safety
var initErrorCodeMsgMap = map[LanguageCode]map[ICodeType]string{
	LanguageEnglish: {
		ErrNone:                      "SUCCESS",
		ErrRequestParamsInvalid:      "Bad request parameters",
		ErrAuthorizationTokenInvalid: "Authorization token invalid",
		ErrInternalServerError:       "Internal server error",
	},
	LanguageEnglishUnitedStates: {
		ErrNone:                      "SUCCESS",
		ErrRequestParamsInvalid:      "Bad request parameters",
		ErrAuthorizationTokenInvalid: "Authorization token invalid",
		ErrInternalServerError:       "Internal server error",
	},
	LanguageChinesePRC: {
		ErrNone:                      "成功",
		ErrRequestParamsInvalid:      "无效的请求参数",
		ErrAuthorizationTokenInvalid: "无效的授权令牌",
		ErrInternalServerError:       "服务器内部错误",
	},
}

// String 获取错误信息字符
func (ct CodeType) String(language ...string) string {
	var lang = LanguageEnglish
	if len(language) > 0 {
		//TODO 这里需要考虑当语言不存在时，是否强制使用某种存在的语言代码
		lang = LanguageCode(language[0])
	}
	ctm.mux.RLock()
	defer ctm.mux.RUnlock()
	if i18nMsg, ok := ctm.maps[lang]; ok {
		if msg, iOk := i18nMsg[ct]; iOk {
			return msg
		}
		return fmt.Sprintf("[Warn] ErrorCode {%d} Language {%s} not defined!", ct, lang)
	}

	return fmt.Sprintf("[Warn] ErrorCode {%d} Language {%s} not defined!", ct, lang)
}

// Int 获取错误码
func (ct CodeType) Int() int {
	return int(ct)
}
