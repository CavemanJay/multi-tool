// +build !release

package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewManager() (*Manager, error) {
	dsn := "root:toor@tcp(127.0.0.1:3306)/gogurt?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// db.AutoMigrate(&User{})

	return &Manager{conn: db}, nil
}
