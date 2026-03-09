package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/inpututil"
	"tarotdaily/backend/internal/platform/timeutil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

const (
	minLoginLength    = 3
	maxLoginLength    = 64
	minPasswordLength = 8
)

func (s *Service) Register(ctx context.Context, login, password string) (*AuthResult, error) {
	normalizedLogin, normalizedPassword, err := validateCredentials(login, password)
	if err != nil {
		return nil, ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(normalizedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &models.User{
		ID:           primitive.NewObjectID(),
		Login:        normalizedLogin,
		PasswordHash: string(passwordHash),
		Role:         models.UserRoleUser,
		Status:       models.UserStatusActive,
		Timezone:     timeutil.DefaultTimezone,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		if mongorepo.IsDuplicateKeyError(err) {
			return nil, ErrLoginAlreadyExists
		}
		return nil, err
	}

	return s.issueAuthResult(user)
}

func (s *Service) Login(ctx context.Context, login, password string) (*AuthResult, error) {
	normalizedLogin, normalizedPassword, err := validateCredentials(login, password)
	if err != nil {
		return nil, ErrInvalidInput
	}

	user, err := s.users.FindByLogin(ctx, normalizedLogin)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(normalizedPassword)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.issueAuthResult(user)
}

func (s *Service) Me(ctx context.Context, userID string) (*UserView, error) {
	userOID, err := inpututil.ParseObjectID(userID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	user, err := s.users.FindByID(ctx, userOID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	view := toUserView(user)
	return &view, nil
}

func (s *Service) issueAuthResult(user *models.User) (*AuthResult, error) {
	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  toUserView(user),
	}, nil
}

func validateCredentials(login, password string) (string, string, error) {
	normalizedLogin, err := inpututil.RequireTrimmed(login)
	if err != nil {
		return "", "", ErrInvalidInput
	}

	normalizedPassword, err := inpututil.RequireTrimmed(password)
	if err != nil {
		return "", "", ErrInvalidInput
	}

	if len(normalizedLogin) < minLoginLength || len(normalizedLogin) > maxLoginLength {
		return "", "", ErrInvalidInput
	}
	if len(strings.TrimSpace(normalizedPassword)) < minPasswordLength {
		return "", "", ErrInvalidInput
	}

	return normalizedLogin, normalizedPassword, nil
}
