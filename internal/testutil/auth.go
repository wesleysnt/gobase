package testutil

import (
	"os"
	"testing"
	"time"

	"github.com/you/gobase/internal/auth"
)

func AuthHeader(t *testing.T, userID string) string {
	t.Helper()

	var secret []byte
	if s, ok := os.LookupEnv("JWT_SECRET"); ok {
		secret = []byte(s)
	} else {
		secret = []byte("test-secret-for-integration-tests")
	}

	token, _, err := auth.GenerateToken(userID, 1*time.Hour, nil, secret)
	if err != nil {
		t.Fatalf("AuthHeader: %v", err)
	}

	return "Bearer " + token
}
