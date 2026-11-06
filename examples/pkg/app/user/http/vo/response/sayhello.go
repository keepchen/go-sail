package response

import (
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// SayHello greeting response data structure
// swagger: model
type SayHello struct {
	dto.Base
	// Data body
	// in: body
	// required: true
	Data string `json:"data" example:"" format:"string" validate:"required"`
}

func (v SayHello) GetData() interface{} {
	return v.Data
}

var _ dto.IResponse = &SayHello{}
