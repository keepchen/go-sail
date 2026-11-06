package response

import (
	modelsEnum "github.com/keepchen/go-sail/v3/examples/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
)

// GetUserInfo response data structure for retrieving user information
// swagger: model
type GetUserInfo struct {
	dto.Base
	// Data body
	// in: body
	// required: true
	Data struct {
		User   UserInfo   `json:"user"`
		Wallet WalletInfo `json:"wallet"`
	} `json:"data" format:"object"`
}

// UserInfo basic user information data structure
// swagger: model
type UserInfo struct {
	UserID int64 `json:"userId" validate:"required"` // User ID
	// User nickname
	// in: body
	// required: true
	Nickname string `json:"userInfo" validate:"required"`
	// Account status
	//
	// UserStatusCodeNormal    = UserStatusCode(0) //Normal
	// UserStatusCodeForbidden = UserStatusCode(1) //Disabled
	//
	// in: body
	// required: true
	Status modelsEnum.UserStatusCode `json:"status" enums:"0,1" validate:"required"`
}

// WalletInfo wallet information data structure
// swagger: model
type WalletInfo struct {
	// Account balance
	// in: body
	// required: true
	Amount float64 `json:"amount" validate:"required"`
	// Wallet status
	//
	// WalletStatusCodeNormal    = WalletStatusCode(0) //Normal
	// WalletStatusCodeForbidden = WalletStatusCode(1) //Disabled
	//
	// in: body
	// required: true
	Status modelsEnum.WalletStatusCode `json:"status" enums:"0,1" validate:"required"`
}

func (v GetUserInfo) GetData() interface{} {
	return v.Data
}

var _ dto.IResponse = &GetUserInfo{}
