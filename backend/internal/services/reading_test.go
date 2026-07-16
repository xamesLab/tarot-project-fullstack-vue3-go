package services

import (
	"reflect"
	"testing"
	"time"

	"tarot-daily/internal/models"
)

func TestSelectDeterministicCardKeyStable(t *testing.T) {
	keys := []string{"ace_of_cups", "strength", "the_fool", "the_magician"}
	seed := buildDeterministicSeed("user-1", "2026-04-02", "Europe/Moscow")

	first := selectDeterministicCardKey(keys, seed)
	second := selectDeterministicCardKey(keys, seed)

	if first == "" {
		t.Fatal("expected non-empty card key")
	}

	if first != second {
		t.Fatalf("expected deterministic selection, got %s and %s", first, second)
	}
}

func TestGuestCardOfDayPoolKeysStable(t *testing.T) {
	expected := []string{
		"the_fool",
		"the_magician",
		"the_high_priestess",
		"the_empress",
		"strength",
		"the_hermit",
		"justice",
		"temperance",
		"the_star",
		"the_sun",
	}
	if !reflect.DeepEqual(guestCardOfDayPoolKeys, expected) {
		t.Fatalf("unexpected guest card pool: %#v", guestCardOfDayPoolKeys)
	}
	for _, key := range expected {
		if !isGuestCardOfDayPoolKey(key) {
			t.Fatalf("expected %s to be allowed for guest card of day", key)
		}
	}
	if isGuestCardOfDayPoolKey("the_tower") {
		t.Fatal("expected cards outside the guest pool to be rejected")
	}
}

func TestSelectRandomUniqueCardKeysReturnsDistinctCards(t *testing.T) {
	keys := []string{"a", "b", "c", "d", "e"}

	selected, err := selectRandomUniqueCardKeys(keys, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(selected) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(selected))
	}

	seen := map[string]struct{}{}
	for _, key := range selected {
		if _, exists := seen[key]; exists {
			t.Fatalf("expected unique cards, duplicate %s found", key)
		}
		seen[key] = struct{}{}
	}
}

func TestCalculateRemainingCount(t *testing.T) {
	if got := calculateRemainingCount(5, 0); got != 4 {
		t.Fatalf("expected 4 remaining after first spread, got %d", got)
	}

	if got := calculateRemainingCount(5, 4); got != 0 {
		t.Fatalf("expected 0 remaining at limit edge, got %d", got)
	}

	if got := calculateRemainingCount(5, 8); got != 0 {
		t.Fatalf("expected 0 remaining when over limit, got %d", got)
	}
}

func TestBuildOldestRetentionDate(t *testing.T) {
	now := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)

	got, err := buildOldestRetentionDate(now, "Europe/Moscow", 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != "2026-03-26" {
		t.Fatalf("expected 2026-03-26, got %s", got)
	}
}

func TestUniqueCardKeys(t *testing.T) {
	readings := []models.Reading{
		{
			Cards: []models.ReadingCard{
				{Key: "the_fool"},
				{Key: "strength"},
			},
		},
		{
			Cards: []models.ReadingCard{
				{Key: "strength"},
				{Key: "the_magician"},
			},
		},
	}

	keys := uniqueCardKeys(readings)
	expected := []string{"strength", "the_fool", "the_magician"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}

	for index, key := range expected {
		if keys[index] != key {
			t.Fatalf("expected key %s at index %d, got %s", key, index, keys[index])
		}
	}
}

func TestBuildReadingCardResponseIncludesCatalogClassification(t *testing.T) {
	tests := []struct {
		name   string
		card   models.Card
		arcana string
		typeEN string
		typeRU string
	}{
		{name: "major", card: models.Card{Key: "the_high_priestess", Arcana: "major", Type: models.LocalizedText{EN: "II", RU: "II"}}, arcana: "major", typeEN: "II", typeRU: "II"},
		{name: "minor", card: models.Card{Key: "ace_of_wands", Arcana: "minor", Type: models.LocalizedText{EN: "Wands", RU: "Жезлы"}}, arcana: "minor", typeEN: "Wands", typeRU: "Жезлы"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := buildReadingCardResponse(test.card, models.ReadingCard{Position: 2, Reversed: true, Viewed: true}, "/face.png", "/back.png")
			if response.Arcana != test.arcana || response.Type != test.typeEN || response.TypeRU != test.typeRU {
				t.Fatalf("unexpected classification: %#v", response)
			}
			if response.Position != 2 || !response.Reversed || !response.Viewed {
				t.Fatalf("reading state was not preserved: %#v", response)
			}
		})
	}
}
