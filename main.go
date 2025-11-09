package main

import (
	db "todo-api/db"
	"todo-api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Conn().Close()
	e := echo.New()
	e.GET("/todos", handlers.GetTodos)
	e.POST("/todos", handlers.CreateTodo)
	e.PUT("/todos/:id", handlers.UpdateTodos)

	log.Fatal(e.Start(":8080"))
}
