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
	"todoApp3/internal/background"
	"todoApp3/internal/database"
	"todoApp3/internal/domain"
	"todoApp3/internal/httpx"
	appMw "todoApp3/internal/middleware"
	"todoApp3/internal/todo"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lmittmann/tint"
	"github.com/robfig/cron/v3"
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

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go auth.StartEmailWorker(ctxApp, smsChannel, &wg, slogLogger)
	}

	repo := background.NewTodoPruneRepoDB(db)
	pruner := background.NewTodoPruner(repo, slogLogger)

	c := cron.New()

	_, err := c.AddFunc("@every 10s", func() {
		ctx, cancel := context.WithTimeout(ctxApp, 5*time.Second)
		defer cancel()

		deleted, err := pruner.Prune(ctx, 30, 1000)
		if err != nil && ctx.Err() == nil {
			slogLogger.Error("todo prune failed", "err", err.Error(), "deleted_count", deleted)
			return
		}
		slogLogger.Debug("todo prune tick", "deleted_count", deleted)
	})
	if err != nil {
		slogLogger.Error("cron add failed", "err", err.Error())
	} else {
		c.Start()
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

	stopCtx := c.Stop()
	<-stopCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}
	close(smsChannel)
	wg.Wait()
	fmt.Println("Graceful shutdown completed.")
}
