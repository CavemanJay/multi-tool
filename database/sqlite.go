// +build release

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewManager() (*Manager, error) {
	db, err := gorm.Open(sqlite.Open("./data/multi-tool.db?mode=rwc"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// db.AutoMigrate(&User{})

	return &Manager{conn: db}, nil
}
