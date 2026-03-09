package decks

import (
	"context"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
)

func (s *Service) toSingleDeckView(ctx context.Context, deck *models.Deck) (*DeckView, error) {
	userMap, err := s.resolveUserNames(ctx, []models.Deck{*deck})
	if err != nil {
		return nil, err
	}

	view := toDeckView(*deck, userMap[deck.CreatedByUserID.Hex()])
	return &view, nil
}

func (s *Service) resolveUserNames(ctx context.Context, decks []models.Deck) (map[string]string, error) {
	userIDs := make([]primitive.ObjectID, 0, len(decks))
	seen := map[string]struct{}{}
	for _, item := range decks {
		if item.CreatedByUserID.IsZero() {
			continue
		}
		key := item.CreatedByUserID.Hex()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		userIDs = append(userIDs, item.CreatedByUserID)
	}

	users, err := s.users.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(users))
	for _, user := range users {
		out[user.ID.Hex()] = models.DisplayNameOrLogin(user.DisplayName, user.Login)
	}
	return out, nil
}

func toDeckView(item models.Deck, createdByName string) DeckView {
	name := strings.TrimSpace(createdByName)
	if name == "" {
		name = "unknown"
	}

	cards := map[string]string{}
	for key, value := range item.Cards {
		cards[key] = value
	}

	return DeckView{
		ID:                 item.ID.Hex(),
		Name:               item.Name,
		Status:             string(item.Status),
		CreatedByUserID:    item.CreatedByUserID.Hex(),
		CreatedByName:      name,
		BackImageURL:       item.BackImageURL,
		Cards:              cards,
		CardsUploadedCount: item.CardsUploadedCount,
		RequiredCardsCount: item.RequiredCardsCount,
		CreatedAt:          item.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:          item.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
