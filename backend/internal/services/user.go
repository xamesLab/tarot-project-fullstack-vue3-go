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
	users *repositories.UserRepository
}

func NewUserService(users *repositories.UserRepository) *UserService {
	return &UserService{users: users}
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
	return s.users.UpdateSettings(ctx, id, req)
}
