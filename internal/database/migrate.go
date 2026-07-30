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
