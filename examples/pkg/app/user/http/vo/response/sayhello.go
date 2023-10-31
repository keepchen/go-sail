package response

import (
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// SayHello 欢迎返回数据结构
// swagger: model
type SayHello struct {
	dto.Base
	// 数据体
	// in: body
	// required: true
	Data string `json:"data" example:"" format:"string" validate:"required"`
}

func (v SayHello) GetData() interface{} {
	return v.Data
}

var _ dto.IResponse = &SayHello{}
