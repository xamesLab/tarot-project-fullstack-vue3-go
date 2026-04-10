package services

import (
	"errors"
	"strings"
	"testing"

	"tarot-daily/internal/models"
)

func TestEnsureReadingReadyForAICardOfDay(t *testing.T) {
	reading := &models.Reading{
		Type: models.ReadingTypeCardOfDay,
		Cards: []models.ReadingCard{
			{Key: "the_fool", Position: 0, Viewed: true},
		},
	}

	cards, err := ensureReadingReadyForAI(reading)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cards) != 1 || cards[0].Key != "the_fool" {
		t.Fatalf("expected revealed card of day, got %#v", cards)
	}
}

func TestEnsureReadingReadyForAIRequiresReveal(t *testing.T) {
	reading := &models.Reading{
		Type: models.ReadingTypeCardOfDay,
		Cards: []models.ReadingCard{
			{Key: "the_fool", Position: 0, Viewed: false},
		},
	}

	_, err := ensureReadingReadyForAI(reading)
	if !errors.Is(err, ErrValidationFailed) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestEnsureReadingReadyForAIThreeCardSortsByPosition(t *testing.T) {
	reading := &models.Reading{
		Type: models.ReadingTypeThreeCard,
		Cards: []models.ReadingCard{
			{Key: "future", Position: 2, Viewed: true},
			{Key: "past", Position: 0, Viewed: true},
			{Key: "present", Position: 1, Viewed: true},
		},
	}

	cards, err := ensureReadingReadyForAI(reading)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cards) != 3 {
		t.Fatalf("expected 3 cards, got %d", len(cards))
	}
	if cards[0].Key != "past" || cards[1].Key != "present" || cards[2].Key != "future" {
		t.Fatalf("expected cards sorted by position, got %#v", cards)
	}
}

func TestNormalizeInterpretationNormalizesWhitespace(t *testing.T) {
	value := "  one   two\nthree \t four  "
	got := normalizeInterpretation(value)

	if strings.Contains(got, "\n") || strings.Contains(got, "\t") || strings.Contains(got, "  ") {
		t.Fatalf("expected normalized whitespace, got %q", got)
	}
	if got != "one two three four" {
		t.Fatalf("unexpected normalized output %q", got)
	}
}

func TestNormalizeInterpretationDoesNotTrimWords(t *testing.T) {
	value := "one two three four five six"
	got := normalizeInterpretation(value)

	if got != "one two three four five six" {
		t.Fatalf("expected full interpretation without truncation, got %q", got)
	}
}

func TestRecommendedMaxTokensForThreeCard(t *testing.T) {
	if got := recommendedMaxTokens(models.ReadingTypeThreeCard); got != 520 {
		t.Fatalf("expected 520 tokens for three-card reading, got %d", got)
	}
}

func TestShouldReuseAskRequestRejectsLengthFinishReason(t *testing.T) {
	request := &models.AskRequest{
		ResponseExcerpt: "partial text",
		ResponseData: map[string]any{
			"finish_reason": "length",
		},
	}

	if shouldReuseAskRequest(request) {
		t.Fatal("expected truncated request to be regenerated instead of reused")
	}
}
