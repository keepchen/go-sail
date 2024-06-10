package orm

import (
	"gorm.io/gorm"
)

func (u *BaseModel) BeforeSave(_ *gorm.DB) (err error) {
	u.CreatedAt = nowTime
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *BaseModel) BeforeCreate(_ *gorm.DB) (err error) {
	u.CreatedAt = nowTime
	u.UpdatedAt = u.CreatedAt

	return nil
}

func (u *BaseModel) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = nowTime

	return nil
}

func (u *BaseModel) BeforeDelete(_ *gorm.DB) (err error) {
	now := nowTime
	u.DeletedAt = &now

	return nil
}
