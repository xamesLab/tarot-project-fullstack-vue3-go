package jwt

import (
	"errors"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"tarotdaily/backend/internal/domain/models"
)

type Claims struct {
	UserID string `json:"uid"`
	Login  string `json:"login"`
	Role   string `json:"role"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *Manager) Generate(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	now := time.Now().UTC()
	claims := Claims{
		UserID: user.ID.Hex(),
		Login:  user.Login,
		Role:   string(user.Role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   user.ID.Hex(),
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	token, err := jwtlib.ParseWithClaims(tokenString, &Claims{}, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
