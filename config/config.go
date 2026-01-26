package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr             string
	SQLitePath           string
	SecretKey            string
	JwtExpirationMinutes int
	AppName              string
	RefReshExpDays       int

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err == nil {
			return i
		}
	}
	return fallback
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddr:             getEnv("HTTP_ADDR", ":8080"),
		SQLitePath:           getEnv("SQLITE_PATH", "./todo.db"),
		SecretKey:            os.Getenv("SECRET_KEY"),
		JwtExpirationMinutes: getEnvAsInt("JWT_EXPIRATION_MINUTES", 60),
		AppName:              os.Getenv("APP_NAME"),
		RefReshExpDays:       getEnvAsInt("REFRESH_TOKEN_TTL_DAYS", 30),

		DBHost:     getEnv("DB_HOST", "db"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "todo_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
	}
	if cfg.SecretKey == "" {
		log.Fatal("CRITICAL ERROR: SECRET_KEY environment variable not set! Application cannot be started.")
	}

	log.Printf("[config] App=%s Addr=%s DBHost=%s", cfg.AppName, cfg.HTTPAddr, cfg.DBHost)

	return cfg
}
