package reading

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
)

func TestCollectCardIDs(t *testing.T) {
	cardOfDayID := primitive.NewObjectID()
	spreadAID := primitive.NewObjectID()
	spreadBID := primitive.NewObjectID()

	daily := &models.DailyReading{
		CardOfDay: models.CardDraw{CardID: cardOfDayID},
		SpreadA:   []models.CardDraw{{CardID: spreadAID}},
		SpreadB:   []models.CardDraw{{CardID: spreadBID}},
	}

	got := collectCardIDs(daily)
	if len(got) != 3 {
		t.Fatalf("collectCardIDs() len = %d, want 3", len(got))
	}
	if got[0] != cardOfDayID || got[1] != spreadAID || got[2] != spreadBID {
		t.Fatalf("collectCardIDs() order mismatch: got %v", got)
	}
}

func TestResolveCardImageURL(t *testing.T) {
	if got := resolveCardImageURL("22", map[string]string{"22": "/uploads/decks/a/cards/22.jpg"}); got != "/uploads/decks/a/cards/22.jpg" {
		t.Fatalf("resolveCardImageURL(asset) = %q", got)
	}
	if got := resolveCardImageURL("22", map[string]string{}); got != "/assets/cards/22.jpg" {
		t.Fatalf("resolveCardImageURL(numeric fallback) = %q", got)
	}
	if got := resolveCardImageURL("ace_of_wands", map[string]string{}); got != "/themes/default/ace_of_wands.jpg" {
		t.Fatalf("resolveCardImageURL(theme fallback) = %q", got)
	}
}

func TestApplyDailyLimit(t *testing.T) {
	view := &TodayView{CanGenerateMore: true}

	applyDailyLimit(view, 3, 2)
	if view.SpreadGenerationsLimit == nil || *view.SpreadGenerationsLimit != 3 {
		t.Fatalf("SpreadGenerationsLimit = %v, want 3", view.SpreadGenerationsLimit)
	}
	if view.SpreadGenerationsRemaining == nil || *view.SpreadGenerationsRemaining != 1 {
		t.Fatalf("SpreadGenerationsRemaining = %v, want 1", view.SpreadGenerationsRemaining)
	}
	if !view.CanGenerateMore {
		t.Fatalf("CanGenerateMore = false, want true")
	}

	unlimited := &TodayView{CanGenerateMore: true}
	applyDailyLimit(unlimited, models.UnlimitedDailySpreadLimit, 99)
	if unlimited.SpreadGenerationsLimit != nil || unlimited.SpreadGenerationsRemaining != nil {
		t.Fatalf("unlimited view should keep limit fields nil")
	}
}
