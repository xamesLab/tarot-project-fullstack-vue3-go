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

func TestBuildCardOfDayAIRequestUsesTemplateData(t *testing.T) {
	card := models.Card{
		Key:    "the_star",
		Arcana: "major",
		Type:   models.LocalizedText{EN: "XVII", RU: "XVII"},
		Names:  models.LocalizedText{EN: "The Star", RU: "Звезда"},
		Interpretations: models.LocalizedText{
			EN: "Hope and renewal.",
			RU: "Надежда и обновление.",
		},
		Tags: models.LocalizedText{EN: "Hope, renewal.", RU: "Надежда, обновление."},
	}

	payload, request := buildCardOfDayAIRequest(card, CardOfDayPromptContext{
		Language:       "ru",
		Timezone:       "Europe/Moscow",
		AccountingDate: "2026-07-07",
		IsGuest:        false,
		Username:       "alice",
	}, "model-main")

	if payload.ReadingType != models.ReadingTypeCardOfDay || payload.Language != "ru" || payload.Cards[0].Key != "the_star" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if request.Model != "model-main" || request.MaxTokens == 0 || len(request.Messages) != 2 {
		t.Fatalf("unexpected provider request: %#v", request)
	}
	userPrompt := request.Messages[1].Content
	for _, expected := range []string{"Звезда", "alice", "2026-07-07", "Europe/Moscow", "\"general\""} {
		if !strings.Contains(userPrompt, expected) {
			t.Fatalf("expected prompt to contain %q, got %q", expected, userPrompt)
		}
	}
}

func TestParseCardOfDayInterpretationValidatesJSON(t *testing.T) {
	interpretation, err := parseCardOfDayInterpretation(`{"general":" One  two ","work":"Work","family":"Family","advice":"Advice"}`)
	if err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if interpretation.General != "One two" || interpretation.Work != "Work" {
		t.Fatalf("unexpected interpretation: %#v", interpretation)
	}

	if _, err := parseCardOfDayInterpretation(`{"general":"Only one field"}`); err == nil {
		t.Fatal("expected missing fields to be rejected")
	}
	if _, err := parseCardOfDayInterpretation("not json"); err == nil {
		t.Fatal("expected invalid json to be rejected")
	}
}

func TestGuestCardOfDayCacheKeyChangesByLanguageAndPromptVersion(t *testing.T) {
	ruKey := guestCardOfDayCacheKey("the_star", "ru")
	enKey := guestCardOfDayCacheKey("the_star", "en")
	if ruKey == enKey {
		t.Fatal("expected language to be part of cache key")
	}
	if !strings.Contains(ruKey, cardOfDayPromptVersion) {
		t.Fatalf("expected prompt version in cache key, got %s", ruKey)
	}
}
