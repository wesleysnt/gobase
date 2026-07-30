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
