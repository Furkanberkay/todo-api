package database

import (
	"fmt"
	"log"
	"log/slog"
	"todoApp3/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresSql struct {
	logger *slog.Logger
	config *config.Config
}

func NewPostgresSql(logger *slog.Logger, cfg *config.Config) *PostgresSql {
	return &PostgresSql{
		logger: logger,
		config: cfg,
	}
}

func (p *PostgresSql) Connect() *gorm.DB {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Istanbul",
		p.config.DBHost,
		p.config.DBUser,
		p.config.DBPassword,
		p.config.DBName,
		p.config.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[database] Failed to connect to Postgres!\nError: %v\nDSN: host=%s port=%s (Password hidden)", err, p.config.DBHost, p.config.DBPort)
	}

	p.logger.Info("[database] postgresql connection successful!", "host", p.config.DBHost)

	return db
}
