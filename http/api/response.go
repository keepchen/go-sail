package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// Responder 响应器
type Responder interface {
	// Builder 组装返回数据
	//
	// Assemble 方法的语法糖
	Builder(code constants.ICodeType, resp dto.IResponse, message ...string) Responder
	// Assemble 组装返回数据
	//
	// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
	Assemble(code constants.ICodeType, resp dto.IResponse, message ...string) Responder
	// Status 指定http状态码
	//
	// 该方法会覆盖 Assemble 解析的http状态码值
	Status(httpCode int) Responder
	// SendWithCode 以指定http状态码响应请求
	SendWithCode(httpCode int)
	// Send 响应请求
	Send()
	// SimpleAssemble 组装返回数据(轻量版)
	//
	// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
	SimpleAssemble(code constants.ICodeType, data interface{}, message ...string) Responder
	// Data 返回数据
	//
	// 此方法直接返回响应数据，其中：
	//
	// 1.http状态码为200
	//
	// 2.当调用 SetupOption 设置了 Option.ErrNoneCode 那么业务code为 Option.ErrNoneCode
	// 否则业务code为 constants.ErrNone
	Data(data interface{})
	// Success 返回成功
	//
	// 此方法直接返回成功状态，响应默认数据(data为空)，其中：
	//
	// 1.http状态码为200
	//
	// 2.当调用 SetupOption 设置了 Option.ErrNoneCode 那么业务code为 Option.ErrNoneCode
	// 否则业务code为 constants.ErrNone
	Success()
	// Failure 返回失败
	//
	// 此方法响应默认数据(data为空)，其中：
	//
	// 1.http状态码为200
	//
	// 2.业务code码为 constants.ErrRequestParamsInvalid
	Failure(message ...string)
	// Failure200 返回失败
	//
	// 此方法响应默认数据(data为空)，其中：
	//
	// 1.http状态码为200
	//
	// 2.业务码为code
	Failure200(code constants.ICodeType, message ...string)
	// Failure400 返回失败
	//
	// 此方法响应默认数据(data为空)，其中：
	//
	// 1.http状态码为400
	//
	// 2.业务码为code
	Failure400(code constants.ICodeType, message ...string)
	// Failure500 返回失败
	//
	// 此方法响应默认数据(data为空)，其中：
	//
	// 1.http状态码为500
	//
	// 2.业务码为code
	Failure500(code constants.ICodeType, message ...string)
}

type responseEngine struct {
	engine   *gin.Context
	httpCode int
	data     interface{}
}

var _ Responder = &responseEngine{}

func New(c *gin.Context) Responder {
	return &responseEngine{
		engine: c,
	}
}

// Response 启动返回实例
//
// # Usage
//
// import . "github.com/keepchen/go-sail/v3/http/api"
//
// Response(c).Builder(...).Send()
//
// New 方法的语法糖
func Response(c *gin.Context) Responder {
	return New(c)
}

// Builder 组装返回数据
//
// Assemble 方法的语法糖
func (a *responseEngine) Builder(code constants.ICodeType, resp dto.IResponse, message ...string) Responder {
	return a.Assemble(code, resp, message...)
}

// Success 返回成功
//
// 此方法直接返回成功状态，响应默认数据(data为空)，其中：
//
// 1.http状态码为200
//
// 2.当调用 SetupOption 设置了 Option.ErrNoneCode 那么业务code为 Option.ErrNoneCode
// 否则业务code为 constants.ErrNone
func (a *responseEngine) Success() {
	a.Data(nil)
}

// Failure 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为200
//
// 2.业务code码为 constants.ErrRequestParamsInvalid
func (a *responseEngine) Failure(message ...string) {
	a.Failure200(constants.ErrRequestParamsInvalid, message...)
}

// Failure200 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为200
//
// 2.业务码为code
func (a *responseEngine) Failure200(code constants.ICodeType, message ...string) {
	a.SimpleAssemble(code, nil, message...).SendWithCode(http.StatusOK)
}

// Failure400 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为400
//
// 2.业务码为code
func (a *responseEngine) Failure400(code constants.ICodeType, message ...string) {
	a.SimpleAssemble(code, nil, message...).SendWithCode(http.StatusBadRequest)
}

// Failure500 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为500
//
// 2.业务码为code
func (a *responseEngine) Failure500(code constants.ICodeType, message ...string) {
	a.SimpleAssemble(code, nil, message...).SendWithCode(http.StatusInternalServerError)
}

// Data 返回数据
//
// 此方法直接返回响应数据，其中：
//
// 1.http状态码为200
//
// 2.当调用 SetupOption 设置了 Option.ErrNoneCode 那么业务code为 Option.ErrNoneCode
// 否则业务code为 constants.ErrNone
func (a *responseEngine) Data(data interface{}) {
	var code constants.ICodeType
	if anotherErrNoneCode != nil {
		code = anotherErrNoneCode
	} else {
		code = constants.ErrNone
	}

	a.SimpleAssemble(code, data).Send()
}

// SimpleAssemble 组装返回数据(轻量版)
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
func (a *responseEngine) SimpleAssemble(code constants.ICodeType, data interface{}, message ...string) Responder {
	return a.mergeBody(code, data, message...)
}

// Assemble 组装返回数据
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
func (a *responseEngine) Assemble(code constants.ICodeType, resp dto.IResponse, message ...string) Responder {
	return a.mergeBody(code, resp, message...)
}

// 合并处理响应体
func (a *responseEngine) mergeBody(code constants.ICodeType, resp interface{}, message ...string) Responder {
	var (
		body      dto.Base
		requestId string
		httpCode  int
		language  = []string{"en"}
	)
	//从上下文中获取语言代码
	if detectAcceptLanguage {
		if acceptLanguage, ok := a.engine.Get("language"); ok {
			if lang, assertOk := acceptLanguage.(string); assertOk {
				language[0] = lang
			}
		}
	}
	//从header中读取X-Request-Id
	if r1Id := a.engine.GetHeader("X-Request-Id"); len(r1Id) > 0 {
		requestId = r1Id
		//从header中读取requestId
	} else if r2Id := a.engine.GetHeader("requestId"); len(r2Id) > 0 {
		requestId = r2Id
	} else {
		//从上下文中读取requestId
		if r3Id, ok := a.engine.Get("requestId"); ok {
			requestId = r3Id.(string)
		}
	}
	body.RequestID = requestId
	body.Code = code.Int()
	if code == constants.ErrNone && anotherErrNoneCode != constants.ErrNone {
		//改写了默认成功code码，且当前code码为None时，需要使用改写后的值
		body.Code = anotherErrNoneCode.Int()
	}
	body.Message = constants.CodeType(body.Code).String(language...)
	if loc != nil {
		//loc可能用于后续其他字段，这里暂时这样调用
		body.Timestamp = time.Now().In(loc).UnixMilli()
	} else {
		body.Timestamp = time.Now().UnixMilli()
	}
	switch code {
	case constants.ErrNone, anotherErrNoneCode:
		body.Success = constants.Success
		httpCode = http.StatusOK
	case constants.ErrRequestParamsInvalid:
		body.Success = constants.Failure
		httpCode = http.StatusBadRequest
	case constants.ErrAuthorizationTokenInvalid:
		body.Success = constants.Failure
		httpCode = http.StatusUnauthorized
	case constants.ErrInternalServerError:
		body.Success = constants.Failure
		httpCode = http.StatusInternalServerError
	default:
		body.Success = constants.Failure
		httpCode = http.StatusBadRequest
	}

	//如果没有单独调用Status方法设置http状态码，则从code码中解析出http状态码
	if a.httpCode == 0 {
		a.httpCode = httpCode
		//开启了强制使用200作为http状态码
		if forceHttpCode200 {
			a.httpCode = http.StatusOK
		}
	}

	//如果message有值，则覆盖默认错误码所代表的错误信息
	if len(message) > 0 {
		var msg bytes.Buffer
		for index, v := range message {
			_, _ = msg.Write([]byte(v))
			if index < len(message)-1 {
				_, _ = msg.Write([]byte(";"))
			}
		}
		body.Message = msg.String()
	}

	if iResp, ok := resp.(dto.IResponse); resp != nil && ok {
		body.Data = iResp.GetData()
	} else if resp != nil {
		body.Data = resp
	} else {
		body.Data = emptyDataField
	}

	a.data = body

	return a
}

// Status 指定http状态码
//
// 该方法会覆盖 Assemble 解析的http状态码值
func (a *responseEngine) Status(httpCode int) Responder {
	a.httpCode = httpCode

	return a
}

// SendWithCode 以指定http状态码响应请求
func (a *responseEngine) SendWithCode(httpCode int) {
	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

// Send 响应请求
func (a *responseEngine) Send() {
	a.SendWithCode(a.httpCode)
}
