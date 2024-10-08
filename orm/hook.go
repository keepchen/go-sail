package orm

import (
	"gorm.io/gorm"
)

func (u *BaseModel) BeforeSave(_ *gorm.DB) (err error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = nowTimeFunc()
		u.UpdatedAt = u.CreatedAt
	} else {
		u.UpdatedAt = nowTimeFunc()
	}

	return nil
}

func (u *BaseModel) BeforeCreate(_ *gorm.DB) (err error) {
	u.CreatedAt = nowTimeFunc()
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *BaseModel) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = nowTimeFunc()

	return nil
}

func (u *BaseModel) BeforeDelete(_ *gorm.DB) (err error) {
	now := nowTimeFunc()
	u.DeletedAt = &now

	return nil
}
