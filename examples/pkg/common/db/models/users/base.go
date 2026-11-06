package users

import "gorm.io/gorm"

//AutoMigrate automatically synchronizes table schemas
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&Wallet{},
	)
	if err != nil {
		panic(err)
	}
}
