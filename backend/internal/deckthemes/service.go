package deckthemes

import (
	"context"
	"errors"
	"strings"

	"tarotdaily/backend/internal/platform/inpututil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrThemeNameConflict = errors.New("theme name already exists")
)

type ThemeView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	repo *mongorepo.DeckThemeRepository
}

func NewService(repo *mongorepo.DeckThemeRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]ThemeView, error) {
	items, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]ThemeView, 0, len(items))
	for _, item := range items {
		out = append(out, ThemeView{
			ID:   item.ID.Hex(),
			Name: item.Name,
		})
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, adminUserID, name string) (*ThemeView, error) {
	normalizedName := strings.TrimSpace(name)
	if len(normalizedName) < 2 || len(normalizedName) > 128 {
		return nil, ErrInvalidInput
	}

	adminOID, err := inpututil.ParseObjectID(adminUserID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	item, err := s.repo.Create(ctx, normalizedName, adminOID)
	if err != nil {
		if mongorepo.IsDuplicateKeyError(err) {
			return nil, ErrThemeNameConflict
		}
		return nil, err
	}

	out := &ThemeView{
		ID:   item.ID.Hex(),
		Name: item.Name,
	}
	return out, nil
}
