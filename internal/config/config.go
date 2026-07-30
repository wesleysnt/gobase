package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	Env         string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	LogLevel    string
	LogFormat   string

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

func Load() (*Config, error) {
	// .env is optional; ignore error if file doesn't exist
	_ = godotenv.Load()

	env := getEnv("ENV", "development")

	cfg := &Config{
		Port:        getEnvInt("PORT", 8080),
		Env:         env,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpiry:   getEnvDuration("JWT_EXPIRY", 24*time.Hour),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		LogFormat:   getEnv("LOG_FORMAT", logFormatDefault(env)),

		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	var errs []string
	if cfg.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		errs = append(errs, "JWT_SECRET is required")
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("configuration errors: %v", errs)
	}

	return cfg, nil
}

func logFormatDefault(env string) string {
	if env == "production" {
		return "json"
	}
	return "text"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
