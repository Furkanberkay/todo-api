package main

import (
	"todo-api/db"
	"todo-api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer db.Conn().Close()
	e := echo.New()
	e.Use(middleware.Recover())
	e.GET("/todos", handlers.GetTodos)
	e.GET("/todos/:id", handlers.GetTodosById)
	e.POST("/todos", handlers.CreateTodo)
	e.PUT("/todos/:id", handlers.UpdateTodos)
	e.PATCH("/todos/:id", handlers.PatchTodo)
	e.DELETE("/todos/:id", handlers.DeleteTodo)
	log.Fatal(e.Start(":8080"))
}
