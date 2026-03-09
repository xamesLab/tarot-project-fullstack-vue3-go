package auth

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/platform/timeutil"
)

func TestToUserViewNormalizesOptionalFields(t *testing.T) {
	t.Parallel()

	preferredDeckID := primitive.NewObjectID()
	deckThemeID := primitive.NewObjectID()

	view := toUserView(&models.User{
		ID:              primitive.NewObjectID(),
		Login:           "alice",
		Role:            models.UserRoleAdmin,
		Status:          "",
		DisplayName:     "  Alice  ",
		Timezone:        "bad/timezone",
		PreferredDeckID: &preferredDeckID,
		DeckThemeID:     &deckThemeID,
	})

	if view.Status != string(models.UserStatusActive) {
		t.Fatalf("unexpected status: %q", view.Status)
	}
	if view.DisplayName != "Alice" {
		t.Fatalf("unexpected displayName: %q", view.DisplayName)
	}
	if view.Timezone != timeutil.DefaultTimezone {
		t.Fatalf("unexpected timezone: %q", view.Timezone)
	}
	if view.PreferredDeckID != preferredDeckID.Hex() {
		t.Fatalf("unexpected preferredDeckId: %q", view.PreferredDeckID)
	}
	if view.DeckThemeID != deckThemeID.Hex() {
		t.Fatalf("unexpected deckThemeId: %q", view.DeckThemeID)
	}
}
