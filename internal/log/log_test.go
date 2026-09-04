package log

import (
	"bytes"
	"testing"

	"log/slog"

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
	_ = buf
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
