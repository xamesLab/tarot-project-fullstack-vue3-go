package utils

import (
	"fmt"
	"time"

	"tarot-daily/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID, role, timezone, secret, ttl string) (string, error) {
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		return "", fmt.Errorf("parse JWT_TTL: %w", err)
	}

	now := time.Now().UTC()
	claims := models.JWTClaims{
		UserID:   userID,
		Role:     role,
		Timezone: timezone,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString, secret string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
