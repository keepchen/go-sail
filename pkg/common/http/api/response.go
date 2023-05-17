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
	Code constants.ICodeType `json:"code" example:"0" format:"int" validate:"required"`
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

var anotherErrNoneCode constants.ICodeType = constants.ErrNone

// OverrideErrNoneCode 改写默认的成功code码
//
// Usage example:
//
//	const AnotherErrNone = constants.CodeType(200)
//
//	api.OverrideErrNoneCode(AnotherErrNone, "some letters")
func OverrideErrNoneCode(errNoneCode constants.ICodeType, msg string) {
	constants.RegisterCode(constants.CodeType(errNoneCode.Int()), msg)
	anotherErrNoneCode = errNoneCode
}

// Assemble 组装返回数据
//
// 该方法会根据传递的code码自动设置http状态、描述信息、当前系统毫秒时间戳以及请求id(需要在路由配置中调用middleware.Before中间件)
func (a API) Assemble(code constants.ICodeType, resp IResponse, message ...string) API {
	var b Base
	if requestId, ok := a.engine.Get("requestId"); ok {
		b.RequestID = requestId.(string)
	}
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

	a.data = b

	return a
}

func (a API) SendWithCode(httpCode int) {
	a.engine.AbortWithStatusJSON(httpCode, a.data)
}

func (a API) Send() {
	a.SendWithCode(a.httpCode)
}
