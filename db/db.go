package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var Db *sql.DB

func InitDB() error {
	var err error

	Db, err = sql.Open("sqlite", "todos1.db")
	if err != nil {
		return fmt.Errorf("db open error: %w", err)
	}

	if err := Db.Ping(); err != nil {
		return fmt.Errorf("db ping error: %w", err)
	}

	return nil
}
