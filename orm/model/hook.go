package model

import (
	"time"

	"gorm.io/gorm"
)

func (u *Base) BeforeSave(_ *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *Base) BeforeCreate(_ *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *Base) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()

	return nil
}

func (u *Base) BeforeDelete(_ *gorm.DB) (err error) {
	now := time.Now()
	u.DeletedAt = &now

	return nil
}
