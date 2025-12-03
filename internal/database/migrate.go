package database

import (
	"log"
	"todoApp3/internal/domain"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&domain.User{}, domain.Todo{}); err != nil {
		log.Fatalf("[database] automigrate failed: %v", err)
	}

	log.Printf("[database] automigrate completed (Users + Todos)")
}
