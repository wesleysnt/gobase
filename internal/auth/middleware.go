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
