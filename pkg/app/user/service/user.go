package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/app/user/http/vo/request"
	"github.com/keepchen/go-sail/pkg/app/user/http/vo/response"
	usersSvc "github.com/keepchen/go-sail/pkg/common/db/service/users"
	"github.com/keepchen/go-sail/pkg/common/http/api"
	"github.com/keepchen/go-sail/pkg/constants"
	"github.com/keepchen/go-sail/pkg/lib/db"
	"github.com/keepchen/go-sail/pkg/lib/logger"
	"gorm.io/gorm"
)

func GetUserInfoSvc(c *gin.Context) {
	var (
		form request.GetUserInfo
		resp response.GetUserInfo
	)
	if err := c.ShouldBind(&form); err != nil {
		api.New(c).Assemble(constants.ErrRequestParamsInvalid, nil).Send()
		return
	}

	if errorCode, err := form.Validator(); err != nil {
		api.New(c).Assemble(errorCode, nil, err.Error()).Send()
		return
	}

	userAndWallet, sqlErr := usersSvc.NewUserSvcImpl(db.GetInstance().R, db.GetInstance().W, logger.GetLogger()).GetUserAndWallet(form.UserID)
	if sqlErr != nil && errors.Is(sqlErr, gorm.ErrRecordNotFound) {
		api.New(c).Assemble(constants.ErrRequestParamsInvalid, nil, "user not found").Send()
		return
	}

	resp.Data.User = response.UserInfo{
		UserID:   userAndWallet.UserID,
		Nickname: userAndWallet.Nickname,
		Status:   userAndWallet.Status,
	}
	resp.Data.Wallet = response.WalletInfo{
		Amount: userAndWallet.Wallet.Amount,
		Status: userAndWallet.Wallet.Status,
	}

	api.New(c).Assemble(constants.ErrNone, resp).Send()
}
