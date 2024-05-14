package orm

import (
	"time"

	"gorm.io/gorm"
)

func (u *BaseModel) BeforeSave(_ *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *BaseModel) BeforeCreate(_ *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *BaseModel) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()

	return nil
}

func (u *BaseModel) BeforeDelete(_ *gorm.DB) (err error) {
	now := time.Now()
	u.DeletedAt = &now

	return nil
}
