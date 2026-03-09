package decks

import (
	"context"
	"errors"
	"path/filepath"

	"tarotdaily/backend/internal/platform/inpututil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

func (s *Service) UploadBack(ctx context.Context, deckID, contentType string, payload []byte) (*DeckView, error) {
	deckOID, err := inpututil.ParseObjectID(deckID)
	if err != nil {
		return nil, ErrInvalidInput
	}
	if err := validateImagePayload(contentType, payload); err != nil {
		return nil, err
	}

	deck, err := s.decks.FindByID(ctx, deckOID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}

	ext := normalizeExtFromContentType(contentType)
	relativePath := filepath.Join("decks", deck.ID.Hex(), "back"+ext)
	assetURL, err := s.storage.Save(relativePath, payload)
	if err != nil {
		return nil, err
	}

	if deck.Cards == nil {
		deck.Cards = map[string]string{}
	}
	status := calculateStatus(assetURL, len(deck.Cards), deck.RequiredCardsCount)

	updated, err := s.decks.UpdateAssets(ctx, deck.ID, assetURL, deck.Cards, status, len(deck.Cards))
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}

	return s.toSingleDeckView(ctx, updated)
}

func (s *Service) UploadCard(ctx context.Context, deckID, cardKey, contentType string, payload []byte) (*DeckView, error) {
	deckOID, err := inpututil.ParseObjectID(deckID)
	if err != nil {
		return nil, ErrInvalidInput
	}
	if err := validateImagePayload(contentType, payload); err != nil {
		return nil, err
	}

	normalizedCardKey, err := normalizeCardKey(cardKey)
	if err != nil {
		return nil, ErrInvalidInput
	}

	deck, err := s.decks.FindByID(ctx, deckOID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}

	if deck.Cards == nil {
		deck.Cards = map[string]string{}
	}
	if _, exists := deck.Cards[normalizedCardKey]; exists {
		return nil, ErrCardAlreadyExists
	}
	if len(deck.Cards) >= deck.RequiredCardsCount {
		return nil, ErrDeckAlreadyFull
	}

	ext := normalizeExtFromContentType(contentType)
	relativePath := filepath.Join("decks", deck.ID.Hex(), "cards", normalizedCardKey+ext)
	assetURL, err := s.storage.Save(relativePath, payload)
	if err != nil {
		return nil, err
	}

	deck.Cards[normalizedCardKey] = assetURL
	cardsCount := len(deck.Cards)
	status := calculateStatus(deck.BackImageURL, cardsCount, deck.RequiredCardsCount)

	updated, err := s.decks.UpdateAssets(ctx, deck.ID, deck.BackImageURL, deck.Cards, status, cardsCount)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrDeckNotFound
		}
		return nil, err
	}

	return s.toSingleDeckView(ctx, updated)
}
