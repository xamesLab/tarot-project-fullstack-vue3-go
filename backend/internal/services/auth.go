package services

import (
	"context"
	"errors"
	"regexp"
	"strconv"
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
	username, err := s.resolveRegistrationUsername(ctx, strings.TrimSpace(req.Username), email)
	if err != nil {
		return nil, "", err
	}

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return nil, "", wrapValidation("Email is already registered")
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
		SelectedDeckID: "",
		SelectedBackID: "",
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

func (s *AuthService) resolveRegistrationUsername(ctx context.Context, rawUsername, email string) (string, error) {
	username := strings.TrimSpace(rawUsername)
	if username == "" {
		return s.ensureUniqueUsername(ctx, buildUsernameFromEmail(email))
	}

	if _, err := s.users.FindByUsername(ctx, username); err == nil {
		return "", wrapValidation("Username unavailable")
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return "", err
	}

	return username, nil
}

func (s *AuthService) ensureUniqueUsername(ctx context.Context, base string) (string, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		base = "user"
	}

	for attempt := 0; attempt < 1000; attempt++ {
		candidate := base
		if attempt > 0 {
			suffix := strconv.Itoa(attempt + 1)
			candidate = truncateUsername(base, 20-len(suffix)) + suffix
		}

		_, err := s.users.FindByUsername(ctx, candidate)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return candidate, nil
		}
		if err != nil {
			return "", err
		}
	}

	return "", wrapValidation("Username unavailable")
}

func buildUsernameFromEmail(email string) string {
	localPart := strings.TrimSpace(strings.SplitN(strings.ToLower(email), "@", 2)[0])
	localPart = regexp.MustCompile(`[^a-z0-9._-]+`).ReplaceAllString(localPart, "")

	switch {
	case len(localPart) >= 3:
		return truncateUsername(localPart, 20)
	case len(localPart) == 2:
		return localPart + "1"
	case len(localPart) == 1:
		return localPart + "01"
	default:
		return "user"
	}
}

func truncateUsername(value string, limit int) string {
	if limit <= 0 {
		return ""
	}
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
