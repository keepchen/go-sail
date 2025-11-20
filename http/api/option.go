package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

var (
	emptyDataField       any = nil                       //空data字段
	forceHttpCode200         = false                     //强制使用200作为http的状态码
	timezone                 = constants.DefaultTimeZone //时区
	detectAcceptLanguage     = false                     //是否检测客户端语言
	languageCode             = constants.LanguageEnglish //语言代码
)

// 可被覆盖的内置错误码
var (
	anotherErrNoneCode                      constants.ICodeType = constants.ErrNone                      //被改写后的成功code码
	anotherErrRequestParamsInvalidCode      constants.ICodeType = constants.ErrRequestParamsInvalid      //被修改后的参数错误code码
	anotherErrAuthorizationTokenInvalidCode constants.ICodeType = constants.ErrAuthorizationTokenInvalid //被修改后的令牌失效code码
	anotherErrInternalServerErrorCode       constants.ICodeType = constants.ErrInternalServerError       //被修改后的服务内部错误code码
)

var (
	loc             *time.Location                                                                                                 //时区
	funcBeforeWrite func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base) //写入响应前的处理函数
)

// Option 配置项
type Option struct {
	//成功code码，默认成功code码为0，配置此项后，成功code码将使用这个值。
	ErrNoneCode constants.ICodeType
	//成功提示语，默认成功提示语为SUCCESS，配置此项后，成功提示语将使用这个值。
	ErrNoneCodeMsg string
	//参数错误code码，默认参数错误code码为100000，配置此项后，参数错误code码将使用这个值。
	ErrRequestParamsInvalidCode constants.ICodeType
	//参数错误提示语，配置此项后，参数错误提示语将使用这个值。
	ErrRequestParamsInvalidCodeMsg string
	//令牌失效code码，默认令牌失效code码为100000，配置此项后，令牌失效code码将使用这个值。
	ErrAuthorizationTokenInvalidCode constants.ICodeType
	//令牌失效提示语，配置此项后，令牌失效提示语将使用这个值。
	ErrAuthorizationTokenInvalidCodeMsg string
	//服务器内部错误code码，默认服务器内部错误code码为999999，配置此项后，服务器内部错误code码将使用这个值。
	ErrInternalServerErrorCode constants.ICodeType
	//服务器内部错误提示语，配置此项后，服务器内部错误提示语将使用这个值。
	ErrInternalServerErrorCodeMsg string
	//空data序列化结构，默认返回的data字段为空时为null值，配置此项后，空data序列化格式将使用这个值。
	EmptyDataStruct int
	//强制使用200作为http的状态码，配置此项后，http状态码将不从业务code码中解析。
	//
	//注意，调用 Responder.Status 方法和 Responder.SendWithCode 方法时的优先级高于此项配置。
	ForceHttpCode200 bool
	//时区
	Timezone string
	//是否检测客户端语言，用于错误码消息返回
	DetectAcceptLanguage bool
	//语言代码
	//
	//当没有启用 DetectAcceptLanguage 时，使用该语言代码
	LanguageCode constants.LanguageCode
	//写入响应前的处理函数
	//
	//request http请求结构体指针
	//
	//entryAtUnixNano 请求到达一刻的纳秒时间戳
	//
	//requestId 请求id
	//
	//spanId 内部调用链的唯一id
	//
	//httpCode http状态码
	//
	//writeData 即将写入的数据
	//
	//# 注意，该函数是同步调用，对于性能敏感的场景，不建议在函数体内做阻塞或耗时的操作。
	FuncBeforeWrite func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base)
}

const (
	DefaultEmptyDataStructNull   int = 0 //空(null)
	DefaultEmptyDataStructObject int = 1 //大(花)括号
	DefaultEmptyDataStructArray  int = 2 //中(方)括号
	DefaultEmptyDataStructString int = 3 //空字符串
)

// SetupOption 设置选项
//
// 目前支持设定:
//
// 1.内置错误码及错误码信息
//
// 2.空数据序列化结构
//
// # 提示
//
// 该方法不支持重入，因此，调用者应当在全局使用且最多调用一次
func SetupOption(opt Option) {
	//覆盖内置错误码
	if opt.ErrNoneCode != nil {
		constants.RegisterCodeSingle(constants.LanguageEnglish, opt.ErrNoneCode, opt.ErrNoneCodeMsg)
		anotherErrNoneCode = opt.ErrNoneCode
	}
	if opt.ErrRequestParamsInvalidCode != nil {
		constants.RegisterCodeSingle(constants.LanguageEnglish, opt.ErrRequestParamsInvalidCode, opt.ErrRequestParamsInvalidCodeMsg)
		anotherErrRequestParamsInvalidCode = opt.ErrRequestParamsInvalidCode
	}
	if opt.ErrAuthorizationTokenInvalidCode != nil {
		constants.RegisterCodeSingle(constants.LanguageEnglish, opt.ErrAuthorizationTokenInvalidCode, opt.ErrAuthorizationTokenInvalidCodeMsg)
		anotherErrAuthorizationTokenInvalidCode = opt.ErrAuthorizationTokenInvalidCode
	}
	if opt.ErrInternalServerErrorCode != nil {
		constants.RegisterCodeSingle(constants.LanguageEnglish, opt.ErrInternalServerErrorCode, opt.ErrInternalServerErrorCodeMsg)
		anotherErrInternalServerErrorCode = opt.ErrInternalServerErrorCode
	}
	switch opt.EmptyDataStruct {
	case DefaultEmptyDataStructNull:
		emptyDataField = nil
	case DefaultEmptyDataStructObject:
		emptyDataField = struct{}{}
	case DefaultEmptyDataStructArray:
		emptyDataField = []bool{}
	case DefaultEmptyDataStructString:
		emptyDataField = ""
	default:
		emptyDataField = nil
	}

	if opt.ForceHttpCode200 {
		forceHttpCode200 = opt.ForceHttpCode200
	}

	if len(opt.Timezone) > 0 {
		timezone = opt.Timezone
	}

	if opt.DetectAcceptLanguage {
		detectAcceptLanguage = opt.DetectAcceptLanguage
	}

	lc, err := time.LoadLocation(timezone)
	if err != nil {
		panic(fmt.Errorf("[GO-SAIL] can not load location: %s", timezone))
	}
	loc = lc

	if opt.FuncBeforeWrite != nil {
		funcBeforeWrite = opt.FuncBeforeWrite
	}

	if len(opt.LanguageCode) != 0 {
		languageCode = opt.LanguageCode
	}
}

// DefaultSetupOption 默认设置
func DefaultSetupOption() *Option {
	return &Option{
		Timezone:                            constants.DefaultTimeZone,
		ErrNoneCode:                         constants.ErrNone,
		ErrNoneCodeMsg:                      constants.ErrNone.String(),
		ErrRequestParamsInvalidCode:         constants.ErrRequestParamsInvalid,
		ErrRequestParamsInvalidCodeMsg:      constants.ErrRequestParamsInvalid.String(),
		ErrAuthorizationTokenInvalidCode:    constants.ErrAuthorizationTokenInvalid,
		ErrAuthorizationTokenInvalidCodeMsg: constants.ErrAuthorizationTokenInvalid.String(),
		ErrInternalServerErrorCode:          constants.ErrInternalServerError,
		ErrInternalServerErrorCodeMsg:       constants.ErrInternalServerError.String(),
		ForceHttpCode200:                    false,
		LanguageCode:                        constants.LanguageEnglish,
	}
}
