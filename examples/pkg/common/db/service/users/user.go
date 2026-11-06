package users

import (
	"github.com/keepchen/go-sail/v3/examples/pkg/common/db/models/users"
	modelsEnum "github.com/keepchen/go-sail/v3/examples/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserSvc interface {
	//CreateUserAndWallet creates a user and the user's wallet
	CreateUserAndWallet(nickname string, amount float64) error
	//CreateUserAndWalletWithUserID creates a user and wallet (with a specified userID)
	CreateUserAndWalletWithUserID(userID int64, nickname string, amount float64) error
	//GetUserAndWallet retrieves user information and balance
	GetUserAndWallet(userID int64) (users.User, error)
}

type UserSvcImpl struct {
	dbr    *gorm.DB //Read instance
	dbw    *gorm.DB //Write instance
	logger *zap.Logger
}

var _ UserSvc = &UserSvcImpl{}

// NewUserSvcImpl instantiates the service
var NewUserSvcImpl = func(dbr *gorm.DB, dbw *gorm.DB, logger *zap.Logger) UserSvc {
	return &UserSvcImpl{
		dbr:    dbr,
		dbw:    dbw,
		logger: logger,
	}
}

func (u UserSvcImpl) CreateUserAndWallet(nickname string, amount float64) error {
	userID := utils.NewTimeWithTimeZone().Now().UnixNano()
	userAndWallet := users.User{
		UserID:   userID,
		Nickname: nickname,
		Status:   modelsEnum.UserStatusCodeNormal,
		Wallet: users.Wallet{
			UserID: userID,
			Amount: amount,
			Status: modelsEnum.WalletStatusCodeNormal,
		},
	}
	err := u.dbw.Create(&userAndWallet).Error
	if err != nil {
		u.logger.Error("数据库操作:CreateUserAndWallet:错误",
			zap.Any("value", logger.MarshalInterfaceValue(userAndWallet)), zap.Errors("errors", []error{err}))
	}

	return err
}

func (u UserSvcImpl) CreateUserAndWalletWithUserID(userID int64, nickname string, amount float64) error {
	userAndWallet := users.User{
		UserID:   userID,
		Nickname: nickname,
		Status:   modelsEnum.UserStatusCodeNormal,
		Wallet: users.Wallet{
			UserID: userID,
			Amount: amount,
			Status: modelsEnum.WalletStatusCodeNormal,
		},
	}
	err := u.dbw.Create(&userAndWallet).Error
	if err != nil {
		u.logger.Error("数据库操作:CreateUserAndWallet:错误",
			zap.Any("value", logger.MarshalInterfaceValue(userAndWallet)), zap.Errors("errors", []error{err}))
	}

	return err
}

func (u UserSvcImpl) GetUserAndWallet(userID int64) (users.User, error) {
	var userAndWallet users.User
	err := u.dbr.Model(&users.User{}).Where(&users.User{UserID: userID}).Preload(userAndWallet.Wallet.Class()).First(&userAndWallet).Error
	if err != nil {
		u.logger.Error("数据库操作:GetUserAndWallet:错误",
			zap.Any("value", userID), zap.Errors("errors", []error{err}))
	}

	return userAndWallet, err
}
