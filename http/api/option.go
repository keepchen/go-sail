package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
)

var (
	emptyDataField       any = nil                       //Empty data field
	forceHttpCode200         = false                     //Force HTTP status code 200
	timezone                 = constants.DefaultTimeZone //Time zone
	detectAcceptLanguage     = false                     //Whether to detect the client language
	languageCode             = constants.LanguageEnglish //Language code
)

// Built-in error codes that can be overridden
var (
	anotherErrNoneCode                      constants.ICodeType = constants.ErrNone                      //Overridden success code
	anotherErrRequestParamsInvalidCode      constants.ICodeType = constants.ErrRequestParamsInvalid      //Overridden invalid parameter code
	anotherErrAuthorizationTokenInvalidCode constants.ICodeType = constants.ErrAuthorizationTokenInvalid //Overridden expired token code
	anotherErrInternalServerErrorCode       constants.ICodeType = constants.ErrInternalServerError       //Overridden internal server error code
)

var (
	loc             *time.Location                                                                                                 //Time zone
	funcBeforeWrite func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base) //Handler invoked before writing the response
	registry        *distribution.Registry                                                                                         //Distributed quantile statistics
)

// Option configuration options
type Option struct {
	//Success code, default: 0. When configured, this value is used as the success code.
	ErrNoneCode constants.ICodeType
	//Success message, default: SUCCESS. When configured, this value is used as the success message.
	ErrNoneCodeMsg string
	//Invalid parameter code, default: 100000. When configured, this value is used as the invalid parameter code.
	ErrRequestParamsInvalidCode constants.ICodeType
	//Invalid parameter message. When configured, this value is used as the invalid parameter message.
	ErrRequestParamsInvalidCodeMsg string
	//Expired token code, default: 100000. When configured, this value is used as the expired token code.
	ErrAuthorizationTokenInvalidCode constants.ICodeType
	//Expired token message. When configured, this value is used as the expired token message.
	ErrAuthorizationTokenInvalidCodeMsg string
	//Internal server error code, default: 999999. When configured, this value is used as the internal server error code.
	ErrInternalServerErrorCode constants.ICodeType
	//Internal server error message. When configured, this value is used as the internal server error message.
	ErrInternalServerErrorCodeMsg string
	//Empty data serialization format, default: null. When configured, this value is used for empty data fields.
	EmptyDataStruct int
	//Force HTTP status code 200. When enabled, the HTTP status code is not derived from the business code.
	//
	//Note that Responder.Status and Responder.SendWithCode take precedence over this option.
	ForceHttpCode200 bool
	//Time zone
	Timezone string
	//Whether to detect the client language for error messages
	DetectAcceptLanguage bool
	//Language code
	//
	//Use this language code when DetectAcceptLanguage is disabled
	LanguageCode constants.LanguageCode
	//Handler invoked before writing the response
	//
	//request pointer to the HTTP request struct
	//
	//entryAtUnixNano nanosecond timestamp when the request arrives
	//
	//requestId request ID
	//
	//spanId unique ID for the internal call chain
	//
	//httpCode HTTP status code
	//
	//writeData data about to be written
	//
	//# Note: this function is called synchronously. Avoid blocking or time-consuming operations in performance-sensitive scenarios.
	FuncBeforeWrite func(request *http.Request, entryAtUnixNano int64, requestId, spanId string, httpCode int, writeData dto.Base)
}

const (
	DefaultEmptyDataStructNull   int = 0 //Empty (null)
	DefaultEmptyDataStructObject int = 1 //Braces (curly brackets)
	DefaultEmptyDataStructArray  int = 2 //Brackets (square brackets)
	DefaultEmptyDataStructString int = 3 //Empty string
)

// SetupOption sets options
//
// Currently supports:
//
// 1.Built-in error codes and error messages
//
// 2.Serialization format for empty data
//
// # Tip
//
// This method is not reentrant; callers should use it globally and call it at most once
func SetupOption(opt Option) {
	//Override built-in error codes
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

// DefaultSetupOption default settings
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

// SetRegistry sets distributed quantile statistics
func SetRegistry(reg *distribution.Registry) {
	registry = reg
}
