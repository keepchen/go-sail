package mock

import (
	"github.com/keepchen/go-sail/v2/pkg/common/db/models"
	modelsUsers "github.com/keepchen/go-sail/v2/pkg/common/db/models/users"
	usersSvc "github.com/keepchen/go-sail/v2/pkg/common/db/service/users"
	modelsEnum "github.com/keepchen/go-sail/v2/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v2/pkg/lib/db"
	"github.com/keepchen/go-sail/v2/pkg/lib/logger"
)

// CreateUserAndWalletData mock用户及钱包数据
//
// 生成三条数据
func CreateUserAndWalletData() {
	users := []modelsUsers.User{
		{
			UserID:   1,
			Nickname: "go-sail",
			Status:   modelsEnum.UserStatusCodeNormal,
			Wallet: modelsUsers.Wallet{
				UserID: 1,
				Amount: 100.00,
				Status: modelsEnum.WalletStatusCodeNormal,
			},
		},
		{
			UserID:   2,
			Nickname: "keepchen",
			Status:   modelsEnum.UserStatusCodeNormal,
			Wallet: modelsUsers.Wallet{
				UserID: 2,
				Amount: 200.00,
				Status: modelsEnum.WalletStatusCodeNormal,
			},
		},
		{
			UserID:   3,
			Nickname: "corgi",
			Status:   modelsEnum.UserStatusCodeNormal,
			Wallet: modelsUsers.Wallet{
				UserID: 3,
				Amount: 300.00,
				Status: modelsEnum.WalletStatusCodeNormal,
			},
		},
	}
	svc := usersSvc.NewUserSvcImpl(db.GetInstance().R, db.GetInstance().W, logger.GetLogger())
	for _, newUser := range users {
		user, _ := svc.GetUserAndWallet(newUser.UserID)
		if user.ID == models.NoneID {
			_ = svc.CreateUserAndWalletWithUserID(newUser.UserID, newUser.Nickname, newUser.Wallet.Amount)
		}
	}
}
