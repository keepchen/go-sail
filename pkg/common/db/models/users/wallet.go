package users

import (
	"github.com/keepchen/go-sail/v2/pkg/common/db/models"
	modelsEnum "github.com/keepchen/go-sail/v2/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v2/pkg/utils"
	"gorm.io/gorm"
)

type Wallet struct {
	models.BaseModel
	UserID int64                       `gorm:"column:user_id;type:bigint;not null;index:,unique;comment:用户ID"`
	Amount float64                     `gorm:"column:amount;type:decimal(14,4);comment:余额"`
	Status modelsEnum.WalletStatusCode `gorm:"column:status;type:tinyint;default:0;comment:用户状态"`
}

func (Wallet) TableName() string {
	return "tb_wallets"
}

func (*Wallet) Class() string {
	return "Wallet"
}

func (w *Wallet) BeforeCreate(_ *gorm.DB) (err error) {
	w.BaseModel.CreatedAt = utils.NewTimeWithTimeZone().Now()
	w.BaseModel.UpdatedAt = w.BaseModel.CreatedAt

	return nil
}

func (w *Wallet) BeforeUpdate(_ *gorm.DB) (err error) {
	w.BaseModel.UpdatedAt = utils.NewTimeWithTimeZone().Now()

	return nil
}
