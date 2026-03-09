package decks

import (
	"context"
	"errors"
	"strings"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/inpututil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

const MaxUploadBytes = 8 * 1024 * 1024 // 8MB

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrDeckNameConflict  = errors.New("deck name already exists")
	ErrDeckNotFound      = errors.New("deck not found")
	ErrUploadTooLarge    = errors.New("upload too large")
	ErrUnsupportedImage  = errors.New("unsupported image")
	ErrCardAlreadyExists = errors.New("card already uploaded")
	ErrDeckAlreadyFull   = errors.New("deck already full")
)

var allowedContentTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

var validCardKeys = buildValidCardKeySet()

type DeckView struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	Status             string            `json:"status"`
	CreatedByUserID    string            `json:"createdByUserId"`
	CreatedByName      string            `json:"createdByName"`
	BackImageURL       string            `json:"backImageUrl,omitempty"`
	Cards              map[string]string `json:"cards,omitempty"`
	CardsUploadedCount int               `json:"cardsUploadedCount"`
	RequiredCardsCount int               `json:"requiredCardsCount"`
	CreatedAt          string            `json:"createdAt"`
	UpdatedAt          string            `json:"updatedAt"`
}

type AvailableDeckView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	decks   *mongorepo.DeckRepository
	users   *mongorepo.UserRepository
	storage *UploadStorage
}

func NewService(decks *mongorepo.DeckRepository, users *mongorepo.UserRepository) *Service {
	return &Service{
		decks:   decks,
		users:   users,
		storage: defaultUploadStorage,
	}
}

func (s *Service) Create(ctx context.Context, adminUserID, name string) (*DeckView, error) {
	normalizedName := strings.TrimSpace(name)
	if len(normalizedName) < 2 || len(normalizedName) > 128 {
		return nil, ErrInvalidInput
	}

	adminOID, err := inpututil.ParseObjectID(adminUserID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	item, err := s.decks.Create(ctx, normalizedName, adminOID)
	if err != nil {
		if mongorepo.IsDuplicateKeyError(err) {
			return nil, ErrDeckNameConflict
		}
		return nil, err
	}

	userMap, err := s.resolveUserNames(ctx, []models.Deck{*item})
	if err != nil {
		return nil, err
	}

	view := toDeckView(*item, userMap[item.CreatedByUserID.Hex()])
	return &view, nil
}

func (s *Service) List(ctx context.Context) ([]DeckView, error) {
	items, err := s.decks.List(ctx)
	if err != nil {
		return nil, err
	}

	userMap, err := s.resolveUserNames(ctx, items)
	if err != nil {
		return nil, err
	}

	out := make([]DeckView, 0, len(items))
	for _, item := range items {
		out = append(out, toDeckView(item, userMap[item.CreatedByUserID.Hex()]))
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, deckID string) (*DeckView, error) {
	deckOID, err := inpututil.ParseObjectID(deckID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	item, err := s.decks.FindByID(ctx, deckOID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}

	userMap, err := s.resolveUserNames(ctx, []models.Deck{*item})
	if err != nil {
		return nil, err
	}

	view := toDeckView(*item, userMap[item.CreatedByUserID.Hex()])
	return &view, nil
}

func (s *Service) ListAvailable(ctx context.Context) ([]AvailableDeckView, error) {
	items, err := s.decks.ListByStatus(ctx, models.DeckStatusDone)
	if err != nil {
		return nil, err
	}

	out := make([]AvailableDeckView, 0, len(items))
	for _, item := range items {
		out = append(out, AvailableDeckView{
			ID:   item.ID.Hex(),
			Name: item.Name,
		})
	}
	return out, nil
}
