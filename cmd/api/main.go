package main

import (
	"log/slog"
	"os"
	"todoApp3/config"
	"todoApp3/internal/auth"
	"todoApp3/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()
	db := database.NewSQLite(cfg.SQLitePath)

	slogHandler := slog.NewJSONHandler(os.Stdout, nil)
	slogLogger := slog.New(slogHandler)

	validate := validator.New()

	repository := auth.NewRepository(db, slogLogger)
	service := auth.NewService(repository)
	handler := auth.NewHandler(service, validate)

	e := echo.New()
	e.Use(middleware.Recover())

	handler.Routes(e)

}
