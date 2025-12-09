package main

import (
	"io"
	"log/slog"
	"os"
	"todoApp3/config"
	"todoApp3/internal/auth"
	"todoApp3/internal/database"
	appMw "todoApp3/internal/middleware"
	"todoApp3/internal/todo"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()
	db := database.NewSQLite(cfg.SQLitePath)

	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Log dosyası oluşturulamadı: " + err.Error())
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	slogHandler := slog.NewJSONHandler(multiWriter, nil)
	slogLogger := slog.New(slogHandler)

	validate := validator.New()

	authRepository := auth.NewRepository(db, slogLogger)
	AuthService := auth.NewService(authRepository, cfg, slogLogger)
	authHandler := auth.NewHandler(AuthService, validate)

	todoRepository := todo.NewRepository(db, slogLogger)
	todoService := todo.NewTodoService(todoRepository, slogLogger)
	todoHandler := todo.NewTodoHandler(todoService, slogLogger, validate)

	authMiddleware := appMw.NewAuthMiddleware(cfg)

	e := echo.New()
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")
	protected := api.Group("")
	protected.Use(authMiddleware.ValidateJwt)

	authHandler.Routes(api)
	todoHandler.Routes(protected)

	e.Start(cfg.HTTPAddr)

}
