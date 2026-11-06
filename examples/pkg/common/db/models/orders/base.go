package orders

import "gorm.io/gorm"

//AutoMigrate automatically synchronizes table schemas
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Order{},
	)
	if err != nil {
		panic(err)
	}
}
