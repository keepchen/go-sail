package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

type API struct {
	engine   *gin.Context
	httpCode int
	data     interface{}
}

func New(c *gin.Context) API {
	return API{
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
func Response(c *gin.Context) API {
	return New(c)
}

// Builder 组装返回数据
//
// Assemble 方法的语法糖
func (a API) Builder(code constants.ICodeType, resp dto.IResponse, message ...string) API {
	return a.Assemble(code, resp, message...)
}

// Assemble 组装返回数据
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用middleware.Before中间件)
func (a API) Assemble(code constants.ICodeType, resp dto.IResponse, message ...string) API {
	var (
		body      dto.Base
		requestId string
		httpCode  int
	)
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
	body.Code = code
	if code == constants.ErrNone && anotherErrNoneCode != constants.ErrNone {
		//改写了默认成功code码，且当前code码为None时，需要使用改写后的值
		body.Code = anotherErrNoneCode
	}
	body.Message = body.Code.String()
	body.Timestamp = time.Now().In(loc).UnixMilli()
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

	if resp != nil {
		body.Data = resp.GetData()
	} else {
		body.Data = emptyDataField
	}

	a.data = body

	return a
}

// Status 指定http状态码
//
// 该方法会覆盖 Assemble 解析的http状态码值
func (a API) Status(httpCode int) API {
	a.httpCode = httpCode

	return a
}

// SendWithCode 以指定http状态码响应请求
func (a API) SendWithCode(httpCode int) {
	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

// Send 响应请求
func (a API) Send() {
	a.SendWithCode(a.httpCode)
}
