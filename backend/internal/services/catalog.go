package services

import (
	"context"

	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"
)

type CatalogService struct {
	decks     *repositories.DeckRepository
	cardBacks *repositories.CardBackRepository
	settings  *repositories.AdminSettingsRepository
}

func NewCatalogService(
	decks *repositories.DeckRepository,
	cardBacks *repositories.CardBackRepository,
	settings *repositories.AdminSettingsRepository,
) *CatalogService {
	return &CatalogService{
		decks:     decks,
		cardBacks: cardBacks,
		settings:  settings,
	}
}

func (s *CatalogService) ListAvailableDecks(ctx context.Context) ([]models.Deck, string, error) {
	decks, err := s.decks.ListAvailable(ctx)
	if err != nil {
		return nil, "", err
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, "", err
	}

	return decks, settings.DefaultDeckID, nil
}

func (s *CatalogService) ListAvailableCardBacks(ctx context.Context) ([]models.CardBack, string, error) {
	cardBacks, err := s.cardBacks.ListAvailable(ctx)
	if err != nil {
		return nil, "", err
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, "", err
	}

	return cardBacks, settings.DefaultBackID, nil
}
