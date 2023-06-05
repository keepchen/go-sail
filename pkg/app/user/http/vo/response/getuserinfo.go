package response

import (
	modelsEnum "github.com/keepchen/go-sail/v2/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v2/pkg/common/http/pojo/dto"
)

// GetUserInfo 获取用户信息返回数据结构
// swagger: model
type GetUserInfo struct {
	dto.Base
	// 数据体
	// in: body
	// required: true
	Data struct {
		User   UserInfo   `json:"user"`
		Wallet WalletInfo `json:"wallet"`
	} `json:"data" format:"object"`
}

// UserInfo 用户基础信息数据结构
// swagger: model
type UserInfo struct {
	UserID int64 `json:"userId" validate:"required"` // 用户id
	// 用户昵称
	// in: body
	// required: true
	Nickname string `json:"userInfo" validate:"required"`
	// 账号状态
	//
	// UserStatusCodeNormal    = UserStatusCode(0) //正常
	// UserStatusCodeForbidden = UserStatusCode(1) //禁用
	//
	// in: body
	// required: true
	Status modelsEnum.UserStatusCode `json:"status" enums:"0,1" validate:"required"`
}

// WalletInfo 钱包信息数据结构
// swagger: model
type WalletInfo struct {
	// 账户余额
	// in: body
	// required: true
	Amount float64 `json:"amount" validate:"required"`
	// 钱包状态
	//
	// WalletStatusCodeNormal    = WalletStatusCode(0) //正常
	// WalletStatusCodeForbidden = WalletStatusCode(1) //禁用
	//
	// in: body
	// required: true
	Status modelsEnum.WalletStatusCode `json:"status" enums:"0,1" validate:"required"`
}

func (v GetUserInfo) GetData() interface{} {
	return v.Data
}

var _ dto.IResponse = &GetUserInfo{}
