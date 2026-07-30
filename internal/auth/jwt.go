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
