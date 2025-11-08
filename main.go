package main

import (
	db "todo-api-db/db"
	gethandler "todo-api-db/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Db.Close()
	e := echo.New()
	e.GET("/todos", gethandler.GetTodos)

	log.Fatal(e.Start(":8080"))
}
