package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	users *repositories.UserRepository
}

func NewAuthService(users *repositories.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest, jwtSecret, jwtTTL string) (*models.User, string, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	username := strings.TrimSpace(req.Username)

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return nil, "", wrapValidation("Email is already registered")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", err
	}

	if _, err := s.users.FindByUsername(ctx, username); err == nil {
		return nil, "", wrapValidation("Username unavailable")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", err
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	user := &models.User{
		Email:          email,
		Username:       username,
		PasswordHash:   passwordHash,
		Role:           "seeker",
		Status:         "active",
		Timezone:       req.Timezone,
		Language:       "en",
		Theme:          "light",
		SelectedDeckID: "default",
		SelectedBackID: "default",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		if isDuplicateKey(err) {
			message := "Email is already registered"
			if strings.Contains(strings.ToLower(err.Error()), "username") {
				message = "Username unavailable"
			}
			return nil, "", wrapValidation(message)
		}
		return nil, "", err
	}

	token, err := utils.GenerateToken(user.ID.Hex(), user.Role, user.Timezone, jwtSecret, jwtTTL)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest, jwtSecret, jwtTTL string) (*models.User, string, error) {
	user, err := s.users.FindByEmail(ctx, strings.TrimSpace(strings.ToLower(req.Email)))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	if user.Status == "disabled" {
		return nil, "", ErrAccountDisabled
	}

	if !utils.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, "", ErrInvalidCredentials
	}

	now := time.Now().UTC()
	if err := s.users.UpdateLastLogin(ctx, user.ID, now); err != nil {
		return nil, "", err
	}
	user.LastLoginAt = &now
	user.UpdatedAt = now

	token, err := utils.GenerateToken(user.ID.Hex(), user.Role, user.Timezone, jwtSecret, jwtTTL)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func isDuplicateKey(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}

func wrapValidation(message string) error {
	return errors.Join(ErrValidationFailed, errors.New(message))
}
