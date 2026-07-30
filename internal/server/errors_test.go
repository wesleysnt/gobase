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
