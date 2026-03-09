package reading

import (
	"context"
	"crypto/rand"

	"tarotdaily/backend/internal/domain/models"
)

func (s *Service) drawCard(ctx context.Context) (models.CardDraw, error) {
	cards, err := s.cards.Sample(ctx, 1)
	if err != nil {
		return models.CardDraw{}, err
	}
	if len(cards) != 1 {
		return models.CardDraw{}, ErrCardsCatalogNotAvailable
	}

	state, err := randomCardState()
	if err != nil {
		return models.CardDraw{}, err
	}

	return models.CardDraw{
		CardID: cards[0].ID,
		State:  state,
	}, nil
}

func (s *Service) drawSpread(ctx context.Context, n int) ([]models.CardDraw, error) {
	cards, err := s.cards.Sample(ctx, n)
	if err != nil {
		return nil, err
	}
	if len(cards) != n {
		return nil, ErrCardsCatalogNotAvailable
	}

	out := make([]models.CardDraw, 0, n)
	for _, card := range cards {
		state, err := randomCardState()
		if err != nil {
			return nil, err
		}

		out = append(out, models.CardDraw{
			CardID: card.ID,
			State:  state,
		})
	}

	return out, nil
}

func randomCardState() (string, error) {
	var b [1]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	if b[0]%2 == 0 {
		return models.CardStateOpen, nil
	}
	return models.CardStateClosed, nil
}
