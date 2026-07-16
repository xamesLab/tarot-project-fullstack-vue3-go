package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"tarot-daily/internal/models"
)

const cardOfDayPromptVersion = "v0.2.0-card-of-day-json"

type CardOfDayPromptContext struct {
	Language       string
	Timezone       string
	AccountingDate string
	IsGuest        bool
	Username       string
}

func buildCardOfDayAIRequest(card models.Card, context CardOfDayPromptContext, model string) (models.AskRequestPayload, AIProviderRequest) {
	language := normalizeLanguage(context.Language)
	cardName := localizedCardName(card, language)
	cardType := localizedText(card.Type, language)
	cardTags := localizedText(card.Tags, language)
	cardMeaning := localizedText(card.Interpretations, language)
	userLabel := "guest"
	if language == "ru" {
		userLabel = "гость"
	}
	if !context.IsGuest && strings.TrimSpace(context.Username) != "" {
		userLabel = strings.TrimSpace(context.Username)
	}

	payload := models.AskRequestPayload{
		ReadingType: models.ReadingTypeCardOfDay,
		Language:    language,
		Cards:       []models.AskRequestCardInput{{Key: card.Key, Position: 0}},
		MaxWords:    260,
	}

	if language == "ru" {
		systemPrompt := "Ты таролог NEBULA-TAROT. Отвечай только на русском языке. Верни только валидный JSON без markdown, комментариев и пояснений. JSON обязан содержать ровно четыре строковых поля: general, work, family, advice."
		userPrompt := fmt.Sprintf(
			"Сформируй интерпретацию карты дня.\nКарта: %s\nКлюч карты: %s\nАркан или масть: %s\nКлючевые слова: %s\nБазовое значение карты: %s\nКонтекст пользователя: %s\nДата учета: %s\nЧасовой пояс: %s\n\nВерни JSON в формате {\"general\":\"2-3 предложения общего смысла карты\",\"work\":\"совет в работе и карьере\",\"family\":\"совет в семье и отношениях\",\"advice\":\"короткий совет дня\"}. Пиши лично, мягко и мистично, но без медицинских, финансовых или фаталистичных гарантий.",
			cardName,
			card.Key,
			cardType,
			cardTags,
			cardMeaning,
			userLabel,
			context.AccountingDate,
			context.Timezone,
		)
		return payload, AIProviderRequest{
			Model:     model,
			Language:  language,
			MaxTokens: 760,
			Messages: []AIProviderMessage{
				{Role: "system", Content: systemPrompt},
				{Role: "user", Content: userPrompt},
			},
		}
	}

	systemPrompt := "You are a NEBULA-TAROT reader. Respond only in English. Return only valid JSON without markdown, comments, or explanations. The JSON must contain exactly four string fields: general, work, family, advice."
	userPrompt := fmt.Sprintf(
		"Create a Card of the Day interpretation.\nCard: %s\nCard key: %s\nArcana or suit: %s\nKeywords: %s\nCatalog meaning: %s\nUser context: %s\nAccounting date: %s\nTimezone: %s\n\nReturn JSON in this shape: {\"general\":\"2-3 sentences with the card's general meaning\",\"work\":\"guidance for work and career\",\"family\":\"guidance for family and relationships\",\"advice\":\"short advice for the day\"}. Write personally, gently, and mystically, without medical, financial, or fatalistic guarantees.",
		cardName,
		card.Key,
		cardType,
		cardTags,
		cardMeaning,
		userLabel,
		context.AccountingDate,
		context.Timezone,
	)

	return payload, AIProviderRequest{
		Model:     model,
		Language:  language,
		MaxTokens: 760,
		Messages: []AIProviderMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}
}

func parseCardOfDayInterpretation(raw string) (models.CardOfDayInterpretation, error) {
	var value models.CardOfDayInterpretation
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &value); err != nil {
		return value, &AIProviderError{
			Code:       "ai_provider_invalid_json",
			Message:    "AI provider returned invalid JSON",
			StatusCode: 502,
		}
	}

	value.General = normalizeInterpretation(value.General)
	value.Work = normalizeInterpretation(value.Work)
	value.Family = normalizeInterpretation(value.Family)
	value.Advice = normalizeInterpretation(value.Advice)
	if value.General == "" || value.Work == "" || value.Family == "" || value.Advice == "" {
		return value, &AIProviderError{
			Code:       "ai_provider_invalid_json",
			Message:    "AI provider returned incomplete JSON",
			StatusCode: 502,
		}
	}

	return value, nil
}

func encodeCardOfDayInterpretation(value models.CardOfDayInterpretation) (string, error) {
	body, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func localizedText(value models.LocalizedText, language string) string {
	if language == "ru" && strings.TrimSpace(value.RU) != "" {
		return value.RU
	}
	return value.EN
}
