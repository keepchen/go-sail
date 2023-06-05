package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v2/pkg/common/http/pojo/dto"
	"github.com/keepchen/go-sail/v2/pkg/constants"
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

// Assemble 组装返回数据
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用middleware.Before中间件)
func (a API) Assemble(code constants.ICodeType, resp dto.IResponse, message ...string) API {
	var (
		b         dto.Base
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
	b.RequestID = requestId
	b.Code = code
	if code == constants.ErrNone && anotherErrNoneCode != constants.ErrNone {
		//改写了默认成功code码，且当前code码为None时，需要使用改写后的值
		b.Code = anotherErrNoneCode
	}
	b.Message = b.Code.String()
	b.Timestamp = time.Now().UnixMilli()
	switch code {
	case constants.ErrNone, anotherErrNoneCode:
		b.Success = constants.Success
		httpCode = http.StatusOK
	case constants.ErrRequestParamsInvalid:
		b.Success = constants.Failure
		httpCode = http.StatusBadRequest
	case constants.ErrAuthorizationTokenInvalid:
		b.Success = constants.Failure
		httpCode = http.StatusUnauthorized
	case constants.ErrInternalSeverError:
		b.Success = constants.Failure
		httpCode = http.StatusInternalServerError
	default:
		b.Success = constants.Failure
		httpCode = http.StatusBadRequest
	}

	//如果没有单独调用Status方法设置http状态码，则从code码中解析出http状态码
	if a.httpCode == 0 {
		a.httpCode = httpCode
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
		b.Message = msg.String()
	}

	if resp != nil {
		b.Data = resp.GetData()
	} else {
		b.Data = emptyDataField
	}

	a.data = b

	return a
}

// Status 指定http状态码
//
// 该方法会覆盖Assemble解析的http状态码值
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
