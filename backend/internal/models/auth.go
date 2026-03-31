package models

import "github.com/golang-jwt/jwt/v5"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
	Password string `json:"password" validate:"required,min=8,password"`
	Timezone string `json:"timezone" validate:"required,timezone"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Role     string `json:"role"`
	Timezone string `json:"timezone"`
	jwt.RegisteredClaims
}
