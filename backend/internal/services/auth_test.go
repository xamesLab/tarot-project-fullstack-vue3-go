package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"tarot-daily/internal/models"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type fakeAuthUsers struct {
	byEmail    map[string]*models.User
	byUsername map[string]*models.User
	created    *models.User
}

func newFakeAuthUsers() *fakeAuthUsers {
	return &fakeAuthUsers{byEmail: map[string]*models.User{}, byUsername: map[string]*models.User{}}
}

func (f *fakeAuthUsers) Create(_ context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	f.created = user
	f.byEmail[user.Email] = user
	f.byUsername[user.Username] = user
	return nil
}

func (f *fakeAuthUsers) FindByEmail(_ context.Context, email string) (*models.User, error) {
	user, ok := f.byEmail[email]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return user, nil
}

func (f *fakeAuthUsers) FindByUsername(_ context.Context, username string) (*models.User, error) {
	user, ok := f.byUsername[username]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return user, nil
}

func (f *fakeAuthUsers) UpdateLastLogin(_ context.Context, _ primitive.ObjectID, _ time.Time) error {
	return nil
}

type fakeAuthSettings struct{ settings models.AdminSettings }

func (f *fakeAuthSettings) Get(context.Context) (*models.AdminSettings, error) {
	return &f.settings, nil
}

func TestBuildUsernameFromEmail(t *testing.T) {
	if got := buildUsernameFromEmail("alice@example.com"); got != "alice" {
		t.Fatalf("expected alice, got %s", got)
	}

	if got := buildUsernameFromEmail("alice.smith@example.com"); got != "alice.smith" {
		t.Fatalf("expected alice.smith, got %s", got)
	}

	if got := buildUsernameFromEmail("a@example.com"); got != "a01" {
		t.Fatalf("expected a01, got %s", got)
	}
}

func TestRegisterRespectsRegistrationEnabled(t *testing.T) {
	users := newFakeAuthUsers()
	service := NewAuthService(users, &fakeAuthSettings{settings: models.AdminSettings{RegistrationEnabled: false, DefaultRole: "seeker"}})
	_, _, err := service.Register(context.Background(), models.RegisterRequest{Email: "alice@example.com", Password: "Secret123", Timezone: "UTC"}, "secret", "1h")
	if !errors.Is(err, ErrRegistrationDisabled) {
		t.Fatalf("expected registration disabled, got %v", err)
	}
	if users.created != nil {
		t.Fatal("user must not be created while registration is disabled")
	}
}

func TestRegisterAppliesSafeDefaultRole(t *testing.T) {
	for _, test := range []struct{ configured, expected string }{{"adept", "adept"}, {"master", "master"}, {"keeper", "seeker"}, {"", "seeker"}} {
		t.Run(test.configured, func(t *testing.T) {
			users := newFakeAuthUsers()
			service := NewAuthService(users, &fakeAuthSettings{settings: models.AdminSettings{RegistrationEnabled: true, DefaultRole: test.configured}})
			user, token, err := service.Register(context.Background(), models.RegisterRequest{Email: "alice@example.com", Password: "Secret123", Timezone: "UTC"}, "secret", "1h")
			if err != nil {
				t.Fatalf("register: %v", err)
			}
			if user.Role != test.expected {
				t.Fatalf("expected role %s, got %s", test.expected, user.Role)
			}
			if token == "" {
				t.Fatal("expected JWT token")
			}
		})
	}
}

func TestRegisterReturnsStructuredDuplicateErrors(t *testing.T) {
	users := newFakeAuthUsers()
	users.byEmail["alice@example.com"] = &models.User{ID: primitive.NewObjectID(), Email: "alice@example.com", Username: "alice"}
	service := NewAuthService(users, &fakeAuthSettings{settings: models.AdminSettings{RegistrationEnabled: true, DefaultRole: "seeker"}})
	_, _, err := service.Register(context.Background(), models.RegisterRequest{Email: "alice@example.com", Username: "other", Password: "Secret123", Timezone: "UTC"}, "secret", "1h")
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if validationErr.Field != "email" || validationErr.Reason != "already_registered" {
		t.Fatalf("unexpected details: %+v", validationErr)
	}

	users = newFakeAuthUsers()
	users.byUsername["alice"] = &models.User{ID: primitive.NewObjectID(), Email: "other@example.com", Username: "alice"}
	service = NewAuthService(users, &fakeAuthSettings{settings: models.AdminSettings{RegistrationEnabled: true, DefaultRole: "seeker"}})
	_, _, err = service.Register(context.Background(), models.RegisterRequest{Email: "new@example.com", Username: "alice", Password: "Secret123", Timezone: "UTC"}, "secret", "1h")
	validationErr = nil
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if validationErr.Field != "username" || validationErr.Reason != "unavailable" {
		t.Fatalf("unexpected details: %+v", validationErr)
	}
}

func TestLoginReturnsUserAndToken(t *testing.T) {
	hash, err := utils.HashPassword("Secret123")
	if err != nil {
		t.Fatal(err)
	}
	users := newFakeAuthUsers()
	users.byEmail["alice@example.com"] = &models.User{ID: primitive.NewObjectID(), Email: "alice@example.com", Username: "alice", PasswordHash: hash, Role: "seeker", Status: "active", Timezone: "UTC"}
	service := NewAuthService(users, &fakeAuthSettings{})
	user, token, err := service.Login(context.Background(), models.LoginRequest{Email: "Alice@example.com", Password: "Secret123"}, "secret", "1h")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if user.Username != "alice" || token == "" {
		t.Fatalf("unexpected login result: %#v %q", user, token)
	}
}
