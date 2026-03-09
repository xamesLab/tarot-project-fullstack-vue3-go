package decks

import (
	"slices"
	"strconv"
	"strings"

	"tarotdaily/backend/internal/domain/models"
)

func calculateStatus(backImageURL string, cardsCount int, requiredCardsCount int) models.DeckStatus {
	required := requiredCardsCount
	if required <= 0 {
		required = models.DefaultRequiredCardsCount
	}
	if strings.TrimSpace(backImageURL) != "" && cardsCount == required {
		return models.DeckStatusDone
	}
	return models.DeckStatusProcess
}

func normalizeCardKey(raw string) (string, error) {
	key := strings.TrimSpace(raw)
	if key == "" {
		return "", ErrInvalidInput
	}
	if _, err := strconv.Atoi(key); err != nil {
		return "", ErrInvalidInput
	}
	if _, ok := validCardKeys[key]; !ok {
		return "", ErrInvalidInput
	}
	return key, nil
}

func normalizeExtFromContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

func validateImagePayload(contentType string, payload []byte) error {
	if len(payload) == 0 {
		return ErrInvalidInput
	}
	if len(payload) > MaxUploadBytes {
		return ErrUploadTooLarge
	}
	if !slices.Contains(allowedContentTypes, strings.ToLower(strings.TrimSpace(contentType))) {
		return ErrUnsupportedImage
	}
	return nil
}

func buildValidCardKeySet() map[string]struct{} {
	out := make(map[string]struct{}, models.DefaultRequiredCardsCount)
	for i := 22; i <= 77; i++ {
		out[strconv.Itoa(i)] = struct{}{}
	}
	for i := 100; i <= 121; i++ {
		out[strconv.Itoa(i)] = struct{}{}
	}
	return out
}
