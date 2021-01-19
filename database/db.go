package database

import (
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("gogurt")

type User struct {
	gorm.Model
	Name  string
	Email string
	test  int
}

type Manager struct {
	conn *gorm.DB
}
