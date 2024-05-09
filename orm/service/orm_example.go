package service

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/orm/model"
)

type User struct {
	model.Base
	UserID   int64  `gorm:"column:user_id;type:bigint;not null;index:,unique;comment:用户ID"`
	Nickname string `gorm:"column:nickname;type:varchar(30);comment:用户昵称"`
	Status   int    `gorm:"column:status;type:tinyint;default:0;comment:用户状态"`
}

func ExampleORMUsage() {
	svc := NewORMSvcImpl(db.GetInstance().R, db.GetInstance().W, logger.GetLogger())

	// ---- read one record
	var user User
	err := svc.R().Where(&User{UserID: 1000}).First(&user)
	fmt.Println(err)

	// ---- create record
	var user0 = User{
		UserID:   1000,
		Nickname: "go-sail",
		Status:   1,
	}
	err = svc.W().Create(&user0)
	fmt.Println(err)

	// ---- force update all fields except some one
	var user1 = User{
		UserID:   1000,
		Nickname: "go-sail",
		Status:   1,
	}
	err = svc.W().Select("*").Omit("deleted_at").Updates(&user1)
	fmt.Println(err)

	// ---- paginate
	var (
		users    []User
		page     = 1
		pageSize = 50
	)
	total, err := svc.R().Paginate(users, page, pageSize)
	fmt.Println(total, err)
}
