package services

import (
	"context"
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
		Tags: models.LocalizedText{EN: "Hope, renewal.", RU: "Надежда, обновление"},
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
	if payload.MaxWords != 350 {
		t.Fatalf("expected max word budget of 350, got %d", payload.MaxWords)
	}
	if request.Model != "model-main" || request.MaxTokens < 2200 || len(request.Messages) != 2 {
		t.Fatalf("unexpected provider request: %#v", request)
	}
	userPrompt := request.Messages[1].Content
	for _, expected := range []string{"Звезда", "alice", "2026-07-07", "Europe/Moscow", "\"general\"", "relationships", "finances", "emotional_state", "career_business", "50"} {
		if !strings.Contains(userPrompt, expected) {
			t.Fatalf("expected prompt to contain %q, got %q", expected, userPrompt)
		}
	}

	_, englishRequest := buildCardOfDayAIRequest(card, CardOfDayPromptContext{Language: "en"}, "model-main")
	if !strings.Contains(englishRequest.Messages[0].Content, "no more than 20 words") || !strings.Contains(englishRequest.Messages[1].Content, "message") {
		t.Fatalf("expected English prompt to request the five long blocks, got %#v", englishRequest.Messages)
	}
}

func TestParseCardOfDayInterpretationValidatesJSON(t *testing.T) {
	valid := validCardOfDayTestInterpretation()
	raw, err := encodeCardOfDayInterpretation(valid)
	if err != nil {
		t.Fatalf("expected test interpretation to encode, got %v", err)
	}
	interpretation, err := parseCardOfDayInterpretation(raw)
	if err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if interpretation.General != valid.General || interpretation.Message != valid.Message || interpretation.Relationships == "" || interpretation.Finances == "" || interpretation.EmotionalState == "" || interpretation.CareerBusiness == "" {
		t.Fatalf("unexpected interpretation: %#v", interpretation)
	}

	if _, err := parseCardOfDayInterpretation(`{"general":"Only one field"}`); err == nil {
		t.Fatal("expected missing fields to be rejected")
	}
	if _, err := parseCardOfDayInterpretation("not json"); err == nil {
		t.Fatal("expected invalid json to be rejected")
	}
}

func TestParseCardOfDayInterpretationRequiresAtLeast50WordsInEveryExpandedBlock(t *testing.T) {
	shortBlockCases := []struct {
		name   string
		update func(*models.CardOfDayInterpretation)
	}{
		{name: "message", update: func(value *models.CardOfDayInterpretation) { value.Message = "one two three" }},
		{name: "relationships", update: func(value *models.CardOfDayInterpretation) { value.Relationships = "one two three" }},
		{name: "finances", update: func(value *models.CardOfDayInterpretation) { value.Finances = "one two three" }},
		{name: "emotional_state", update: func(value *models.CardOfDayInterpretation) { value.EmotionalState = "one two three" }},
		{name: "career_business", update: func(value *models.CardOfDayInterpretation) { value.CareerBusiness = "one two three" }},
	}

	for _, testCase := range shortBlockCases {
		t.Run(testCase.name, func(t *testing.T) {
			value := validCardOfDayTestInterpretation()
			testCase.update(&value)
			raw, err := encodeCardOfDayInterpretation(value)
			if err != nil {
				t.Fatalf("expected test interpretation to encode, got %v", err)
			}
			if _, err := parseCardOfDayInterpretation(raw); err == nil {
				t.Fatalf("expected short %s block to be rejected", testCase.name)
			} else {
				providerErr, ok := err.(*AIProviderError)
				if !ok || !strings.Contains(providerErr.Message, "at least 50 are required") {
					t.Fatalf("expected short %s block error message, got %#v", testCase.name, err)
				}
			}
		})
	}
}

func TestParseCardOfDayInterpretationLimitsGeneralTo20Words(t *testing.T) {
	interpretation := validCardOfDayTestInterpretation()
	interpretation.General = strings.TrimSpace(strings.Repeat("brief ", 21))
	raw, err := encodeCardOfDayInterpretation(interpretation)
	if err != nil {
		t.Fatalf("expected interpretation to encode, got %v", err)
	}
	if _, err := parseCardOfDayInterpretation(raw); err == nil {
		t.Fatal("expected general block longer than 20 words to be rejected")
	} else if providerErr, ok := err.(*AIProviderError); !ok || !strings.Contains(providerErr.Message, "no more than 20") {
		t.Fatalf("expected general length error message, got %#v", err)
	}
}

func TestParseCardOfDayInterpretationCountsRussianWords(t *testing.T) {
	block := strings.TrimSpace(strings.Repeat("\u0441\u043b\u043e\u0432\u043e ", 50))
	interpretation := models.CardOfDayInterpretation{
		General:        "\u043e\u0431\u0449\u0438\u0439 \u0432\u044b\u0432\u043e\u0434",
		Message:        block,
		Relationships:  block,
		Finances:       block,
		EmotionalState: block,
		CareerBusiness: block,
	}
	raw, err := encodeCardOfDayInterpretation(interpretation)
	if err != nil {
		t.Fatalf("expected Russian interpretation to encode, got %v", err)
	}
	if _, err := parseCardOfDayInterpretation(raw); err != nil {
		t.Fatalf("expected 50-word Russian blocks to be accepted, got %v", err)
	}
}

func TestGenerateCardOfDayInterpretationRetriesWhenAnyBlockIsTooShort(t *testing.T) {
	short := validCardOfDayTestInterpretation()
	short.Finances = "A brief financial reflection."
	shortJSON, err := encodeCardOfDayInterpretation(short)
	if err != nil {
		t.Fatalf("expected short interpretation to encode, got %v", err)
	}
	validJSON, err := encodeCardOfDayInterpretation(validCardOfDayTestInterpretation())
	if err != nil {
		t.Fatalf("expected valid interpretation to encode, got %v", err)
	}
	provider := &sequenceAIProvider{results: []AIProviderResult{
		{Content: shortJSON, FinishReason: "stop"},
		{Content: validJSON, FinishReason: "stop"},
	}}
	service := &AskService{provider: provider}

	result, interpretation, err := service.generateCardOfDayInterpretationWithRetry(context.Background(), AIProviderRequest{Language: "en", MaxTokens: 2200})
	if err != nil {
		t.Fatalf("expected short response to be retried, got %v", err)
	}
	if result == nil || interpretation.Finances != validCardOfDayTestInterpretation().Finances {
		t.Fatalf("expected valid retry response, got %#v", interpretation)
	}
	if len(provider.requests) != 2 || !strings.Contains(provider.requests[1].Messages[len(provider.requests[1].Messages)-1].Content, "at least 50 words") {
		t.Fatalf("expected corrective retry prompt, got %#v", provider.requests)
	}
}

func validCardOfDayTestInterpretation() models.CardOfDayInterpretation {
	block := strings.TrimSpace(strings.Repeat("reflection ", 50))
	return models.CardOfDayInterpretation{
		General:        "A concise overall reading.",
		Message:        block,
		Relationships:  block,
		Finances:       block,
		EmotionalState: block,
		CareerBusiness: block,
	}
}

type sequenceAIProvider struct {
	results  []AIProviderResult
	requests []AIProviderRequest
}

func (provider *sequenceAIProvider) GenerateInterpretation(_ context.Context, request AIProviderRequest) (*AIProviderResult, error) {
	provider.requests = append(provider.requests, request)
	if len(provider.results) == 0 {
		return nil, errors.New("unexpected extra provider call")
	}
	result := provider.results[0]
	provider.results = provider.results[1:]
	return &result, nil
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
