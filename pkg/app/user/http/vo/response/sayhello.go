package response

import (
	"github.com/keepchen/go-sail/pkg/common/http/api"
)

// SayHello 欢迎返回数据结构
// swagger: model
type SayHello struct {
	api.Base
	// 数据体
	// in: body
	// required: true
	Data string `json:"data" validate:"required"`
}

func (v SayHello) GetData() interface{} {
	return v.Data
}

var _ api.IResponse = &SayHello{}
