package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteStorage(storageName string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(storageName), &gorm.Config{})
}
