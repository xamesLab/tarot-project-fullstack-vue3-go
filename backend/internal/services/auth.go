package services

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"tarot-daily/internal/models"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	users    authUserRepository
	settings authSettingsRepository
}

type authUserRepository interface {
	Create(context.Context, *models.User) error
	FindByEmail(context.Context, string) (*models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	UpdateLastLogin(context.Context, primitive.ObjectID, time.Time) error
}

type authSettingsRepository interface {
	Get(context.Context) (*models.AdminSettings, error)
}

var invalidUsernameCharacters = regexp.MustCompile(`[^a-z0-9._-]+`)

func NewAuthService(users authUserRepository, settings authSettingsRepository) *AuthService {
	return &AuthService{users: users, settings: settings}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest, jwtSecret, jwtTTL string) (*models.User, string, error) {
	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, "", err
	}
	if !settings.RegistrationEnabled {
		return nil, "", ErrRegistrationDisabled
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	username, err := s.resolveRegistrationUsername(ctx, strings.TrimSpace(req.Username), email)
	if err != nil {
		return nil, "", err
	}

	if _, err := s.users.FindByEmail(ctx, email); err == nil {
		return nil, "", newValidationError("email", "already_registered", "Email is already registered")
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
		Role:           publicDefaultRole(settings.DefaultRole),
		Status:         "active",
		Timezone:       req.Timezone,
		Language:       "ru",
		Theme:          "dark",
		SelectedDeckID: "",
		SelectedBackID: "",
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.users.Create(ctx, user); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			message := "Email is already registered"
			if strings.Contains(strings.ToLower(err.Error()), "username") {
				message = "Username unavailable"
			}
			field := "email"
			reason := "already_registered"
			if message == "Username unavailable" {
				field = "username"
				reason = "unavailable"
			}
			return nil, "", newValidationError(field, reason, message)
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

func (s *AuthService) resolveRegistrationUsername(ctx context.Context, rawUsername, email string) (string, error) {
	username := strings.TrimSpace(rawUsername)
	if username == "" {
		return s.ensureUniqueUsername(ctx, buildUsernameFromEmail(email))
	}

	if _, err := s.users.FindByUsername(ctx, username); err == nil {
		return "", newValidationError("username", "unavailable", "Username unavailable")
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

	return "", newValidationError("username", "unavailable", "Username unavailable")
}

func publicDefaultRole(role string) string {
	switch strings.TrimSpace(role) {
	case "seeker", "adept", "master":
		return strings.TrimSpace(role)
	default:
		return "seeker"
	}
}

func buildUsernameFromEmail(email string) string {
	localPart := strings.TrimSpace(strings.SplitN(strings.ToLower(email), "@", 2)[0])
	localPart = invalidUsernameCharacters.ReplaceAllString(localPart, "")

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
