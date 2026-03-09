package auth

import (
	"strings"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/timeutil"
)

func toUserView(user *models.User) UserView {
	view := UserView{
		ID:          user.ID.Hex(),
		Login:       user.Login,
		Role:        string(user.Role),
		Status:      string(models.NormalizeUserStatus(user.Status)),
		DisplayName: strings.TrimSpace(user.DisplayName),
		Timezone:    normalizeUserTimezone(user.Timezone),
	}

	if user.PreferredDeckID != nil {
		view.PreferredDeckID = user.PreferredDeckID.Hex()
	}
	if user.DeckThemeID != nil {
		view.DeckThemeID = user.DeckThemeID.Hex()
	}

	return view
}

func normalizeUserTimezone(timezone string) string {
	normalized, err := timeutil.NormalizeTimezone(timezone)
	if err != nil {
		return timeutil.DefaultTimezone
	}
	return normalized
}
