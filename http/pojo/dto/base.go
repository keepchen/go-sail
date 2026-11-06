package dto

// Base common response parameters
// swagger: model
type Base struct {
	// Request ID
	// in: body
	// required: true
	RequestID string `json:"requestId" example:"5686efa5-c747-4f63-8657-e6052f8181a9" format:"string" validate:"required"`
	// Error codes
	// in: body
	// required: true
	Code int `json:"code" format:"int" example:"0" validate:"required"`
	// Whether successful
	// in: body
	// required: true
	Success bool `json:"success" example:"true" format:"bool" validate:"required"`
	// Message
	// in: body
	// required: true
	Message string `json:"message" example:"SUCCESS" format:"string" validate:"required"`
	// Server time (timestamp in milliseconds)
	// in: body
	// required: true
	Timestamp int64 `json:"ts" example:"1670899688591" format:"int64" validate:"required"`
	// Business data
	// in: body
	// required: true
	Data any `json:"data" format:"object|array|string|number|boolean" validate:"required"`
}
