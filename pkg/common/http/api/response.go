package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/constants"
)

type IResponse interface {
	GetData() interface{}
}

// Base 公共返回参数
// swagger: model
type Base struct {
	// 请求id
	// in: body
	// required: false
	RequestID string `json:"requestId" example:"1234567890123456789" format:"string"`
	// 错误码
	// in: body
	// required: true
	Code int `json:"code" example:"0" format:"int" validate:"required"`
	// 是否成功
	// in: body
	// required: true
	Success bool `json:"success" example:"true" format:"bool" validate:"required"`
	// 提示信息
	// in: body
	// required: true
	Message string `json:"message" example:"SUCCESS" format:"string" validate:"required"`
	// 服务器时间(毫秒时间戳)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// 业务数据
	// in: body
	// required: true
	Data interface{} `json:"data" validate:"required"`
}

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

func (a API) Assemble(code constants.ICodeType, resp IResponse, message ...string) API {
	var b Base
	b.Code = code.Int()
	b.Message = code.String()
	b.Timestamp = time.Now().UnixMilli()
	switch code {
	case constants.ErrNone:
		b.Success = constants.Success
		a.httpCode = http.StatusOK
	case constants.ErrRequestParamsInvalid:
		b.Success = constants.Failure
		a.httpCode = http.StatusBadRequest
	case constants.ErrAuthorizationTokenInvalid:
		b.Success = constants.Failure
		a.httpCode = http.StatusUnauthorized
	case constants.ErrInternalSeverError:
		b.Success = constants.Failure
		a.httpCode = http.StatusInternalServerError
	default:
		b.Success = constants.Failure
		a.httpCode = http.StatusBadRequest
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
	}

	requestID := a.engine.GetHeader("X-Request-Id")
	if len(requestID) != 0 {
		b.RequestID = requestID
	}

	a.data = b

	return a
}

func (a API) SendWithCode(httpCode int) {
	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

func (a API) Send() {
	a.SendWithCode(a.httpCode)
}
