package orm

import (
	"errors"

	"gorm.io/gorm"
)

// AutoMigrate 自动同步表结构
func AutoMigrate(db *gorm.DB, tables ...interface{}) error {
	err := db.AutoMigrate(tables...)
	return err
}

// Paginate 分页器
//
// page默认为0，pageSize默认为10
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

// IgnoreErrRecordNotFound 忽略记录未找到的错误
//
// @docs https://gorm.io/docs/v2_release_note.html#ErrRecordNotFound
//
// Example:
//
// err := IgnoreErrRecordNotFound(db.First())
func IgnoreErrRecordNotFound(db *gorm.DB) error {
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		if db.Statement != nil {
			db.Statement.RaiseErrorOnNotFound = false
		}
		db.Error = nil
	}

	return db.Error
}
