package main

import (
	"database/sql"
	"log"
	"os"
	"todo-api/handlers"
	"todo-api/repository"
	"todo-api/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "todos1.db")
	if err != nil {
		log.Fatalf("db open failed: %v", err)
	}
	defer db.Close()

	logger := log.New(os.Stdout, "[todo-api] ", log.LstdFlags|log.Lshortfile)

	repo := repository.NewSqliteTodoRepository(db, logger)
	todoService := service.NewTodoService(repo)
	todoHandler := handlers.NewTodoHandler(todoService)

	e := echo.New()
	e.Use(middleware.Recover())

	todoHandler.ClientRouters(e)

	log.Fatal(e.Start(":8080"))
}

//katman service repository, yapılandır handler,DI
//config dosyası
