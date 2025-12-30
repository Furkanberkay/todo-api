package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"todoApp3/config"
	"todoApp3/internal/auth"
	"todoApp3/internal/database"
	"todoApp3/internal/domain"
	"todoApp3/internal/httpx"
	appMw "todoApp3/internal/middleware"
	"todoApp3/internal/todo"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
)

func main() {

	ctxApp, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

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

	smsChannel := make(chan domain.SmsJob, 100)
	wg := sync.WaitGroup{}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go auth.StartSmsWorker(ctxApp, smsChannel, &wg, slogLogger)
	}

	authRepository := auth.NewRepository(db, slogLogger)
	AuthService := auth.NewService(authRepository, cfg, slogLogger, smsChannel)
	authHandler := auth.NewHandler(AuthService, validate, slogLogger)

	todoRepository := todo.NewRepository(db, slogLogger)
	todoService := todo.NewTodoService(todoRepository, slogLogger)
	todoHandler := todo.NewTodoHandler(todoService, slogLogger, validate)

	//authMiddleware := appMw.NewAuthMiddleware(cfg)
	jwtVerify := httpx.NewJwtVerify(cfg.SecretKey, slogLogger)
	authenticationMiddleware := appMw.NewAuthenticationMiddleware(jwtVerify)

	e := echo.New()
	e.Use(middleware.Recover())

	api := e.Group("/api/v1")
	protected := api.Group("")
	protected.Use(authenticationMiddleware.Authenticate)

	authHandler.Routes(api)
	todoHandler.Routes(protected)

	go func() {
		if err := e.Start(cfg.HTTPAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(ctxApp, 10*time.Minute)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}
	close(smsChannel)
	wg.Wait()
	fmt.Println("Graceful shutdown completed.")
}
