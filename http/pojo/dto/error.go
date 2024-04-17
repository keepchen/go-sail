package dto

// Error400 错误返回参数
// swagger: model
type Error400 struct {
	// 请求id
	// in: body
	// required: false
	RequestID string `json:"requestId" example:"5686efa5-c747-4f63-8657-e6052f8181a9" format:"string"`
	// 错误码
	// in: body
	// required: true
	Code int `json:"code" example:"100000" format:"int" validate:"required"`
	// 是否成功
	// in: body
	// required: true
	Success bool `json:"success" example:"false" format:"bool" validate:"required"`
	// 提示信息
	// in: body
	// required: true
	Message string `json:"message" example:"Bad request parameters" format:"string" validate:"required"`
	// 服务器时间(毫秒时间戳)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// 业务数据
	// in: body
	// required: true
	Data interface{} `json:"data" format:"object|array|string" validate:"required"`
}

// Error500 错误返回参数
// swagger: model
type Error500 struct {
	// 请求id
	// in: body
	// required: false
	RequestID string `json:"requestId" example:"5686efa5-c747-4f63-8657-e6052f8181a9" format:"string"`
	// 错误码
	// in: body
	// required: true
	Code int `json:"code" example:"999999" format:"int" validate:"required"`
	// 是否成功
	// in: body
	// required: true
	Success bool `json:"success" example:"false" format:"bool" validate:"required"`
	// 提示信息
	// in: body
	// required: true
	Message string `json:"message" example:"Internal server error" format:"string" validate:"required"`
	// 服务器时间(毫秒时间戳)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// 业务数据
	// in: body
	// required: true
	Data interface{} `json:"data" format:"object|array|string" validate:"required"`
}
