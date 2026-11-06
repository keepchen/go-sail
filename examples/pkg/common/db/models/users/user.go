package users

import (
	"github.com/keepchen/go-sail/v3/examples/pkg/common/db/models"
	"github.com/keepchen/go-sail/v3/examples/pkg/common/db/models/orders"
	modelsEnum "github.com/keepchen/go-sail/v3/examples/pkg/common/enum/models"
	"github.com/keepchen/go-sail/v3/utils"
	"gorm.io/gorm"
)

type User struct {
	models.BaseModel
	UserID   int64                     `gorm:"column:user_id;type:bigint;not null;index:,unique;comment:用户ID"`
	Nickname string                    `gorm:"column:nickname;type:varchar(30);comment:用户昵称"`
	Status   modelsEnum.UserStatusCode `gorm:"column:status;type:tinyint;default:0;comment:用户状态"`
	Wallet   Wallet                    `gorm:"foreignKey:user_id;references:user_id"`
	Orders   []orders.Order            `gorm:"foreignKey:user_id;references:user_id"`
}

func (*User) TableName() string {
	return "tb_users"
}

func (*User) Class() string {
	return "User"
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.CreatedAt = utils.NewTimeWithTimeZone().Now()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = utils.NewTimeWithTimeZone().Now()

	return nil
}
