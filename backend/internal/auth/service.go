package auth

import (
	"errors"

	jwtplatform "tarotdaily/backend/internal/platform/jwt"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

var (
	ErrInvalidInput           = errors.New("invalid input")
	ErrLoginAlreadyExists     = errors.New("login already exists")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrUserNotFound           = errors.New("user not found")
	ErrDeckThemeNotFound      = errors.New("deck theme not found")
	ErrPreferredDeckNotUsable = errors.New("preferred deck is not usable")
)

type UserView struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	Role            string `json:"role"`
	Status          string `json:"status"`
	DisplayName     string `json:"displayName,omitempty"`
	Timezone        string `json:"timezone"`
	PreferredDeckID string `json:"preferredDeckId,omitempty"`
	DeckThemeID     string `json:"deckThemeId,omitempty"`
}

type AuthResult struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

type UpdateSettingsInput struct {
	DisplayName     *string
	Timezone        *string
	PreferredDeckID *string
	DeckThemeID     *string
}

type Service struct {
	users      *mongorepo.UserRepository
	decks      *mongorepo.DeckRepository
	deckThemes *mongorepo.DeckThemeRepository
	jwtManager *jwtplatform.Manager
}

func NewService(
	users *mongorepo.UserRepository,
	decks *mongorepo.DeckRepository,
	deckThemes *mongorepo.DeckThemeRepository,
	jwtManager *jwtplatform.Manager,
) *Service {
	return &Service{
		users:      users,
		decks:      decks,
		deckThemes: deckThemes,
		jwtManager: jwtManager,
	}
}
