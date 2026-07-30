package database

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/you/gobase/internal/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sqlx.DB, error) {
	driver, err := parseDriver(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Connect(driver, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("database connect: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping: %w", err)
	}

	return db, nil
}

func parseDriver(dsn string) (string, error) {
	if dsn == "" {
		return "", fmt.Errorf("empty DATABASE_URL")
	}

	switch {
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return "postgres", nil
	case strings.HasPrefix(dsn, "mysql://"):
		return "mysql", nil
	case strings.HasPrefix(dsn, "sqlite://"), strings.HasPrefix(dsn, "sqlite3://"),
		strings.HasPrefix(dsn, "file:"):
		return "sqlite3", nil
	default:
		return "", fmt.Errorf("unsupported database URL scheme; supported schemes: postgres://, postgresql://, mysql://, sqlite://, sqlite3://, file:")
	}
}
