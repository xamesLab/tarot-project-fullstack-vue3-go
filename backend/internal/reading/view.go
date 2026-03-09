package reading

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
)

const defaultBackImageURL = "/assets/cards/back.jpg"

func (s *Service) toTodayView(
	ctx context.Context,
	daily *models.DailyReading,
	limit int,
	deckID *primitive.ObjectID,
) (*TodayView, error) {
	cards, err := s.cards.FindByIDs(ctx, collectCardIDs(daily))
	if err != nil {
		return nil, err
	}

	cardByID := indexCardsByID(cards)
	assets, backImageURL := s.resolveReadingAssets(ctx, deckID)

	cardOfDay, err := toCardView(daily.CardOfDay, cardByID, assets, backImageURL)
	if err != nil {
		return nil, err
	}

	spreadA, err := toSpreadView(daily.SpreadA, cardByID, assets, backImageURL)
	if err != nil {
		return nil, err
	}

	spreadB, err := toSpreadView(daily.SpreadB, cardByID, assets, backImageURL)
	if err != nil {
		return nil, err
	}

	view := &TodayView{
		Date:                  daily.Date,
		TimezoneUsed:          daily.TimezoneUsed,
		CardOfDay:             cardOfDay,
		SpreadA:               spreadA,
		SpreadB:               spreadB,
		SpreadGenerationsUsed: daily.SpreadGenerationsCount,
		CanGenerateMore:       true,
	}

	applyDailyLimit(view, limit, daily.SpreadGenerationsCount)

	return view, nil
}

func (s *Service) resolveReadingAssets(ctx context.Context, refDeckID *primitive.ObjectID) (map[string]string, string) {
	if refDeckID == nil {
		return map[string]string{}, defaultBackImageURL
	}

	deck, err := s.decks.FindByID(ctx, *refDeckID)
	if err == nil {
		backImageURL := strings.TrimSpace(deck.BackImageURL)
		if backImageURL == "" {
			backImageURL = defaultBackImageURL
		}
		if deck.Status == models.DeckStatusDone && deck.Cards != nil {
			return deck.Cards, backImageURL
		}
	}

	theme, err := s.deckThemes.FindByID(ctx, *refDeckID)
	if err == nil && theme.Assets != nil {
		return theme.Assets, defaultBackImageURL
	}

	return map[string]string{}, defaultBackImageURL
}

func collectCardIDs(daily *models.DailyReading) []primitive.ObjectID {
	cardIDs := make([]primitive.ObjectID, 0, 1+len(daily.SpreadA)+len(daily.SpreadB))
	cardIDs = append(cardIDs, daily.CardOfDay.CardID)

	for _, draw := range daily.SpreadA {
		cardIDs = append(cardIDs, draw.CardID)
	}

	for _, draw := range daily.SpreadB {
		cardIDs = append(cardIDs, draw.CardID)
	}

	return cardIDs
}

func indexCardsByID(cards []models.Card) map[string]models.Card {
	out := make(map[string]models.Card, len(cards))
	for _, card := range cards {
		out[card.ID.Hex()] = card
	}
	return out
}

func applyDailyLimit(view *TodayView, limit int, used int) {
	if limit == models.UnlimitedDailySpreadLimit {
		return
	}

	view.SpreadGenerationsLimit = &limit
	remaining := limit - used
	if remaining < 0 {
		remaining = 0
	}
	view.SpreadGenerationsRemaining = &remaining
	view.CanGenerateMore = remaining > 0
}

func toSpreadView(
	draws []models.CardDraw,
	cardByID map[string]models.Card,
	assets map[string]string,
	backImageURL string,
) ([]CardView, error) {
	out := make([]CardView, 0, len(draws))
	for _, draw := range draws {
		cardView, err := toCardView(draw, cardByID, assets, backImageURL)
		if err != nil {
			return nil, err
		}
		out = append(out, cardView)
	}
	return out, nil
}

func toCardView(
	draw models.CardDraw,
	cardByID map[string]models.Card,
	assets map[string]string,
	backImageURL string,
) (CardView, error) {
	card, ok := cardByID[draw.CardID.Hex()]
	if !ok {
		return CardView{}, fmt.Errorf("card not found for id: %s", draw.CardID.Hex())
	}

	interpretation := card.InterpretationOpen
	if draw.State == models.CardStateClosed {
		interpretation = card.InterpretationClosed
	}

	return CardView{
		ID:             card.ID.Hex(),
		Name:           card.Name,
		State:          draw.State,
		ImageURL:       resolveCardImageURL(card.Code, assets),
		BackImageURL:   backImageURL,
		Description:    card.Description,
		Interpretation: interpretation,
	}, nil
}

func resolveCardImageURL(cardCode string, assets map[string]string) string {
	if imageURL := assets[cardCode]; imageURL != "" {
		return imageURL
	}

	if _, err := strconv.Atoi(strings.TrimSpace(cardCode)); err == nil {
		return fmt.Sprintf("/assets/cards/%s.jpg", cardCode)
	}

	return fmt.Sprintf("/themes/default/%s.jpg", cardCode)
}
