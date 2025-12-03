package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("[database] failed to open sqlite db (%s): %v", path, err)
	}
	log.Printf("[database] sqlite connected (%s)", path)
	return db
}
