package reading

import (
	"context"
	"errors"

	"tarotdaily/backend/internal/domain/models"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
	"tarotdaily/backend/internal/roles"
)

const dayResetHour = 3

var (
	ErrInvalidInput             = errors.New("invalid input")
	ErrUserNotFound             = errors.New("user not found")
	ErrRolePolicyNotFound       = errors.New("role policy not found")
	ErrDailySpreadLimitReached  = errors.New("daily spread limit reached")
	ErrCardsCatalogNotAvailable = errors.New("cards catalog is not available")
)

type CardView struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	State          string `json:"state"`
	ImageURL       string `json:"imageUrl"`
	BackImageURL   string `json:"backImageUrl"`
	Description    string `json:"description"`
	Interpretation string `json:"interpretation"`
}

type TodayView struct {
	Date                       string     `json:"date"`
	TimezoneUsed               string     `json:"timezoneUsed"`
	CardOfDay                  CardView   `json:"cardOfDay"`
	SpreadA                    []CardView `json:"spreadA"`
	SpreadB                    []CardView `json:"spreadB"`
	SpreadGenerationsUsed      int        `json:"spreadGenerationsUsed"`
	SpreadGenerationsLimit     *int       `json:"spreadGenerationsLimit,omitempty"`
	SpreadGenerationsRemaining *int       `json:"spreadGenerationsRemaining,omitempty"`
	CanGenerateMore            bool       `json:"canGenerateMore"`
}

type Service struct {
	users      *mongorepo.UserRepository
	cards      *mongorepo.CardRepository
	decks      *mongorepo.DeckRepository
	readings   *mongorepo.DailyReadingRepository
	deckThemes *mongorepo.DeckThemeRepository
	roles      *roles.Service
}

func NewService(
	users *mongorepo.UserRepository,
	cards *mongorepo.CardRepository,
	decks *mongorepo.DeckRepository,
	readings *mongorepo.DailyReadingRepository,
	deckThemes *mongorepo.DeckThemeRepository,
	rolesService *roles.Service,
) *Service {
	return &Service{
		users:      users,
		cards:      cards,
		decks:      decks,
		readings:   readings,
		deckThemes: deckThemes,
		roles:      rolesService,
	}
}

func (s *Service) GetToday(ctx context.Context, userID string) (*TodayView, error) {
	user, _, _, daily, limit, err := s.ensureDaily(ctx, userID)
	if err != nil {
		return nil, err
	}

	// First open of the day gets the first spreads automatically.
	if len(daily.SpreadA) == 0 || len(daily.SpreadB) == 0 {
		if limit != models.UnlimitedDailySpreadLimit && daily.SpreadGenerationsCount >= limit {
			return s.toTodayView(ctx, daily, limit, resolveReadingDeckID(user))
		}
		updated, err := s.generateAndStoreSpreads(ctx, daily, resolveReadingDeckID(user))
		if err != nil {
			return nil, err
		}
		daily = updated
	}

	return s.toTodayView(ctx, daily, limit, resolveReadingDeckID(user))
}

func (s *Service) GenerateSpreads(ctx context.Context, userID string) (*TodayView, error) {
	user, _, _, daily, limit, err := s.ensureDaily(ctx, userID)
	if err != nil {
		return nil, err
	}

	if limit != models.UnlimitedDailySpreadLimit && daily.SpreadGenerationsCount >= limit {
		return nil, ErrDailySpreadLimitReached
	}

	updated, err := s.generateAndStoreSpreads(ctx, daily, resolveReadingDeckID(user))
	if err != nil {
		return nil, err
	}

	return s.toTodayView(ctx, updated, limit, resolveReadingDeckID(user))
}
