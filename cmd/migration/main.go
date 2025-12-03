package main

import (
	"todoApp3/config"
	"todoApp3/internal/database"
)

func main() {
	cfg := config.Load()

	db := database.NewSQLite(cfg.SQLitePath)
	db.AutoMigrate(db)
}
