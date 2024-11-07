package service

import (
	"errors"

	"go.uber.org/zap"

	"github.com/keepchen/go-sail/v3/sail"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/vo/request"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/vo/response"
	usersSvc "github.com/keepchen/go-sail/v3/examples/pkg/common/db/service/users"
	"github.com/keepchen/go-sail/v3/lib/db"
	"gorm.io/gorm"
)

func GetUserInfoSvc(c *gin.Context) {
	var (
		form      request.GetUserInfo
		resp      response.GetUserInfo
		logTracer = sail.LogTrace(c)
	)
	if err := c.ShouldBind(&form); err != nil {
		sail.Response(c).Builder(constants.ErrRequestParamsInvalid, nil).Send()
		return
	}

	if errorCode, err := form.Validator(); err != nil {
		sail.Response(c).Builder(errorCode, nil, err.Error()).Send()
		return
	}

	userAndWallet, sqlErr := usersSvc.NewUserSvcImpl(db.GetInstance().R, db.GetInstance().W, logTracer.GetLogger()).GetUserAndWallet(form.UserID)
	if sqlErr != nil && errors.Is(sqlErr, gorm.ErrRecordNotFound) {
		sail.Response(c).Builder(constants.ErrRequestParamsInvalid, nil, "user not found").Send()
		return
	} else {
		logTracer.Warn("query failed, cause: ", zap.String("err", sqlErr.Error()))
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

	sail.Response(c).Builder(constants.ErrNone, resp).Send()
}
