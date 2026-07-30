package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Clear relevant env vars (only those with defaults — required fields are set)
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("JWT_EXPIRY")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")
	os.Unsetenv("DB_MAX_OPEN_CONNS")
	os.Unsetenv("DB_MAX_IDLE_CONNS")
	os.Unsetenv("DB_CONN_MAX_LIFETIME")
	// Set required fields so Load() can succeed when testing defaults
	os.Setenv("DATABASE_URL", "postgres://localhost/defaultdb")
	os.Setenv("JWT_SECRET", "default-secret")
	defer os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Port = %d, want 8080", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Errorf("Env = %s, want development", cfg.Env)
	}
	if cfg.JWTExpiry != 24*time.Hour {
		t.Errorf("JWTExpiry = %v, want 24h", cfg.JWTExpiry)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %s, want info", cfg.LogLevel)
	}
	if cfg.DBMaxOpenConns != 25 {
		t.Errorf("DBMaxOpenConns = %d, want 25", cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns != 5 {
		t.Errorf("DBMaxIdleConns = %d, want 5", cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime != 5*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want 5m", cfg.DBConnMaxLifetime)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/db")
	os.Setenv("JWT_SECRET", "super-secret")
	os.Setenv("JWT_EXPIRY", "72h")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "json")
	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "10")
	os.Setenv("DB_CONN_MAX_LIFETIME", "10m")
	defer os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Port != 3000 {
		t.Errorf("Port = %d, want 3000", cfg.Port)
	}
	if cfg.Env != "production" {
		t.Errorf("Env = %s, want production", cfg.Env)
	}
	if cfg.DatabaseURL != "postgres://localhost:5432/db" {
		t.Errorf("DatabaseURL = %s, want postgres://localhost:5432/db", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "super-secret" {
		t.Errorf("JWTSecret = %s, want super-secret", cfg.JWTSecret)
	}
	if cfg.JWTExpiry != 72*time.Hour {
		t.Errorf("JWTExpiry = %v, want 72h", cfg.JWTExpiry)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %s, want debug", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("LogFormat = %s, want json", cfg.LogFormat)
	}
	if cfg.DBMaxOpenConns != 50 {
		t.Errorf("DBMaxOpenConns = %d, want 50", cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns != 10 {
		t.Errorf("DBMaxIdleConns = %d, want 10", cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime != 10*time.Minute {
		t.Errorf("DBConnMaxLifetime = %v, want 10m", cfg.DBConnMaxLifetime)
	}
}

func TestLoadMissingRequired(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	_, err := Load()
	if err == nil {
		t.Error("Load() expected error for missing required fields, got nil")
	}
}

func TestLogFormatDefaultByEnv(t *testing.T) {
	os.Unsetenv("LOG_FORMAT")
	os.Setenv("ENV", "production")
	os.Setenv("DATABASE_URL", "postgres://localhost/db")
	os.Setenv("JWT_SECRET", "secret")
	defer os.Clearenv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("LogFormat = %s, want json (production default)", cfg.LogFormat)
	}
}
