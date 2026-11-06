package orm

import (
	"time"
)

// BaseModel base model fields
//
// docs: https://gorm.io/docs/models.html
type BaseModel struct {
	ID        uint64     `gorm:"column:id;type:bigint UNSIGNED;primaryKey;AUTO_INCREMENT;comment:主键ID"`
	CreatedAt time.Time  `gorm:"column:created_at;type:datetime;comment:创建时间;autoCreateTime;<-:create"` //Creation time
	UpdatedAt time.Time  `gorm:"column:updated_at;type:datetime;comment:更新时间;autoUpdateTime"`           //Update time
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;comment:(软)删除时间"`                       //(Soft) deletion time
}

// NoneID empty ID
const NoneID = uint64(0)

var nowTimeFunc = func() time.Time {
	return time.Now()
}

// SetHookTimeFunc sets the time hook function
func SetHookTimeFunc(ntf func() time.Time) {
	nowTimeFunc = ntf
}
