package orm

import (
	"errors"

	"gorm.io/gorm"
)

// AutoMigrate automatically synchronizes table schemas
func AutoMigrate(db *gorm.DB, tables ...any) error {
	err := db.AutoMigrate(tables...)
	return err
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

// IgnoreErrRecordNotFound ignores record-not-found errors
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
