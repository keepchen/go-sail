package api

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
)

// Responder response handler
type Responder interface {
	// Builder assembles response data
	//
	// This method automatically sets the HTTP status, message, current system timestamp in milliseconds, and request ID based on the given code (requires middleware.LogTrace in the route configuration)
	Builder(code constants.ICodeType, resp dto.IResponse, message ...string) Responder
	// Status specifies the HTTP status code
	//
	// This method overrides the HTTP status code derived by Assemble, Builder, SimpleAssemble, and Wrap
	Status(httpCode int) Responder
	// SendWithCode responds with the specified HTTP status code
	SendWithCode(httpCode int)
	// Send responds to the request
	Send()
	// Wrap assembles response data
	//
	// Unlike Builder, this method does not require the data argument to implement dto.IResponse
	//
	// This method automatically sets the HTTP status, message, current system timestamp in milliseconds, and request ID based on the given code (requires middleware.LogTrace in the route configuration)
	Wrap(code constants.ICodeType, data any, message ...string) Responder
	// Bundle wraps response data
	//
	// This method provides more flexible response data assembly
	//
	// # Tip
	//
	// Codes must be registered with sail.Code().Register() beforehand to work correctly
	//
	// For example:
	//
	// sail.Code().Register("en", 12345, "Code description...") //Register the code
	//
	// Response(c).Bundle(12345, nil).Send() //Use the code
	Bundle(code int, data any, message ...string) Responder
	// Data returns data
	//
	// This method directly returns response data, where:
	//
	// 1.The HTTP status code is 200
	//
	// 2.If Option.ErrNoneCode is set through SetupOption, the business code is Option.ErrNoneCode
	// Otherwise, the business code is constants.ErrNone
	Data(data any)
	// Success returns success
	//
	// This method directly returns a success status with default data (empty data), where:
	//
	// 1.The HTTP status code is 200
	//
	// 2.If Option.ErrNoneCode is set through SetupOption, the business code is Option.ErrNoneCode
	// Otherwise, the business code is constants.ErrNone
	Success()
	// Failure returns failure
	//
	// This method responds with default data (empty data), where:
	//
	// 1.The HTTP status code is 200
	//
	// 2.The business code is constants.ErrRequestParamsInvalid
	Failure(message ...string)
	// Failure200 returns failure
	//
	// This method responds with default data (empty data), where:
	//
	// 1.The HTTP status code is 200
	//
	// 2.The business code is code
	Failure200(code constants.ICodeType, message ...string)
	// Failure400 returns failure
	//
	// This method responds with default data (empty data), where:
	//
	// 1.The HTTP status code is 400
	//
	// 2.The business code is code
	Failure400(code constants.ICodeType, message ...string)
	// Failure500 returns failure
	//
	// This method responds with default data (empty data), where:
	//
	// 1.The HTTP status code is 500
	//
	// 2.The business code is code
	Failure500(code constants.ICodeType, message ...string)
	// WithHeaders sets response headers
	WithHeaders(headers map[string]string) Responder
	// WithCookies sets response cookies
	WithCookies(cookies []CookieStd) Responder
}

type responseEngine struct {
	engine    *gin.Context
	httpCode  int
	data      any
	requestId string
	spanId    string
	entryAt   int64
	language  string
}

var _ Responder = &responseEngine{}

// Define sync.Pool
var responseEnginePool = sync.Pool{
	New: func() any {
		return &responseEngine{}
	},
}

// Return responseEngine to the pool
func release(re *responseEngine) {
	re.engine = nil
	re.data = nil
	responseEnginePool.Put(re)
}

func New(c *gin.Context) Responder {
	// Get responseEngine from the pool
	re := responseEnginePool.Get().(*responseEngine)

	// Reset fields to avoid data contamination
	re.engine = c
	re.language = languageCode.String()
	re.requestId = ""
	re.spanId = ""
	re.entryAt = 0

	if val, ok := c.Get("requestId"); ok {
		re.requestId = val.(string)
	}
	if val, ok := c.Get("spanId"); ok {
		re.spanId = val.(string)
	}
	if val, ok := c.Get("entryAt"); ok {
		re.entryAt = val.(int64)
	}
	//Get the language code from the context
	if detectAcceptLanguage {
		if acceptLanguage, ok := c.Get("language"); ok {
			if lang, assertOk := acceptLanguage.(string); assertOk {
				re.language = lang
			}
		}
	}

	return re
}

// Response initializes a response instance
//
// # Usage
//
// import . "github.com/keepchen/go-sail/v3/http/api"
//
// Response(c).Builder(...).Send()
//
// Shorthand for New
func Response(c *gin.Context) Responder {
	return New(c)
}

// Builder assembles response data
//
// Shorthand for Assemble
func (a *responseEngine) Builder(code constants.ICodeType, resp dto.IResponse, message ...string) Responder {
	return a.mergeBody(code, resp, message...)
}

// Wrap assembles response data (lightweight version)
//
// Unlike Builder, this method does not require the data argument to implement dto.IResponse
//
// This method automatically sets the HTTP status, message, current system timestamp in milliseconds, and request ID based on the given code (requires middleware.LogTrace in the route configuration)
func (a *responseEngine) Wrap(code constants.ICodeType, data any, message ...string) Responder {
	return a.mergeBody(code, data, message...)
}

// Bundle wraps response data
//
// This method provides more flexible response data assembly
//
// # Tip
//
// Codes must be registered with sail.Code().Register() beforehand to work correctly
//
// For example:
//
// sail.Code().Register("en", 12345, "Code description...") //Register the code
//
// Response(c).Bundle(12345, nil).Send() //Use the code
func (a *responseEngine) Bundle(code int, resp any, message ...string) Responder {
	return a.mergeBody(constants.CodeType(code), resp, message...)
}

// WithHeaders sets response headers
func (a *responseEngine) WithHeaders(headers map[string]string) Responder {
	for k, v := range headers {
		a.engine.Header(k, v)
	}

	return a
}

// WithCookies sets response cookies
func (a *responseEngine) WithCookies(cookies []CookieStd) Responder {
	for _, cookie := range cookies {
		a.engine.SetCookie(cookie.Name, cookie.Value, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
	}

	return a
}

// Success returns success
//
// This method directly returns a success status with default data (empty data), where:
//
// 1.The HTTP status code is 200
//
// 2.If Option.ErrNoneCode is set through SetupOption, the business code is Option.ErrNoneCode
// Otherwise, the business code is constants.ErrNone
func (a *responseEngine) Success() {
	a.Data(nil)
}

// Failure returns failure
//
// This method responds with default data (empty data), where:
//
// 1.The HTTP status code is 200
//
// 2.The business code is constants.ErrRequestParamsInvalid
func (a *responseEngine) Failure(message ...string) {
	a.Failure200(constants.ErrRequestParamsInvalid, message...)
}

// Failure200 returns failure
//
// This method responds with default data (empty data), where:
//
// 1.The HTTP status code is 200
//
// 2.The business code is code
func (a *responseEngine) Failure200(code constants.ICodeType, message ...string) {
	a.Wrap(code, nil, message...).SendWithCode(http.StatusOK)
}

// Failure400 returns failure
//
// This method responds with default data (empty data), where:
//
// 1.The HTTP status code is 400
//
// 2.The business code is code
func (a *responseEngine) Failure400(code constants.ICodeType, message ...string) {
	a.Wrap(code, nil, message...).SendWithCode(http.StatusBadRequest)
}

// Failure500 returns failure
//
// This method responds with default data (empty data), where:
//
// 1.The HTTP status code is 500
//
// 2.The business code is code
func (a *responseEngine) Failure500(code constants.ICodeType, message ...string) {
	a.Wrap(code, nil, message...).SendWithCode(http.StatusInternalServerError)
}

// Data returns data
//
// This method directly returns response data, where:
//
// 1.The HTTP status code is 200
//
// 2.If Option.ErrNoneCode is set through SetupOption, the business code is Option.ErrNoneCode
// Otherwise, the business code is constants.ErrNone
func (a *responseEngine) Data(data any) {
	var code constants.ICodeType
	if anotherErrNoneCode != nil {
		code = anotherErrNoneCode
	} else {
		code = constants.ErrNone
	}

	a.Wrap(code, data).Send()
}

// Consolidate response body processing
func (a *responseEngine) mergeBody(code constants.ICodeType, resp any, message ...string) Responder {
	var (
		body     dto.Base
		httpCode int
	)

	body.RequestID = a.requestId
	body.Code = code.Int()
	if code == constants.ErrNone && anotherErrNoneCode != constants.ErrNone {
		//Use the overridden success code when the default success code has been changed and the current code is None
		body.Code = anotherErrNoneCode.Int()
	}
	body.Message = constants.CodeType(body.Code).String(a.language)

	switch code {
	case anotherErrNoneCode:
		body.Success = constants.Success
		httpCode = http.StatusOK
	case anotherErrRequestParamsInvalidCode:
		body.Success = constants.Failure
		httpCode = http.StatusBadRequest
	case anotherErrAuthorizationTokenInvalidCode:
		body.Success = constants.Failure
		httpCode = http.StatusUnauthorized
	case anotherErrInternalServerErrorCode:
		body.Success = constants.Failure
		httpCode = http.StatusInternalServerError
	default:
		body.Success = constants.Failure
		httpCode = http.StatusBadRequest
	}

	//Derive the HTTP status code from the business code unless Status was explicitly called
	if a.httpCode == 0 {
		a.httpCode = httpCode
		//Forcing HTTP status code 200 is enabled
		if forceHttpCode200 {
			a.httpCode = http.StatusOK
		}
	}

	//Override the default error message when message is provided
	if len(message) > 0 {
		var msg = strings.Builder{}
		msg.Grow(len(message) * 16)
		for index, v := range message {
			_, _ = msg.WriteString(v)
			if index < len(message)-1 {
				_, _ = msg.WriteString(";")
			}
		}
		body.Message = msg.String()
	}

	//Set the user response body
	//body.Data = resp
	if resp == nil {
		body.Data = emptyDataField
	} else if !isTypedNil(resp) {
		if iResp, ok := resp.(dto.IResponse); ok && iResp != nil {
			body.Data = iResp.GetData()
		} else {
			body.Data = resp
		}
	} else {
		body.Data = nil
	}

	if loc != nil {
		//loc may be used for other fields later; use this call for now
		body.Timestamp = time.Now().In(loc).UnixMilli()
	} else {
		body.Timestamp = time.Now().UnixMilli()
	}

	a.data = body

	return a
}

// Status specifies the HTTP status code
//
// This method overrides the HTTP status code derived by Assemble, Builder, SimpleAssemble, and Wrap
func (a *responseEngine) Status(httpCode int) Responder {
	a.httpCode = httpCode

	return a
}

// SendWithCode responds with the specified HTTP status code
//
// After calling this method:
//
// 1.The automatically derived status code is overridden
//
// 2.The Option.ForceHttpCode200 setting is ignored
func (a *responseEngine) SendWithCode(httpCode int) {
	defer release(a)

	if funcBeforeWrite != nil {
		var data dto.Base
		if val, ok := a.data.(dto.Base); ok {
			data = val
		} else {
			data.Data = a.data
		}

		funcBeforeWrite(a.engine.Request, a.entryAt, a.requestId, a.spanId, httpCode, data)
	}

	//Record quantile statistics
	if registry != nil {
		status := httpCode / 100
		registry.ObserveDuration(
			"http.server.duration",
			time.Since(time.Unix(0, a.entryAt)),
			distribution.Labels{
				"method": a.engine.Request.Method,
				"route":  a.engine.FullPath(),
				"status": strconv.Itoa(status) + "xx",
			},
		)
	}

	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

// Send responds to the request
func (a *responseEngine) Send() {
	a.SendWithCode(a.httpCode)
}

// Check whether an interface is a typed nil (its type is non-nil but its underlying value is nil)
func isTypedNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Chan, reflect.Interface, reflect.UnsafePointer,
		reflect.Map, reflect.Func, reflect.Pointer, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
