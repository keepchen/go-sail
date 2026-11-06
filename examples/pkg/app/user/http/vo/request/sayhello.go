package request

import "github.com/keepchen/go-sail/v3/constants"

// SayHello greeting request parameters
// swagger: model
type SayHello struct {
	Nickname string `json:"nickname" form:"nickname" query:"nickname"` // Nickname
}

func (v SayHello) Validator() (constants.ICodeType, error) {
	return constants.ErrNone, nil
}
