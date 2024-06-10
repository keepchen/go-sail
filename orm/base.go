package orm

import (
	"time"
)

// BaseModel 基础模型字段
//
// docs: https://gorm.io/docs/models.html
type BaseModel struct {
	ID        uint64     `gorm:"column:id;type:bigint UNSIGNED;primaryKey;AUTO_INCREMENT;comment:主键ID"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;comment:创建时间"`    //创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;comment:更新时间"`    //更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;comment:(软)删除时间"` //(软)删除时间
}

// NoneID 空ID
const NoneID = uint64(0)

var nowTime = time.Now()

// SetHookTime 设置勾子函数的时间对象
func SetHookTime(now time.Time) {
	nowTime = now
}
