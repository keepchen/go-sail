package request

import (
	"fmt"

	"github.com/keepchen/go-sail/v2/pkg/constants"
)

// GetUserInfo 获取用户信息请求参数
// swagger: model
type GetUserInfo struct {
	UserID int64 `json:"userId" form:"userId" validate:"required"` // 用户id
}

func (v GetUserInfo) Validator() (constants.ICodeType, error) {
	if v.UserID < 1 {
		return constants.ErrRequestParamsInvalid, fmt.Errorf("field [userId], value:{%d} is invalid", v.UserID)
	}

	return constants.ErrNone, nil
}
