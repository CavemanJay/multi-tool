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

func (m *Manager) Close() {
	db, err := m.conn.DB()
	if err != nil {
		log.Errorf("Error retriving database connection to close: %s", err)
		return
	}

	db.Close()
}
