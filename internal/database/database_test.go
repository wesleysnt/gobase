package database

import (
	"os"
	"testing"

	"github.com/you/gobase/internal/config"
)

func TestParseDriver(t *testing.T) {
	tests := []struct {
		url     string
		driver  string
		wantErr bool
	}{
		{"postgres://localhost:5432/db", "postgres", false},
		{"postgresql://localhost:5432/db", "postgres", false},
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
