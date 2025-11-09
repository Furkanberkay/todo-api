package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() error {
	var err error

	db, err = sql.Open("sqlite", "todos1.db")
	if err != nil {
		return fmt.Errorf("db open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db ping error: %w", err)
	}

	return nil
}

func Conn() *sql.DB {
	return db
}
