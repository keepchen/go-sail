package api

import (
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// Responder 响应器
type Responder interface {
	// Builder 组装返回数据
	//
	// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
	Builder(code constants.ICodeType, resp dto.IResponse, message ...string) Responder
	// Assemble 组装返回数据
	//
	// Deprecated: Assemble is deprecated,it will be removed in the future.
	//
	// Please use Builder instead.
	//
	// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
	Assemble(code constants.ICodeType, resp dto.IResponse, message ...string) Responder
	// Status 指定http状态码
	//
	// 该方法会覆盖 Assemble, Builder, SimpleAssemble, Wrap 解析的http状态码值
	Status(httpCode int) Responder
	// SendWithCode 以指定http状态码响应请求
	SendWithCode(httpCode int)
	// Send 响应请求
	Send()
	// Wrap 组装返回数据
	//
	// 该方法与 Builder 的区别在于data参数不需要实现 dto.IResponse 接口
	//
	// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
	Wrap(code constants.ICodeType, data interface{}, message ...string) Responder
	// SimpleAssemble 组装返回数据(轻量版)
	//
	// Deprecated: SimpleAssemble is deprecated,it will be removed in the future.
	//
	// Please use Wrap instead.
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
	engine    *gin.Context
	httpCode  int
	data      interface{}
	requestId string
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
	return a.mergeBody(code, resp, message...)
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
	a.Wrap(code, nil, message...).SendWithCode(http.StatusOK)
}

// Failure400 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为400
//
// 2.业务码为code
func (a *responseEngine) Failure400(code constants.ICodeType, message ...string) {
	a.Wrap(code, nil, message...).SendWithCode(http.StatusBadRequest)
}

// Failure500 返回失败
//
// 此方法响应默认数据(data为空)，其中：
//
// 1.http状态码为500
//
// 2.业务码为code
func (a *responseEngine) Failure500(code constants.ICodeType, message ...string) {
	a.Wrap(code, nil, message...).SendWithCode(http.StatusInternalServerError)
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

	a.Wrap(code, data).Send()
}

// Wrap 组装返回数据(轻量版)
//
// 该方法与 Builder 的区别在于data参数不需要实现 dto.IResponse 接口
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
func (a *responseEngine) Wrap(code constants.ICodeType, data interface{}, message ...string) Responder {
	return a.mergeBody(code, data, message...)
}

// SimpleAssemble 组装返回数据(轻量版)
//
// Deprecated: SimpleAssemble is deprecated,it will be removed in the future.
//
// Please use Wrap instead.
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用 middleware.LogTrace 中间件)
func (a *responseEngine) SimpleAssemble(code constants.ICodeType, data interface{}, message ...string) Responder {
	return a.mergeBody(code, data, message...)
}

// Assemble 组装返回数据
//
// Deprecated: Assemble is deprecated,it will be removed in the future.
//
// Please use Builder instead.
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
		language  = []string{languageCode.String()}
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
		var msg = strings.Builder{}
		for index, v := range message {
			_, _ = msg.WriteString(v)
			if index < len(message)-1 {
				_, _ = msg.WriteString(";")
			}
		}
		body.Message = msg.String()
	}

	//设置用户响应体
	body.Data = resp
	if !isTypedNil(resp) {
		if iResp, ok := resp.(dto.IResponse); ok && !isTypedNil(iResp) {
			body.Data = iResp.GetData()
		} else {
			body.Data = resp
		}
	} else {
		body.Data = nil
	}

	//当空data配置项不为nil时，需要判断用户响应体数据类型并覆盖空data配置项
	if emptyDataField != nil {
		vf := reflect.ValueOf(body.Data)
		switch true {
		//为空
		case body.Data == nil:
			body.Data = emptyDataField
		//为指针类型且为空
		case (vf.Kind() == reflect.Pointer ||
			vf.Kind() == reflect.Slice ||
			vf.Kind() == reflect.Map ||
			vf.Kind() == reflect.Interface ||
			vf.Kind() == reflect.Func ||
			vf.Kind() == reflect.Chan ||
			vf.Kind() == reflect.UnsafePointer) && vf.IsNil():
			body.Data = emptyDataField
		//数组或切片类型且长度为0
		case (vf.Kind() == reflect.Slice || vf.Kind() == reflect.Array) && vf.Len() == 0:
			body.Data = emptyDataField
		}
	}

	if loc != nil {
		//loc可能用于后续其他字段，这里暂时这样调用
		body.Timestamp = time.Now().In(loc).UnixMilli()
	} else {
		body.Timestamp = time.Now().UnixMilli()
	}

	a.requestId = requestId
	a.data = body

	return a
}

// Status 指定http状态码
//
// 该方法会覆盖 Assemble, Builder, SimpleAssemble, Wrap 解析的http状态码值
func (a *responseEngine) Status(httpCode int) Responder {
	a.httpCode = httpCode

	return a
}

// SendWithCode 以指定http状态码响应请求
//
// 调用此方法后：
//
// 1.会覆盖自动推导的状态码
//
// 2.会忽略 Option.ForceHttpCode200 设置
func (a *responseEngine) SendWithCode(httpCode int) {
	if funcBeforeWrite != nil {
		var (
			spanId  string
			entryAt int64
		)
		if val, ok := a.engine.Get("spanId"); ok {
			spanId = val.(string)
		}
		if val, ok := a.engine.Get("entryAt"); ok {
			entryAt = val.(int64)
		}

		var data dto.Base
		if val, ok := a.data.(dto.Base); ok {
			data = val
		} else {
			data.Data = a.data
		}

		funcBeforeWrite(a.engine.Request, entryAt, a.requestId, spanId, httpCode, data)
	}
	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

// Send 响应请求
func (a *responseEngine) Send() {
	a.SendWithCode(a.httpCode)
}

// 判断一个interface是否为typed nil（即类型不为nil，但内部值为nil）
func isTypedNil(v interface{}) bool {
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
