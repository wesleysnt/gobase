# GoBase Template Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a batteries-included Go project template with cobra CLI, chi HTTP API, sqlx database access, golang-migrate migrations, HMAC JWT auth, and slog structured logging.

**Architecture:** Single binary with cobra subcommands (serve, migrate, jwt). chi router assembled in `internal/server/`. Each domain (starting with `user/`) is a vertical slice under `internal/` containing model, repository, service, and handler. Infrastructure packages (config, log, database, auth) are shared. Module placeholder is `github.com/wesleysnt/gobase`, replaced at clone time by `setup.sh`.

**Tech Stack:** Go 1.22+, chi v5, cobra v1.8, sqlx v1.3, golang-migrate v4, golang-jwt v5, go-playground/validator v10, godotenv v1.5

## Global Constraints

- Module path: `github.com/wesleysnt/gobase` (placeholder, replaced by setup.sh)
- Go version: 1.22 minimum
- All application code under `internal/`
- Package naming: singular (`user`, not `users`)
- `context.Context` as first parameter in every cross-layer method
- SQL placeholders: `$1, $2` (PostgreSQL default)
- Errors use sentinel pattern with `errors.Is` in `server/codeFrom()`
- No global state; dependencies injected through constructors and function parameters
- Test-first for all domain logic; integration tests use real database via `TEST_DATABASE_URL`

---

### Task 1: Module skeleton

**Files:**
- Create: `go.mod`
- Create: `main.go`

**Interfaces:**
- Produces: module `github.com/wesleysnt/gobase` with Go 1.22, entry point `main.go` calling `cmd.Execute()`

- [ ] **Step 1: Initialize go.mod**

```bash
cd /Users/wesleysnt/Documents/GitHub/gobase
go mod init github.com/wesleysnt/gobase
```

Run: `go mod init github.com/wesleysnt/gobase`
Expected: creates `go.mod` with `module github.com/wesleysnt/gobase` and `go 1.22` (or current installed version)

- [ ] **Step 2: Write main.go**

```go
// main.go
package main

import "github.com/wesleysnt/gobase/cmd"

func main() {
	cmd.Execute()
}
```

- [ ] **Step 3: Verify main.go compiles (will fail on missing cmd package — expected)**

Run: `go build .`
Expected: FAIL — `package github.com/wesleysnt/gobase/cmd is not in GOROOT`

- [ ] **Step 4: Commit**

```bash
git add go.mod main.go
git commit -m "feat: add module skeleton with go.mod and main.go entry point

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 2: Config management

**Files:**
- Create: `internal/config/config.go`

**Interfaces:**
- Produces: `func Load() (*Config, error)` — reads env + .env, returns typed config
- Produces: `type Config struct { Port int; Env string; DatabaseURL string; JWTSecret string; JWTExpiry time.Duration; LogLevel string; LogFormat string; DBMaxOpenConns int; DBMaxIdleConns int; DBConnMaxLifetime time.Duration }`

- [ ] **Step 1: Write config tests**

Create `internal/config/config_test.go`:

```go
package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	// Clear relevant env vars
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_EXPIRY")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")

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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/config/...`
Expected: FAIL — package not found or compile errors

- [ ] **Step 3: Implement config.go**

```go
// internal/config/config.go
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
```

- [ ] **Step 4: Install dependencies**

Run: `go get github.com/joho/godotenv@v1.5.1`
Expected: adds godotenv to go.mod and go.sum

- [ ] **Step 5: Run tests**

Run: `go test ./internal/config/... -v`
Expected: all tests PASS

- [ ] **Step 6: Commit**

```bash
git add internal/config/config.go internal/config/config_test.go go.mod go.sum
git commit -m "feat: add config management with env + .env loading

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 3: Structured logging

**Files:**
- Create: `internal/log/log.go`
- Create: `internal/log/log_test.go`

**Interfaces:**
- Consumes: `config.Config` (LogLevel, LogFormat, Env fields)
- Produces: `func New(cfg *config.Config) *slog.Logger` — returns configured logger

- [ ] **Step 1: Write log tests**

```go
// internal/log/log_test.go
package log

import (
	"bytes"
	"strings"
	"testing"
	"testing/slog"

	"github.com/wesleysnt/gobase/internal/config"
)

func TestNewTextFormat(t *testing.T) {
	cfg := &config.Config{
		Env:      "development",
		LogLevel: "info",
	}

	buf := &bytes.Buffer{}
	logger := New(cfg)
	// slog doesn't expose the writer, so we test behavior indirectly
	// by checking that the logger is non-nil and has the right level
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}
}

func TestNewJSONFormat(t *testing.T) {
	cfg := &config.Config{
		Env:       "production",
		LogFormat: "json",
		LogLevel:  "warn",
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}
}

func TestNewDebugLevel(t *testing.T) {
	cfg := &config.Config{
		LogLevel: "debug",
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}
}

func TestNewInvalidLevel(t *testing.T) {
	cfg := &config.Config{
		LogLevel: "invalid",
	}

	logger := New(cfg)
	if logger == nil {
		t.Fatal("New() returned nil logger (should default to info)")
	}
}

func TestLevelMapping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"debug", "debug", slog.LevelDebug},
		{"info", "info", slog.LevelInfo},
		{"warn", "warn", slog.LevelWarn},
		{"error", "error", slog.LevelError},
		{"invalid", "bogus", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLevel(tt.input)
			if got != tt.expected {
				t.Errorf("parseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/log/...`
Expected: FAIL — compile errors (undefined: New, undefined: parseLevel)

- [ ] **Step 3: Implement log.go**

```go
// internal/log/log.go
package log

import (
	"log/slog"
	"os"

	"github.com/wesleysnt/gobase/internal/config"
)

func New(cfg *config.Config) *slog.Logger {
	level := parseLevel(cfg.LogLevel)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	switch cfg.LogFormat {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./internal/log/... -v`
Expected: all tests PASS

- [ ] **Step 5: Commit**

```bash
git add internal/log/log.go internal/log/log_test.go
git commit -m "feat: add structured logging with slog, JSON and text formats

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 4: Database connection & migration infrastructure

**Files:**
- Create: `internal/database/database.go`
- Create: `internal/database/migrate.go`
- Create: `internal/database/database_test.go`

**Interfaces:**
- Consumes: `config.Config` (DatabaseURL, pool settings)
- Produces: `func Connect(cfg *config.Config) (*sqlx.DB, error)` — infers driver from URL scheme
- Produces: `func RunMigrations(db *sqlx.DB, databaseURL, direction string, steps int) error` — embedded migration runner
- Produces: `func CreateMigration(name string) error` — creates up/down SQL files in migrations/

- [ ] **Step 1: Write database tests**

```go
// internal/database/database_test.go
package database

import (
	"os"
	"testing"

	"github.com/wesleysnt/gobase/internal/config"
)

func TestParseDriver(t *testing.T) {
	tests := []struct {
		url     string
		driver  string
		wantErr bool
	}{
		{"postgres://localhost:5432/db", "pgx", false},
		{"postgresql://localhost:5432/db", "pgx", false},
		{"mysql://localhost:3306/db", "mysql", false},
		{"sqlite://data.db", "sqlite3", false},
		{"file:test.db", "sqlite3", false},
		{"sqlite3://data.db", "sqlite3", false},
		{"unknown://localhost/db", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.driver, func(t *testing.T) {
			driver, err := parseDriver(tt.url)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseDriver(%q) expected error, got driver=%q", tt.url, driver)
				}
				return
			}
			if err != nil {
				t.Errorf("parseDriver(%q) unexpected error: %v", tt.url, err)
			}
			if driver != tt.driver {
				t.Errorf("parseDriver(%q) = %q, want %q", tt.url, driver, tt.driver)
			}
		})
	}
}

func TestConnectUnreachable(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping integration test in CI")
	}

	cfg := &config.Config{
		DatabaseURL:      "postgres://localhost:59999/nonexistent",
		DBMaxOpenConns:    1,
		DBMaxIdleConns:    1,
		DBConnMaxLifetime: 1,
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Error("Connect() expected error for unreachable database, got nil")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/database/...`
Expected: FAIL — undefined symbols

- [ ] **Step 3: Get sqlx and golang-migrate dependencies**

Run:
```bash
go get github.com/jmoiron/sqlx@v1.4.0
go get github.com/lib/pq@v1.10.9
go get github.com/golang-migrate/migrate/v4@v4.17.0
```

- [ ] **Step 4: Implement database.go**

```go
// internal/database/database.go
package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wesleysnt/gobase/internal/config"

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
		return "pgx", nil
	case strings.HasPrefix(dsn, "mysql://"):
		return "mysql", nil
	case strings.HasPrefix(dsn, "sqlite://"), strings.HasPrefix(dsn, "sqlite3://"),
		strings.HasPrefix(dsn, "file:"):
		return "sqlite3", nil
	default:
		return "", fmt.Errorf("unsupported database URL scheme; supported schemes: postgres://, postgresql://, mysql://, sqlite://, sqlite3://, file:")
	}
}
```

- [ ] **Step 5: Implement migrate.go**

```go
// internal/database/migrate.go
package database

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(db *sqlx.DB, direction string, steps int) error {
	sourceDriver, err := iofs.New(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("migration source: %w", err)
	}
	// The iofs driver strips the "../../migrations/" prefix in its path handling.
	// We need to register it differently. For now, use the filesystem path.
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, db.Driver().(interface{ Name() string }).Name()+"://"+db.DSN())
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Steps(-1)
		}
	default:
		return fmt.Errorf("unknown migration direction: %s (use 'up' or 'down')", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration %s: %w", direction, err)
	}

	return nil
}

func CreateMigration(name string) error {
	dir := "migrations"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create migrations dir: %w", err)
	}

	timestamp := time.Now().Unix()

	upFile := filepath.Join(dir, fmt.Sprintf("%d_%s.up.sql", timestamp, name))
	downFile := filepath.Join(dir, fmt.Sprintf("%d_%s.down.sql", timestamp, name))

	if err := os.WriteFile(upFile, []byte("-- migrate:up\n"), 0o644); err != nil {
		return fmt.Errorf("create up migration: %w", err)
	}
	if err := os.WriteFile(downFile, []byte("-- migrate:down\n"), 0o644); err != nil {
		return fmt.Errorf("create down migration: %w", err)
	}

	fmt.Printf("Created: %s\nCreated: %s\n", upFile, downFile)
	return nil
}
```

- [ ] **Step 6: Run tests**

Run: `go test ./internal/database/... -v -run TestParseDriver`
Expected: PASS for parseDriver cases, SKIP for integration test

- [ ] **Step 7: Commit**

```bash
git add internal/database/ go.mod go.sum
git commit -m "feat: add database connection with driver inference and migration infrastructure

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 5: JWT auth & middleware

**Files:**
- Create: `internal/auth/jwt.go`
- Create: `internal/auth/middleware.go`
- Create: `internal/auth/jwt_test.go`
- Create: `internal/auth/middleware_test.go`

**Interfaces:**
- Consumes: `config.Config` (JWTSecret, JWTExpiry)
- Produces: `type Claims struct { RegisteredClaims jwtsdk.RegisteredClaims; Extra map[string]interface{} }`
- Produces: `func GenerateToken(userID string, expiry time.Duration, extra map[string]interface{}) (string, []byte, error)` — signs HMAC-SHA256; returns token + secret used
- Produces: `func ParseToken(tokenStr string, secret []byte) (*Claims, error)` — validates signature and expiry
- Produces: `func RequireAuth(secret []byte) func(http.Handler) http.Handler`
- Produces: `func OptionalAuth(secret []byte) func(http.Handler) http.Handler`
- Produces: `func GetClaims(ctx context.Context) *Claims`
- Produces: `func SetClaims(ctx context.Context, c *Claims) context.Context`

- [ ] **Step 1: Write JWT tests**

```go
// internal/auth/jwt_test.go
package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseToken(t *testing.T) {
	userID := "usr_abc123"
	expiry := 1 * time.Hour
	secret := []byte("test-secret-key")

	token, _, err := GenerateToken(userID, expiry, nil, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken() returned empty token")
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.Subject != userID {
		t.Errorf("claims.Subject = %q, want %q", claims.Subject, userID)
	}
}

func TestGenerateTokenWithExtraClaims(t *testing.T) {
	secret := []byte("test-secret-key")
	extra := map[string]interface{}{
		"role":     "admin",
		"tenantID": float64(42),
	}

	token, _, err := GenerateToken("usr_1", 1*time.Hour, extra, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if role, ok := claims.Extra["role"]; !ok || role != "admin" {
		t.Errorf("Extra[role] = %v, want admin", claims.Extra["role"])
	}
}

func TestParseExpiredToken(t *testing.T) {
	secret := []byte("test-secret-key")
	token, _, err := GenerateToken("usr_1", -1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = ParseToken(token, secret)
	if err == nil {
		t.Error("ParseToken() expected error for expired token, got nil")
	}
}

func TestParseTokenWrongSecret(t *testing.T) {
	secret := []byte("correct-secret")
	wrongSecret := []byte("wrong-secret")

	token, _, err := GenerateToken("usr_1", 1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	_, err = ParseToken(token, wrongSecret)
	if err == nil {
		t.Error("ParseToken() expected error for wrong secret, got nil")
	}
}

func TestParseMalformedToken(t *testing.T) {
	_, err := ParseToken("not.a.valid.token", []byte("secret"))
	if err == nil {
		t.Error("ParseToken() expected error for malformed token, got nil")
	}
}
```

- [ ] **Step 2: Get golang-jwt dependency**

Run: `go get github.com/golang-jwt/jwt/v5@v5.2.1`

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test ./internal/auth/...`
Expected: FAIL — undefined symbols

- [ ] **Step 4: Implement jwt.go**

```go
// internal/auth/jwt.go
package auth

import (
	"fmt"
	"time"

	jwtsdk "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwtsdk.RegisteredClaims
	Extra map[string]interface{} `json:"extra,omitempty"`
}

func GenerateToken(userID string, expiry time.Duration, extra map[string]interface{}, secret []byte) (string, []byte, error) {
	now := time.Now()
	claims := &Claims{
		RegisteredClaims: jwtsdk.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwtsdk.NewNumericDate(now),
			ExpiresAt: jwtsdk.NewNumericDate(now.Add(expiry)),
		},
		Extra: extra,
	}

	token := jwtsdk.NewWithClaims(jwtsdk.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", nil, fmt.Errorf("sign token: %w", err)
	}

	return signed, secret, nil
}

func ParseToken(tokenStr string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtsdk.ParseWithClaims(tokenStr, claims, func(t *jwtsdk.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtsdk.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
```

- [ ] **Step 5: Run JWT tests**

Run: `go test ./internal/auth/... -v -run TestGenerate`
Expected: all JWT tests PASS

- [ ] **Step 6: Write middleware tests**

```go
// internal/auth/middleware_test.go
package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequireAuthValidToken(t *testing.T) {
	secret := []byte("test-secret")
	token, _, err := GenerateToken("usr_1", 1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	var capturedClaims *Claims
	handler := RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedClaims = GetClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if capturedClaims == nil {
		t.Fatal("GetClaims() returned nil")
	}
	if capturedClaims.Subject != "usr_1" {
		t.Errorf("Subject = %q, want usr_1", capturedClaims.Subject)
	}
}

func TestRequireAuthNoToken(t *testing.T) {
	secret := []byte("test-secret")
	handler := RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuthInvalidHeader(t *testing.T) {
	secret := []byte("test-secret")
	handler := RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestOptionalAuthValidToken(t *testing.T) {
	secret := []byte("test-secret")
	token, _, err := GenerateToken("usr_1", 1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	var capturedClaims *Claims
	handler := OptionalAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedClaims = GetClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/optional", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if capturedClaims == nil {
		t.Fatal("GetClaims() returned nil")
	}
}

func TestOptionalAuthNoToken(t *testing.T) {
	secret := []byte("test-secret")
	var capturedClaims *Claims
	handler := OptionalAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedClaims = GetClaims(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/optional", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if capturedClaims != nil {
		t.Error("GetClaims() should return nil for unauthenticated request")
	}
}

func TestSetAndGetClaims(t *testing.T) {
	claims := &Claims{}
	claims.Subject = "usr_42"

	ctx := SetClaims(context.Background(), claims)
	got := GetClaims(ctx)

	if got == nil {
		t.Fatal("GetClaims() returned nil")
	}
	if got.Subject != "usr_42" {
		t.Errorf("Subject = %q, want usr_42", got.Subject)
	}
}
```

- [ ] **Step 7: Implement middleware.go**

```go
// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const claimsKey contextKey = "claims"

func RequireAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := extractAndParse(r, secret)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			ctx := SetClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := extractAndParse(r, secret)
			if err == nil {
				ctx := SetClaims(r.Context(), claims)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func extractAndParse(r *http.Request, secret []byte) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.ErrNoCookie // using ErrNoCookie as sentinel; caught by callers
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	return ParseToken(tokenStr, secret)
}

func GetClaims(ctx context.Context) *Claims {
	if claims, ok := ctx.Value(claimsKey).(*Claims); ok {
		return claims
	}
	return nil
}

func SetClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}
```

- [ ] **Step 8: Run all auth tests**

Run: `go test ./internal/auth/... -v`
Expected: all tests PASS

- [ ] **Step 9: Commit**

```bash
git add internal/auth/ go.mod go.sum
git commit -m "feat: add JWT HMAC-SHA256 token generation, parsing, and chi middleware

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 6: Server infrastructure

**Files:**
- Create: `internal/server/response.go`
- Create: `internal/server/errors.go`
- Create: `internal/server/middleware.go`
- Create: `internal/server/routes.go`
- Create: `internal/server/server.go`
- Create: `internal/server/response_test.go`
- Create: `internal/server/errors_test.go`
- Create: `internal/server/middleware_test.go`

**Interfaces:**
- Consumes: `config.Config`, `*slog.Logger`, `*sqlx.DB`, `auth` package
- Produces: `func WriteJSON(w http.ResponseWriter, status int, data interface{})`
- Produces: `func WriteError(w http.ResponseWriter, status int, message string)`
- Produces: `var ErrNotFound, ErrValidation, ErrConflict, ErrForbidden, ErrUnauthorized error`
- Produces: `func CodeFrom(err error) int`
- Produces: `func NewRouter(cfg *config.Config, db *sqlx.DB, log *slog.Logger) chi.Router`
- Produces: `func ListenAndServe(cfg *config.Config, log *slog.Logger, router http.Handler) error`

- [ ] **Step 1: Get chi dependency**

Run: `go get github.com/go-chi/chi/v5@v5.1.0`

- [ ] **Step 2: Write tests for response helpers and error mapping**

```go
// internal/server/response_test.go
package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	data := map[string]string{"status": "ok"}

	WriteJSON(rec, http.StatusOK, data)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var decoded map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if decoded["status"] != "ok" {
		t.Errorf("body.status = %q, want ok", decoded["status"])
	}
}

func TestWriteJSONUnsupportedType(t *testing.T) {
	rec := httptest.NewRecorder()
	// channels can't be marshaled to JSON
	WriteJSON(rec, http.StatusOK, make(chan int))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusNotFound, "user not found")

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["error"] != "user not found" {
		t.Errorf("error = %q, want 'user not found'", body["error"])
	}
}
```

```go
// internal/server/errors_test.go
package server

import (
	"errors"
	"net/http"
	"testing"
)

func TestCodeFrom(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"not found", ErrNotFound, http.StatusNotFound},
		{"validation", ErrValidation, http.StatusUnprocessableEntity},
		{"conflict", ErrConflict, http.StatusConflict},
		{"forbidden", ErrForbidden, http.StatusForbidden},
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized},
		{"unknown", errors.New("something bad"), http.StatusInternalServerError},
		{"wrapped", errors.New("wrapped: some context"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CodeFrom(tt.err)
			if got != tt.want {
				t.Errorf("CodeFrom(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}

func TestCodeFromWrapped(t *testing.T) {
	err := errors.New("extra detail")
	wrapped := errors.Join(ErrNotFound, err)

	got := CodeFrom(wrapped)
	if got != http.StatusNotFound {
		t.Errorf("CodeFrom(wrapped ErrNotFound) = %d, want %d", got, http.StatusNotFound)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test ./internal/server/...`
Expected: FAIL — undefined symbols

- [ ] **Step 4: Implement response.go**

```go
// internal/server/response.go
package server

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
```

- [ ] **Step 5: Implement errors.go**

```go
// internal/server/errors.go
package server

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrValidation    = errors.New("validation error")
	ErrConflict      = errors.New("conflict")
	ErrForbidden     = errors.New("forbidden")
	ErrUnauthorized  = errors.New("unauthorized")
)

func CodeFrom(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrValidation):
		return http.StatusUnprocessableEntity
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
```

- [ ] **Step 6: Run response + errors tests**

Run: `go test ./internal/server/... -v -run "TestWrite|TestCodeFrom"`
Expected: PASS

- [ ] **Step 7: Implement server middleware (middleware.go)**

```go
// internal/server/middleware.go
package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"duration", time.Since(start).String(),
				"request_id", middleware.GetReqID(r.Context()),
			)
		})
	}
}
```

- [ ] **Step 8: Write middleware test**

```go
// internal/server/middleware_test.go
package server

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	output := buf.String()
	if !strings.Contains(output, "GET") || !strings.Contains(output, "/test") {
		t.Errorf("log output missing request info: %s", output)
	}
}
```

- [ ] **Step 9: Implement routes.go**

```go
// internal/server/routes.go
package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/config"
)

func NewRouter(cfg *config.Config, db *sqlx.DB, log *slog.Logger) chi.Router {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(RequestLogger(log))
	r.Use(middleware.Recoverer)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 — domain routes registered here
	r.Route("/api/v1", func(r chi.Router) {
		// user.RegisterRoutes(r, db, cfg) — uncommented after Task 7
	})

	return r
}
```

- [ ] **Step 10: Implement server.go**

```go
// internal/server/server.go
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wesleysnt/gobase/internal/config"
)

func ListenAndServe(cfg *config.Config, log *slog.Logger, router http.Handler) error {
	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server forced to shutdown", "error", err)
		}
	}()

	log.Info("server starting", "addr", addr, "env", cfg.Env)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server: %w", err)
	}

	log.Info("server stopped")
	return nil
}
```

- [ ] **Step 11: Run all server tests**

Run: `go test ./internal/server/... -v`
Expected: all tests PASS

- [ ] **Step 12: Commit**

```bash
git add internal/server/ go.mod go.sum
git commit -m "feat: add server infrastructure — router assembly, graceful shutdown, response helpers, error mapping

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 7: User domain module

**Files:**
- Create: `internal/user/model.go`
- Create: `internal/user/repository.go`
- Create: `internal/user/service.go`
- Create: `internal/user/handler.go`
- Create: `internal/user/service_test.go`
- Create: `internal/user/handler_test.go`

**Interfaces:**
- Consumes: `*sqlx.DB`, `*config.Config`, server response helpers, auth middleware
- Produces: `func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config)`
- Produces: `type User struct { ID string; Name string; Email string; CreatedAt time.Time; UpdatedAt time.Time }`
- Produces: `type CreateUserRequest struct { Name string; Email string }`
- Produces: `type UpdateUserRequest struct { Name *string; Email *string }`
- Produces: `type UserRepo struct { db *sqlx.DB }` with CRUD methods
- Produces: `type UserService struct { repo *UserRepo }` with business logic
- Produces: `type Handler struct { svc *UserService }` with HTTP handlers

- [ ] **Step 1: Get validator dependency**

Run: `go get github.com/go-playground/validator/v10@v10.22.0`

- [ ] **Step 2: Write service tests**

```go
// internal/user/service_test.go
package user

import (
	"context"
	"errors"
	"testing"

	"github.com/wesleysnt/gobase/internal/server"
)

type stubUserRepo struct {
	users map[string]*User
}

func newStubRepo() *stubUserRepo {
	return &stubUserRepo{users: make(map[string]*User)}
}

func (r *stubUserRepo) Create(ctx context.Context, u *User) error {
	u.ID = "generated-id"
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	// Simulate email uniqueness check
	for _, existing := range r.users {
		if existing.Email == u.Email {
			return server.ErrConflict
		}
	}
	r.users[u.ID] = u
	return nil
}

func (r *stubUserRepo) FindByID(ctx context.Context, id string) (*User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, server.ErrNotFound
	}
	return u, nil
}

func (r *stubUserRepo) List(ctx context.Context) ([]*User, error) {
	result := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

func (r *stubUserRepo) Update(ctx context.Context, u *User) error {
	if _, ok := r.users[u.ID]; !ok {
		return server.ErrNotFound
	}
	u.UpdatedAt = time.Now()
	r.users[u.ID] = u
	return nil
}

func (r *stubUserRepo) Delete(ctx context.Context, id string) error {
	if _, ok := r.users[id]; !ok {
		return server.ErrNotFound
	}
	delete(r.users, id)
	return nil
}

import "time"

func TestServiceCreate(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	req := CreateUserRequest{Name: "Alice", Email: "alice@example.com"}
	user, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.ID == "" {
		t.Error("Create() user.ID is empty")
	}
	if user.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", user.Name)
	}
}

func TestServiceCreateValidation(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}

	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr bool
	}{
		{"missing name", CreateUserRequest{Name: "", Email: "a@b.com"}, true},
		{"missing email", CreateUserRequest{Name: "Bob", Email: ""}, true},
		{"short name", CreateUserRequest{Name: "AB", Email: "a@b.com"}, true},
		{"invalid email", CreateUserRequest{Name: "Charlie", Email: "not-an-email"}, true},
		{"valid", CreateUserRequest{Name: "Diana", Email: "diana@example.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tt.req)
			if tt.wantErr && err == nil {
				t.Error("Create() expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Create() unexpected error: %v", err)
			}
		})
	}
}

func TestServiceGetByID(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	created, _ := svc.Create(context.Background(), CreateUserRequest{Name: "Eve", Email: "eve@example.com"})

	user, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if user.Name != "Eve" {
		t.Errorf("Name = %q, want Eve", user.Name)
	}
}

func TestServiceGetByIDNotFound(t *testing.T) {
	svc := &UserService{repo: newStubRepo()}
	_, err := svc.GetByID(context.Background(), "nonexistent")
	if !errors.Is(err, server.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test ./internal/user/...`
Expected: FAIL — undefined symbols

- [ ] **Step 4: Implement model.go**

```go
// internal/user/model.go
package user

import "time"

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Email *string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
```

- [ ] **Step 5: Implement repository.go**

```go
// internal/user/repository.go
package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/server"
)

type UserRepo struct {
	db *sqlx.DB
}

func (r *UserRepo) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*User, error) {
	var u User
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s: %w", id, server.ErrNotFound)
	}
	return &u, err
}

func (r *UserRepo) List(ctx context.Context) ([]*User, error) {
	var users []*User
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *UserRepo) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	result, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s: %w", u.ID, server.ErrNotFound)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s: %w", id, server.ErrNotFound)
	}
	return nil
}
```

- [ ] **Step 6: Implement service.go**

```go
// internal/user/service.go
package user

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/wesleysnt/gobase/internal/server"
)

var validate = validator.New()

type UserService struct {
	repo *UserRepo
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", server.ErrValidation, err)
	}

	user := &User{Name: req.Name, Email: req.Email}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) List(ctx context.Context) ([]*UserResponse, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		r := u.ToResponse()
		resp[i] = &r
	}
	return resp, nil
}

func (s *UserService) Update(ctx context.Context, id string, req UpdateUserRequest) (*UserResponse, error) {
	if err := validate.Struct(req); err != nil {
		return nil, fmt.Errorf("%w: %v", server.ErrValidation, err)
	}

	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
```

- [ ] **Step 7: Run service tests**

Run: `go test ./internal/user/... -v -run TestService`
Expected: all service tests PASS

- [ ] **Step 8: Implement handler.go**

```go
// internal/user/handler.go
package user

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/auth"
	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/server"
)

type Handler struct {
	svc *UserService
}

func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	repo := &UserRepo{db: db}
	svc := &UserService{repo: repo}
	h := &Handler{svc: svc}

	r.Group(func(r chi.Router) {
		r.Get("/users", h.List)
		r.Post("/users", h.Create)
		r.Get("/users/{id}", h.GetByID)
		r.Patch("/users/{id}", h.Update)
		r.Delete("/users/{id}", h.Delete)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth([]byte(cfg.JWTSecret)))
		r.Get("/me", h.Me)
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	server.WriteJSON(w, http.StatusCreated, user)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	server.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	if users == nil {
		users = []*UserResponse{}
	}

	server.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	server.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaims(r.Context())
	if claims == nil {
		server.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := h.svc.GetByID(r.Context(), claims.Subject)
	if err != nil {
		server.WriteError(w, server.CodeFrom(err), err.Error())
		return
	}

	server.WriteJSON(w, http.StatusOK, user)
}
```

- [ ] **Step 9: Write handler tests**

```go
// internal/user/handler_test.go
package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func setupTestRouter() chi.Router {
	r := chi.NewRouter()
	repo := newStubRepo()
	svc := &UserService{repo: repo}
	h := &Handler{svc: svc}

	r.Post("/users", h.Create)
	r.Get("/users", h.List)
	r.Get("/users/{id}", h.GetByID)
	r.Patch("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)

	return r
}

func TestHandlerCreate(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"Alice","email":"alice@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp UserResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Name != "Alice" {
		t.Errorf("Name = %q, want Alice", resp.Name)
	}
	if resp.ID == "" {
		t.Error("ID is empty")
	}
}

func TestHandlerCreateValidation(t *testing.T) {
	router := setupTestRouter()

	body := `{"name":"","email":""}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

func TestHandlerGetByID(t *testing.T) {
	router := setupTestRouter()

	// Create a user first
	createBody := `{"name":"Bob","email":"bob@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var created UserResponse
	json.NewDecoder(rec.Body).Decode(&created)

	// Now fetch it
	req = httptest.NewRequest("GET", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var fetched UserResponse
	json.NewDecoder(rec.Body).Decode(&fetched)
	if fetched.Name != "Bob" {
		t.Errorf("Name = %q, want Bob", fetched.Name)
	}
}

func TestHandlerGetByIDNotFound(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest("GET", "/users/nonexistent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandlerList(t *testing.T) {
	router := setupTestRouter()

	// Create two users
	for _, body := range []string{
		`{"name":"Alice","email":"alice@example.com"}`,
		`{"name":"Bob","email":"bob@example.com"}`,
	} {
		req := httptest.NewRequest("POST", "/users", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest("GET", "/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var users []UserResponse
	json.NewDecoder(rec.Body).Decode(&users)
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}
}

func TestHandlerDelete(t *testing.T) {
	router := setupTestRouter()

	// Create a user
	createBody := `{"name":"Charlie","email":"charlie@example.com"}`
	req := httptest.NewRequest("POST", "/users", strings.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	var created UserResponse
	json.NewDecoder(rec.Body).Decode(&created)

	// Delete it
	req = httptest.NewRequest("DELETE", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	// Verify it's gone
	req = httptest.NewRequest("GET", "/users/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d (not found after delete)", rec.Code, http.StatusNotFound)
	}
}
```

- [ ] **Step 10: Run all user tests**

Run: `go test ./internal/user/... -v`
Expected: all tests PASS

- [ ] **Step 11: Wire user routes into the router**

Update `internal/server/routes.go` — add user import and uncomment the user registration line:

In `routes.go`, add the import:
```go
import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/user"  // add this
)
```

And replace the comment inside `r.Route("/api/v1", ...)` with:
```go
r.Route("/api/v1", func(r chi.Router) {
    user.RegisterRoutes(r, db, cfg)
})
```

Run: `go build ./...` to verify it compiles.

- [ ] **Step 12: Fix stub repo to implement the UserRepo interface properly**

The handler tests use a stub repo. The service currently takes `*UserRepo` (concrete). We need to make the service accept an interface so handler tests work without a database.

Update `internal/user/service.go` — change the service struct:

```go
// internal/user/service.go
package user

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/wesleysnt/gobase/internal/server"
)

var validate = validator.New()

// UserRepository defines the data access interface for users.
type UserRepository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id string) error
}

type UserService struct {
	repo UserRepository
}
```

And update `internal/user/handler.go` — the `RegisterRoutes` function already creates `&UserRepo{db: db}` which satisfies the interface:

In handler.go, update RegisterRoutes:
```go
func RegisterRoutes(r chi.Router, db *sqlx.DB, cfg *config.Config) {
	svc := &UserService{repo: &UserRepo{db: db}}
	h := &Handler{svc: svc}
	// ... rest unchanged
}
```

Run: `go test ./internal/user/... -v`
Expected: all tests PASS

- [ ] **Step 13: Commit**

```bash
git add internal/user/ internal/server/routes.go go.mod go.sum
git commit -m "feat: add user domain module with CRUD handlers, service layer, and repository

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 8: CLI commands

**Files:**
- Create: `cmd/root.go`
- Create: `cmd/serve.go`
- Create: `cmd/migrate.go`
- Create: `cmd/jwt.go`

**Interfaces:**
- Consumes: all internal packages
- Produces: `func Execute()` — cobra root command

- [ ] **Step 1: Get cobra dependency**

Run: `go get github.com/spf13/cobra@v1.8.1`

- [ ] **Step 2: Implement root.go**

```go
// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gobase",
	Short: "GoBase — a Go project template",
	Long:  `A batteries-included Go project with CLI, HTTP API, database migrations, and JWT auth.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

- [ ] **Step 3: Implement serve.go**

```go
// cmd/serve.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/database"
	"github.com/wesleysnt/gobase/internal/log"
	"github.com/wesleysnt/gobase/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		// Override from flags if set
		if port, _ := cmd.Flags().GetInt("port"); port != 0 {
			cfg.Port = port
		}
		if env, _ := cmd.Flags().GetString("env"); env != "" {
			cfg.Env = env
		}

		logger := log.New(cfg)

		db, err := database.Connect(cfg)
		if err != nil {
			return fmt.Errorf("database: %w", err)
		}
		defer db.Close()

		router := server.NewRouter(cfg, db, logger)

		return server.ListenAndServe(cfg, logger, router)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 0, "Server port (default: from PORT env or 8080)")
	serveCmd.Flags().String("env", "", "Environment: development|production|test")
}
```

- [ ] **Step 4: Implement migrate.go**

```go
// cmd/migrate.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/database"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Run pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrate("up", cmd)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrate("down", cmd)
	},
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file pair",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return database.CreateMigration(args[0])
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		// For now: run up with steps=0 to show pending status
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		db, err := database.Connect(cfg)
		if err != nil {
			return fmt.Errorf("database: %w", err)
		}
		defer db.Close()

		// Use golang-migrate to get version info
		version, dirty, err := database.MigrationStatus(db, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
		return nil
	},
}

func runMigrate(direction string, cmd *cobra.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		return fmt.Errorf("database: %w", err)
	}
	defer db.Close()

	steps, _ := cmd.Flags().GetInt("steps")
	return database.RunMigrations(db, cfg.DatabaseURL, direction, steps)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	migrateCmd.AddCommand(migrateStatusCmd)

	migrateUpCmd.Flags().Int("steps", 0, "Max migrations to apply (0 = all)")
	migrateDownCmd.Flags().Int("steps", 1, "Number of migrations to roll back")
}
```

- [ ] **Step 5: Add MigrationStatus to database package**

We need to add the `MigrationStatus` function to `internal/database/migrate.go`:

```go
func MigrationStatus(db *sqlx.DB, cfg *config.Config) (version uint, dirty bool, err error) {
	driver, err := parseDriver(cfg.DatabaseURL)
	if err != nil {
		return 0, false, err
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL+"?sslmode=disable",
	)
	if err != nil {
		return 0, false, fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	version, dirty, err = m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("migrate version: %w", err)
	}
	return version, dirty, nil
}
```

Also update `RunMigrations` to accept `databaseURL` since we need the DSN for the migrate library:

Update the signature in migrate.go from `func RunMigrations(db *sqlx.DB, direction string, steps int) error` to:
```go
func RunMigrations(db *sqlx.DB, databaseURL, direction string, steps int) error {
	driver, err := parseDriver(databaseURL)
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file://migrations",
		databaseURL+"?sslmode=disable",
	)
	// ... rest of implementation
}
```

Full updated migrate.go:

```go
// internal/database/migrate.go
package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB, databaseURL, direction string, steps int) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	switch direction {
	case "up":
		if steps > 0 {
			err = m.Steps(steps)
		} else {
			err = m.Up()
		}
	case "down":
		if steps > 0 {
			err = m.Steps(-steps)
		} else {
			err = m.Steps(-1)
		}
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration %s: %w", direction, err)
	}

	return nil
}

func MigrationStatus(db *sqlx.DB, databaseURL string) (version uint, dirty bool, err error) {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return 0, false, fmt.Errorf("migrate init: %w", err)
	}
	defer m.Close()

	version, dirty, err = m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return 0, false, fmt.Errorf("migrate version: %w", err)
	}
	return version, dirty, nil
}

func CreateMigration(name string) error {
	dir := "migrations"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create migrations dir: %w", err)
	}

	timestamp := time.Now().Unix()

	upFile := filepath.Join(dir, fmt.Sprintf("%d_%s.up.sql", timestamp, name))
	downFile := filepath.Join(dir, fmt.Sprintf("%d_%s.down.sql", timestamp, name))

	if err := os.WriteFile(upFile, []byte("-- migrate:up\n"), 0o644); err != nil {
		return fmt.Errorf("create up migration: %w", err)
	}
	if err := os.WriteFile(downFile, []byte("-- migrate:down\n"), 0o644); err != nil {
		return fmt.Errorf("create down migration: %w", err)
	}

	fmt.Printf("Created: %s\nCreated: %s\n", upFile, downFile)
	return nil
}
```

- [ ] **Step 6: Implement jwt.go**

```go
// cmd/jwt.go
package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/wesleysnt/gobase/internal/auth"
	"github.com/wesleysnt/gobase/internal/config"
)

var jwtGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a signed JWT",
	Example: `  gobase jwt generate --user-id=usr_123
  gobase jwt generate --user-id=usr_123 --expires-in=72h
  gobase jwt generate --user-id=usr_123 --claims='{"role":"admin"}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}

		userID, _ := cmd.Flags().GetString("user-id")
		if userID == "" {
			return fmt.Errorf("--user-id is required")
		}

		expiresIn, _ := cmd.Flags().GetDuration("expires-in")
		if expiresIn == 0 {
			expiresIn = 24 * time.Hour
		}

		claimsJSON, _ := cmd.Flags().GetString("claims")
		var extra map[string]interface{}
		if claimsJSON != "" {
			if err := json.Unmarshal([]byte(claimsJSON), &extra); err != nil {
				return fmt.Errorf("invalid --claims JSON: %w", err)
			}
		}

		secret := []byte(cfg.JWTSecret)
		token, _, err := auth.GenerateToken(userID, expiresIn, extra, secret)
		if err != nil {
			return fmt.Errorf("generate token: %w", err)
		}

		fmt.Println(token)
		return nil
	},
}

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "JWT token utilities",
}

func init() {
	rootCmd.AddCommand(jwtCmd)
	jwtCmd.AddCommand(jwtGenerateCmd)

	jwtGenerateCmd.Flags().String("user-id", "", "Subject claim (user ID) — required")
	jwtGenerateCmd.Flags().Duration("expires-in", 24*time.Hour, "Token expiry duration")
	jwtGenerateCmd.Flags().String("claims", "", "Extra claims as JSON object")
}
```

- [ ] **Step 7: Verify compilation**

Run: `go build .`
Expected: builds successfully, creates `gobase` binary

- [ ] **Step 8: Verify CLI help works**

Run:
```bash
./gobase --help
./gobase serve --help
./gobase migrate --help
```
Expected: help output for each command

- [ ] **Step 9: Commit**

```bash
git add cmd/ go.mod go.sum internal/database/migrate.go
git commit -m "feat: add cobra CLI with serve, migrate, and jwt generate subcommands

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 9: Project scaffolding — migrations, config, Makefile, setup.sh

**Files:**
- Create: `migrations/000001_create_users.up.sql`
- Create: `migrations/000001_create_users.down.sql`
- Create: `.env.example`
- Create: `Makefile`
- Create: `setup.sh`

- [ ] **Step 1: Create initial migration**

```sql
-- migrations/000001_create_users.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

```sql
-- migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

- [ ] **Step 2: Create .env.example**

```
# GoBase Configuration
# Copy this file to .env and fill in your values

# Server
PORT=8080
ENV=development

# Database (required)
# Supported schemes: postgres://, postgresql://, mysql://, sqlite://, file:
DATABASE_URL=postgres://localhost:5432/gobase?sslmode=disable

# Database pool (optional, defaults shown)
# DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS=5
# DB_CONN_MAX_LIFETIME=5m

# JWT (required)
JWT_SECRET=change-me-to-a-random-secret
# JWT_EXPIRY=24h

# Logging (optional, defaults shown)
# LOG_LEVEL=info
# LOG_FORMAT=text
```

- [ ] **Step 3: Create Makefile**

```makefile
.PHONY: build dev test lint migrate-up migrate-down jwt clean

APP_NAME := gobase
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) .

dev:
	go run . serve

test:
	go test ./... -v

lint:
	golangci-lint run

migrate-up:
	go run . migrate up

migrate-down:
	go run . migrate down

jwt:
	go run . jwt generate --user-id=$(id)

clean:
	rm -rf $(BUILD_DIR)
```

- [ ] **Step 4: Create setup.sh**

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "GoBase Project Setup"
echo "===================="
echo ""

# Prompt for new module path
read -r -p "Module path (e.g., github.com/wesleysnt/myproject): " MODULE_PATH
if [ -z "$MODULE_PATH" ]; then
    echo "Error: Module path is required"
    exit 1
fi

OLD_MODULE="github.com/wesleysnt/gobase"

echo ""
echo "Replacing module references..."
echo "  $OLD_MODULE → $MODULE_PATH"

# Replace in go.mod
if [ -f go.mod ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|${OLD_MODULE}|${MODULE_PATH}|g" go.mod
    else
        sed -i "s|${OLD_MODULE}|${MODULE_PATH}|g" go.mod
    fi
fi

# Replace in all .go files
find . -name "*.go" -type f | while read -r file; do
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|${OLD_MODULE}|${MODULE_PATH}|g" "$file"
    else
        sed -i "s|${OLD_MODULE}|${MODULE_PATH}|g" "$file"
    fi
done

echo "Module path set to: $MODULE_PATH"

# Copy .env.example to .env
if [ -f .env.example ] && [ ! -f .env ]; then
    cp .env.example .env
    echo ".env created from .env.example — edit with your settings"
elif [ -f .env ]; then
    echo ".env already exists, skipping"
fi

# Run go mod tidy
echo ""
echo "Running go mod tidy..."
go mod tidy

echo ""
echo "===================="
echo "Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Edit .env with your database URL and JWT secret"
echo "  2. Run: make dev"
```

Make it executable:
```bash
chmod +x setup.sh
```

- [ ] **Step 5: Verify setup.sh works**

Run:
```bash
./setup.sh
```
Input: `github.com/test/demo`
Expected: replaces module references, copies .env, runs go mod tidy

Then restore the original module name:
```bash
# Replace back to original for the template
find . -name "*.go" -type f -exec sed -i '' 's|github.com/test/demo|github.com/wesleysnt/gobase|g' {} +
sed -i '' 's|github.com/test/demo|github.com/wesleysnt/gobase|g' go.mod
go mod tidy
```

- [ ] **Step 6: Build and verify binary**

Run: `make build`
Expected: binary created at `bin/gobase`

Run: `./bin/gobase --help`
Expected: cobra help output

- [ ] **Step 7: Commit**

```bash
git add migrations/ .env.example Makefile setup.sh
git commit -m "feat: add project scaffolding — migrations, env config, Makefile, and setup script

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 10: Test utilities & final verification

**Files:**
- Create: `internal/testutil/db.go`
- Create: `internal/testutil/auth.go`

**Interfaces:**
- Produces: `func SetupTestDB(t *testing.T) *sqlx.DB` — connects to `TEST_DATABASE_URL`, runs migrations
- Produces: `func AuthHeader(t *testing.T, userID string) string` — generates JWT and returns `Bearer <token>`

- [ ] **Step 1: Implement testutil/db.go**

```go
// internal/testutil/db.go
package testutil

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/wesleysnt/gobase/internal/config"
	"github.com/wesleysnt/gobase/internal/database"
)

func SetupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	cfg := &config.Config{
		DatabaseURL:       dsn,
		DBMaxOpenConns:    5,
		DBMaxIdleConns:    2,
		DBConnMaxLifetime: 1,
	}

	db, err := database.Connect(cfg)
	if err != nil {
		t.Fatalf("SetupTestDB: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db, dsn, "up", 0); err != nil {
		db.Close()
		t.Fatalf("SetupTestDB migrations: %v", err)
	}

	// Clean up after test
	t.Cleanup(func() {
		// Truncate all tables between tests
		db.Exec("DELETE FROM users")
		db.Close()
	})

	return db
}
```

- [ ] **Step 2: Implement testutil/auth.go**

```go
// internal/testutil/auth.go
package testutil

import (
	"os"
	"testing"
	"time"

	"github.com/wesleysnt/gobase/internal/auth"
)

func AuthHeader(t *testing.T, userID string) string {
	t.Helper()

	secret, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		secret = "test-secret-for-integration-tests"
	}

	token, _, err := auth.GenerateToken(userID, 1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("AuthHeader: %v", err)
	}

	return "Bearer " + token
}
```

- [ ] **Step 3: Verify everything compiles**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 4: Run all tests**

Run: `go test ./... -v`
Expected: all tests PASS

- [ ] **Step 5: Run go vet**

Run: `go vet ./...`
Expected: no issues

- [ ] **Step 6: Final build**

Run:
```bash
make build
./bin/gobase --help
```

Expected: clean build, help output

- [ ] **Step 7: Commit**

```bash
git add internal/testutil/
git commit -m "feat: add test utilities for database setup and auth headers

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 11: .gitignore

**Files:**
- Create: `.gitignore`

- [ ] **Step 1: Create .gitignore**

```
# Binary
/gobase
/bin/

# Environment
.env

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Test
coverage.out
```

- [ ] **Step 2: Commit**

```bash
git add .gitignore
git commit -m "chore: add .gitignore

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### Task 12: Final integration verification

- [ ] **Step 1: Clean test run**

```bash
go clean -testcache
go test ./... -v -count=1
```

Expected: all tests PASS

- [ ] **Step 2: Verify build is clean**

```bash
go build -o /dev/null .
go vet ./...
```

Expected: no errors, no warnings

- [ ] **Step 3: Verify migration create CLI**

```bash
go run . migrate create test_migration
ls migrations/*test_migration*
```

Expected: two new migration files created, then clean up:
```bash
rm migrations/*test_migration*
```

- [ ] **Step 4: Verify setup.sh re-run works**

```bash
# Test with a dry run to verify it doesn't break
grep -r "github.com/wesleysnt/gobase" --include="*.go" | head -5
```

Expected: shows files with the placeholder module path

- [ ] **Step 5: Verify all files are committed**

```bash
git status
```

Expected: clean working tree, nothing unstaged
