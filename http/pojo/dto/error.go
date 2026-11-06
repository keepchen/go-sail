package dto

// Error400 error response parameters
// swagger: model
type Error400 struct {
	// Request ID
	// in: body
	// required: false
	RequestID string `json:"requestId" example:"5686efa5-c747-4f63-8657-e6052f8181a9" format:"string" validate:"required"`
	// Error codes
	// in: body
	// required: true
	Code int `json:"code" example:"100000" format:"int" validate:"required"`
	// Whether successful
	// in: body
	// required: true
	Success bool `json:"success" example:"false" format:"bool" validate:"required"`
	// Message
	// in: body
	// required: true
	Message string `json:"message" example:"Bad request parameters" format:"string" validate:"required"`
	// Server time (timestamp in milliseconds)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// Business data
	// in: body
	// required: true
	Data any `json:"data" format:"object|array|string|boolean|number" validate:"required"`
}

// Error500 error response parameters
// swagger: model
type Error500 struct {
	// Request ID
	// in: body
	// required: false
	RequestID string `json:"requestId" example:"5686efa5-c747-4f63-8657-e6052f8181a9" format:"string" validate:"required"`
	// Error codes
	// in: body
	// required: true
	Code int `json:"code" example:"999999" format:"int" validate:"required"`
	// Whether successful
	// in: body
	// required: true
	Success bool `json:"success" example:"false" format:"bool" validate:"required"`
	// Message
	// in: body
	// required: true
	Message string `json:"message" example:"Internal server error" format:"string" validate:"required"`
	// Server time (timestamp in milliseconds)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// Business data
	// in: body
	// required: true
	Data any `json:"data" format:"object|array|string|boolean|number" validate:"required"`
}
