package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"tarot-daily/internal/models"
)

const cardOfDayPromptVersion = "v0.4.0-card-of-day-six-blocks-short-general"

var cardOfDayWordPattern = regexp.MustCompile(`\p{L}[\p{L}\p{N}]*(?:[-’'][\p{L}\p{N}]+)*`)

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
		MaxWords:    350,
	}

	var systemPrompt, userPrompt string
	if language == "ru" {
		systemPrompt = `Ты таролог NEBULA-TAROT. Отвечай только по-русски и верни только валидный JSON без markdown и пояснений. JSON должен содержать ровно шесть строковых полей: general, message, relationships, finances, emotional_state, career_business. Поле general — отдельный краткий общий вывод о карте, не более 20 слов. В каждом из остальных полей напиши не менее 50 слов; стремись к 50–65 словам в каждом блоке. Не смешивай general с message и не сокращай развёрнутые блоки ради общего объёма.`
		userPrompt = fmt.Sprintf(`Создай интерпретацию карты дня для конкретного человека.
Карта: %s
Ключ карты: %s
Аркан или масть: %s
Ключевые слова карты: %s
Базовое значение карты: %s
Контекст пользователя: %s
Дата учёта: %s
Часовой пояс: %s


Верни JSON в формате {"general":"краткий общий вывод о карте, не более 20 слов","message":"послание дня, не менее 50 слов","relationships":"интерпретация отношений, не менее 50 слов","finances":"интерпретация финансов, не менее 50 слов","emotional_state":"интерпретация душевного состояния, не менее 50 слов","career_business":"интерпретация карьеры или бизнеса, не менее 50 слов"}. Поле general должно суммировать основную тему карты в одном коротком выводе и не повторять полный текст послания. Каждый из остальных блоков — самостоятельный связный абзац из 50–65 слов с конкретной, бережной мыслью, связанной с картой. Не повторяй один совет во всех разделах. В финансовом блоке не обещай доход и не давай инвестиционных гарантий; в блоке душевного состояния не ставь диагнозы. Не делай медицинских или фаталистичных предсказаний.`,
			cardName,
			card.Key,
			cardType,
			cardTags,
			cardMeaning,
			userLabel,
			context.AccountingDate,
			context.Timezone,
		)
	} else {
		systemPrompt = `You are a NEBULA-TAROT reader. Respond only in English and return valid JSON without markdown or explanations. The JSON must contain exactly six string fields: general, message, relationships, finances, emotional_state, career_business. The general field is a separate concise overall card reading of no more than 20 words. Every other field must contain at least 50 words; aim for 50–65 words per block. Keep general distinct from message and do not shorten any long block to meet a total length.`
		userPrompt = fmt.Sprintf(`Create a personalized Card of the Day interpretation.
Card: %s
Card key: %s
Arcana or suit: %s
Card keywords: %s
Base card meaning: %s
User context: %s
Accounting date: %s
Timezone: %s

Return JSON in this shape: {"general":"concise overall card reading, no more than 20 words","message":"message of the day, at least 50 words","relationships":"relationships interpretation, at least 50 words","finances":"finances interpretation, at least 50 words","emotional_state":"emotional wellbeing interpretation, at least 50 words","career_business":"career or business interpretation, at least 50 words"}. The general field is one brief summary of the card's central theme and must remain distinct from the full message. Each other field must be a separate, coherent paragraph of 50–65 words with a specific, considerate insight connected to the card. Do not repeat the same advice across sections. In the finance block, do not promise income or give investment guarantees; in the emotional state block, do not diagnose. Avoid medical or fatalistic predictions.`,
			cardName,
			card.Key,
			cardType,
			cardTags,
			cardMeaning,
			userLabel,
			context.AccountingDate,
			context.Timezone,
		)
	}

	return payload, AIProviderRequest{
		Model:     model,
		Language:  language,
		MaxTokens: 2200,
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
	value.Message = normalizeInterpretation(value.Message)
	value.Relationships = normalizeInterpretation(value.Relationships)
	value.Finances = normalizeInterpretation(value.Finances)
	value.EmotionalState = normalizeInterpretation(value.EmotionalState)
	value.CareerBusiness = normalizeInterpretation(value.CareerBusiness)
	if value.General == "" || value.Message == "" || value.Relationships == "" || value.Finances == "" || value.EmotionalState == "" || value.CareerBusiness == "" {
		return value, &AIProviderError{
			Code:       "ai_provider_invalid_json",
			Message:    "AI provider returned incomplete JSON",
			StatusCode: 502,
		}
	}
	if count := len(cardOfDayWordPattern.FindAllString(value.General, -1)); count > 20 {
		return value, &AIProviderError{
			Code:       "ai_provider_invalid_json",
			Message:    fmt.Sprintf("AI provider returned general block with %d words; no more than 20 are allowed", count),
			StatusCode: 502,
		}
	}
	blocks := []struct {
		name  string
		value string
	}{
		{name: "message", value: value.Message},
		{name: "relationships", value: value.Relationships},
		{name: "finances", value: value.Finances},
		{name: "emotional_state", value: value.EmotionalState},
		{name: "career_business", value: value.CareerBusiness},
	}
	for _, block := range blocks {
		if count := len(cardOfDayWordPattern.FindAllString(block.value, -1)); count < 50 {
			return value, &AIProviderError{
				Code:       "ai_provider_invalid_json",
				Message:    fmt.Sprintf("AI provider returned %s block with %d words; at least 50 are required", block.name, count),
				StatusCode: 502,
			}
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
