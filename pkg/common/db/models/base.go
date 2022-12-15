package models

import (
	"time"
)

//BaseModel 基础模型字段
type BaseModel struct {
	ID        uint64     `gorm:"column:id;primary_key; AUTO_INCREMENT"`
	CreatedAt time.Time  `gorm:"column:created_at;comment:创建时间"`    //创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at;comment:更新时间"`    //更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at;comment:(软)删除时间"` //(软)删除时间
}

//NoneID 空ID
const NoneID = uint64(0)
