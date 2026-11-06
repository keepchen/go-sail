package request

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/constants"
)

// GetUserInfo request parameters for retrieving user information
// swagger: model
type GetUserInfo struct {
	UserID int64 `json:"userId" form:"userId" validate:"required"` // User ID
}

func (v GetUserInfo) Validator() (constants.ICodeType, error) {
	if v.UserID < 1 {
		return constants.ErrRequestParamsInvalid, fmt.Errorf("field [userId], value:{%d} is invalid", v.UserID)
	}

	return constants.ErrNone, nil
}
