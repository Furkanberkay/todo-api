package main

import (
	"log/slog"
	"os"
	"time"
	"todoApp3/config"
	"todoApp3/internal/auth"
	"todoApp3/internal/database"
	appMw "todoApp3/internal/middleware"
	"todoApp3/internal/todo"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

func main() {
	cfg := config.Load()
	db := database.NewSQLite(cfg.SQLitePath)

	slogHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
		AddSource:  true,
	})

	slogLogger := slog.New(slogHandler)
	slog.SetDefault(slogLogger)
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
