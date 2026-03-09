package auth

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/inpututil"
	"tarotdaily/backend/internal/platform/timeutil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

const maxDisplayNameLength = 64

func (s *Service) UpdateSettings(ctx context.Context, userID string, input UpdateSettingsInput) (*UserView, error) {
	userOID, err := inpututil.ParseObjectID(userID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	setFields := bson.M{}
	unsetFields := bson.M{}

	if input.DisplayName != nil {
		displayName := strings.TrimSpace(*input.DisplayName)
		if len(displayName) > maxDisplayNameLength {
			return nil, ErrInvalidInput
		}
		if displayName == "" {
			unsetFields["displayName"] = ""
		} else {
			setFields["displayName"] = displayName
		}
	}

	if input.Timezone != nil {
		timezone, err := timeutil.NormalizeTimezone(*input.Timezone)
		if err != nil {
			return nil, ErrInvalidInput
		}
		setFields["timezone"] = timezone
	}

	if input.PreferredDeckID != nil {
		deckID, err := inpututil.ParseOptionalObjectID(*input.PreferredDeckID)
		if err != nil {
			return nil, ErrInvalidInput
		}
		if deckID == nil {
			unsetFields["preferredDeckId"] = ""
		} else {
			deck, err := s.decks.FindByID(ctx, *deckID)
			if err != nil {
				if errors.Is(err, mongorepo.ErrNotFound) {
					return nil, ErrPreferredDeckNotUsable
				}
				return nil, err
			}
			if deck.Status != models.DeckStatusDone {
				return nil, ErrPreferredDeckNotUsable
			}
			setFields["preferredDeckId"] = *deckID
		}
	}

	if input.DeckThemeID != nil {
		themeID, err := inpututil.ParseOptionalObjectID(*input.DeckThemeID)
		if err != nil {
			return nil, ErrInvalidInput
		}
		if themeID == nil {
			unsetFields["deckThemeId"] = ""
		} else {
			if _, err := s.deckThemes.FindByID(ctx, *themeID); err != nil {
				if errors.Is(err, mongorepo.ErrNotFound) {
					return nil, ErrDeckThemeNotFound
				}
				return nil, err
			}
			setFields["deckThemeId"] = *themeID
		}
	}

	user, err := s.users.UpdateSettings(ctx, userOID, setFields, unsetFields)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	view := toUserView(user)
	return &view, nil
}
