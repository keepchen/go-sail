package constants

import (
	"fmt"
)

// Error codes
const (
	ErrNone                      = CodeType(0)      //No error; placeholder
	ErrRequestParamsInvalid      = CodeType(100000) //Invalid request parameters
	ErrAuthorizationTokenInvalid = CodeType(100001) //Token has expired
	ErrInternalServerError       = CodeType(999999) //Internal server error
)

// Error code message table
//
// # READONLY for concurrency safety
var initErrorCodeMsgMap = MMBox{
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

// String returns the error message string
func (ct CodeType) String(language ...string) string {
	var lang = LanguageEnglish // Use English by default
	if len(language) > 0 {
		lang = LanguageCode(language[0])
	}

	if val, ok := GetCodeMsg(lang, ct); ok {
		return val
	}
	//Fall back to English if the requested language is not English and cannot be found
	if lang != LanguageEnglish {
		if val, ok := GetCodeMsg(LanguageEnglish, ct); ok {
			return val
		}
	}

	return fmt.Sprintf("[Warn] ErrorCode {%d} under Language {%s} not defined!", ct, lang)
}

// Int returns the error code
func (ct CodeType) Int() int {
	return int(ct)
}
