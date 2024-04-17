package dto

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
	Code int `json:"code" format:"int" example:"200" validate:"required"`
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
	Data interface{} `json:"data" format:"object|array|string" validate:"required"`
}
