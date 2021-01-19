package database

import (
	"github.com/JayCuevas/gogurt/sync"
	"github.com/op/go-logging"
	"gorm.io/gorm"
)

var log = logging.MustGetLogger("gogurt")

type file struct {
	gorm.Model
	sync.File
}

type Manager struct {
	conn *gorm.DB
}

func (m *Manager) ApplyMigrations() error {
	return m.conn.AutoMigrate(&file{})
}

func (m *Manager) All() []sync.File {
	var files []sync.File
	m.conn.Select([]string{"path", "hash"}).Find(&files)
	return files
}

func (m *Manager) Upsert(f *sync.File) error {
	targetFile := &file{}
	m.conn.Find(targetFile, "path = ?", f.Path)

	// If the file does not exist
	if targetFile.Path == "" && targetFile.Hash == "" {
		result := m.conn.Create(&file{File: *f})

		return result.Error
	}

	if targetFile.Hash != f.Hash {
		targetFile.Hash = f.Hash
		m.conn.Save(targetFile)
	}

	return nil
}

func (m *Manager) Delete(paths []string) error {
	return nil
}

func (m *Manager) Close() {
	db, err := m.conn.DB()
	if err != nil {
		log.Errorf("Error retriving database connection to close: %s", err)
		return
	}

	db.Close()
}
