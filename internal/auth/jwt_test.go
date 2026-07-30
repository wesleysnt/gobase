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
