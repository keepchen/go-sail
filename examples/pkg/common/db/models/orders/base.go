package orders

import "gorm.io/gorm"

//AutoMigrate 自动同步表结构
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Order{},
	)
	if err != nil {
		panic(err)
	}
}
