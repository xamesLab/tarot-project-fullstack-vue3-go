package services

import (
	"context"
	"errors"
	"strings"

	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	users     *repositories.UserRepository
	decks     *repositories.DeckRepository
	cardBacks *repositories.CardBackRepository
}

func NewUserService(
	users *repositories.UserRepository,
	decks *repositories.DeckRepository,
	cardBacks *repositories.CardBackRepository,
) *UserService {
	return &UserService{
		users:     users,
		decks:     decks,
		cardBacks: cardBacks,
	}
}

func (s *UserService) GetUserProfile(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.users.FindByID(ctx, id)
}

func (s *UserService) UpdateProfile(ctx context.Context, id primitive.ObjectID, req models.UpdateProfileRequest) (*models.User, error) {
	if req.Username != "" {
		existing, err := s.users.FindByUsername(ctx, strings.TrimSpace(req.Username))
		if err == nil && existing.ID != id {
			return nil, errors.Join(ErrValidationFailed, errors.New("Username unavailable"))
		}
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	}

	return s.users.UpdateProfile(ctx, id, strings.TrimSpace(req.Username), req.Timezone)
}

func (s *UserService) UpdateSettings(ctx context.Context, id primitive.ObjectID, req models.UpdateSettingsRequest) (*models.User, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.SelectedDeckID != nil || req.SelectedBackID != nil {
		if user.Role != "adept" && user.Role != "master" && user.Role != "keeper" {
			return nil, wrapValidation("Custom deck settings are available only for adept, master, and keeper roles")
		}
	}

	if req.SelectedDeckID != nil {
		deckID := strings.TrimSpace(*req.SelectedDeckID)
		if deckID != "" {
			deck, err := s.decks.FindByID(ctx, deckID)
			if err != nil {
				return nil, err
			}
			if deck.Status != "done" {
				return nil, wrapValidation("Selected deck is not available")
			}
		}
		req.SelectedDeckID = &deckID
	}

	if req.SelectedBackID != nil {
		backID := strings.TrimSpace(*req.SelectedBackID)
		if backID != "" {
			cardBack, err := s.cardBacks.FindByID(ctx, backID)
			if err != nil {
				return nil, err
			}
			if cardBack.Status != "done" {
				return nil, wrapValidation("Selected card back is not available")
			}
		}
		req.SelectedBackID = &backID
	}

	return s.users.UpdateSettings(ctx, id, req)
}
