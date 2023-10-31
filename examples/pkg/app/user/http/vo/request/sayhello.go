package request

import "github.com/keepchen/go-sail/v3/constants"

// SayHello 欢迎请求参数
// swagger: model
type SayHello struct {
	Nickname string `json:"nickname" form:"nickname" query:"nickname"` // 昵称
}

func (v SayHello) Validator() (constants.ICodeType, error) {
	return constants.ErrNone, nil
}
