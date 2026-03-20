package reading

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/inpututil"
	"tarotdaily/backend/internal/platform/timeutil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
	"tarotdaily/backend/internal/roles"
)

func (s *Service) ensureDaily(
	ctx context.Context,
	userID string,
) (*models.User, string, string, *models.DailyReading, int, error) {
	oid, err := inpututil.ParseObjectID(userID)
	if err != nil {
		return nil, "", "", nil, 0, ErrInvalidInput
	}

	user, err := s.users.FindByID(ctx, oid)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, "", "", nil, 0, ErrUserNotFound
		}
		return nil, "", "", nil, 0, err
	}

	timezoneUsed, err := timeutil.NormalizeTimezone(user.Timezone)
	if err != nil {
		timezoneUsed = timeutil.DefaultTimezone
	}

	dateKey, err := timeutil.BusinessDateAtReset(time.Now().UTC(), timezoneUsed, dayResetHour)
	if err != nil {
		return nil, "", "", nil, 0, ErrInvalidInput
	}

	limit, err := s.roles.ResolveDailySpreadLimit(ctx, string(user.Role))
	if err != nil {
		if errors.Is(err, roles.ErrRolePolicyNotFound) {
			return nil, "", "", nil, 0, ErrRolePolicyNotFound
		}
		if errors.Is(err, roles.ErrInvalidInput) {
			return nil, "", "", nil, 0, ErrInvalidInput
		}
		return nil, "", "", nil, 0, err
	}

	daily, err := s.readings.FindByUserAndDate(ctx, user.ID, dateKey)
	if err != nil {
		if !errors.Is(err, mongorepo.ErrNotFound) {
			return nil, "", "", nil, 0, err
		}

		daily, err = s.createDailyReading(ctx, user, dateKey, timezoneUsed)
		if err != nil {
			return nil, "", "", nil, 0, err
		}
	}

	return user, dateKey, timezoneUsed, daily, limit, nil
}

func (s *Service) createDailyReading(
	ctx context.Context,
	user *models.User,
	dateKey string,
	timezoneUsed string,
) (*models.DailyReading, error) {
	cardOfDay, err := s.drawCard(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	item := &models.DailyReading{
		ID:                     primitive.NewObjectID(),
		UserID:                 user.ID,
		Date:                   dateKey,
		TimezoneUsed:           timezoneUsed,
		DeckThemeID:            resolveReadingDeckID(user),
		CardOfDay:              cardOfDay,
		SpreadA:                []models.CardDraw{},
		SpreadB:                []models.CardDraw{},
		SpreadGenerationsCount: 0,
		GeneratedAt:            now,
		UpdatedAt:              now,
	}

	if err := s.readings.Create(ctx, item); err != nil {
		if !mongorepo.IsDuplicateKeyError(err) {
			return nil, err
		}

		return s.readings.FindByUserAndDate(ctx, user.ID, dateKey)
	}

	return item, nil
}

func (s *Service) generateAndStoreSpreads(
	ctx context.Context,
	daily *models.DailyReading,
	currentDeckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	spreadA, err := s.drawSpread(ctx, 3)
	if err != nil {
		return nil, err
	}

	spreadB, err := s.drawSpread(ctx, 5)
	if err != nil {
		return nil, err
	}

	return s.readings.UpdateSpreads(ctx, daily.ID, spreadA, spreadB, currentDeckThemeID)
}

func (s *Service) generateCardOfDayOnly(
	ctx context.Context,
	daily *models.DailyReading,
	currentDeckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	cardOfDay, err := s.drawCard(ctx)
	if err != nil {
		return nil, err
	}

	return s.readings.UpdateCardOfDay(ctx, daily.ID, cardOfDay, currentDeckThemeID)
}

func (s *Service) generateSpreadAOnly(
	ctx context.Context,
	daily *models.DailyReading,
	currentDeckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	spreadA, err := s.drawSpread(ctx, 3)
	if err != nil {
		return nil, err
	}

	return s.readings.UpdateSpreadA(ctx, daily.ID, spreadA, currentDeckThemeID)
}

func (s *Service) generateSpreadBOnly(
	ctx context.Context,
	daily *models.DailyReading,
	currentDeckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	spreadB, err := s.drawSpread(ctx, 5)
	if err != nil {
		return nil, err
	}

	return s.readings.UpdateSpreadB(ctx, daily.ID, spreadB, currentDeckThemeID)
}

func resolveReadingDeckID(user *models.User) *primitive.ObjectID {
	if user == nil {
		return nil
	}
	if user.PreferredDeckID != nil {
		return user.PreferredDeckID
	}
	return user.DeckThemeID
}
