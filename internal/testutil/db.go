package testutil

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/you/gobase/internal/config"
	"github.com/you/gobase/internal/database"
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

	if err := database.RunMigrations(db, dsn, "up", 0); err != nil {
		db.Close()
		t.Fatalf("SetupTestDB migrations: %v", err)
	}

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Close()
	})

	return db
}
