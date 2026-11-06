package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel base model fields
type BaseModel struct {
	ID        uint64     `gorm:"column:id;type:bigint;primary_key; AUTO_INCREMENT"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;comment:创建时间"`    //Creation time
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;comment:更新时间"`    //Update time
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;comment:(软)删除时间"` //(Soft) deletion time
}

// NoneID empty ID
const NoneID = uint64(0)

// AutoMigrate automatically synchronizes table schemas
func AutoMigrate(db *gorm.DB, tables ...interface{}) {
	err := db.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}

// Paginate paginator
//
// page defaults to 0 and pageSize defaults to 10
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 0
		} else {
			page--
		}

		if pageSize < 1 {
			pageSize = 10
		}

		return db.Offset(page * pageSize).Limit(pageSize)
	}
}
