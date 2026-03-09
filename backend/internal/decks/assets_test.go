package decks

import (
	"testing"

	"tarotdaily/backend/internal/domain/models"
)

func TestCalculateStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		backImageURL       string
		cardsCount         int
		requiredCardsCount int
		want               models.DeckStatus
	}{
		{
			name:               "done when back and all cards are present",
			backImageURL:       "/uploads/decks/1/back.jpg",
			cardsCount:         models.DefaultRequiredCardsCount,
			requiredCardsCount: models.DefaultRequiredCardsCount,
			want:               models.DeckStatusDone,
		},
		{
			name:               "process without back image",
			backImageURL:       "",
			cardsCount:         models.DefaultRequiredCardsCount,
			requiredCardsCount: models.DefaultRequiredCardsCount,
			want:               models.DeckStatusProcess,
		},
		{
			name:               "fallback required cards count",
			backImageURL:       "/uploads/decks/1/back.jpg",
			cardsCount:         models.DefaultRequiredCardsCount,
			requiredCardsCount: 0,
			want:               models.DeckStatusDone,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := calculateStatus(tt.backImageURL, tt.cardsCount, tt.requiredCardsCount)
			if got != tt.want {
				t.Fatalf("unexpected status: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeCardKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid minor arcana", input: "22", want: "22"},
		{name: "valid major arcana", input: "121", want: "121"},
		{name: "trim spaces", input: " 77 ", want: "77"},
		{name: "reject empty", input: "", wantErr: true},
		{name: "reject non numeric", input: "ace_of_wands", wantErr: true},
		{name: "reject out of catalog", input: "78", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeCardKey(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected card key: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestValidateImagePayload(t *testing.T) {
	t.Parallel()

	if err := validateImagePayload("image/jpeg", []byte{1, 2, 3}); err != nil {
		t.Fatalf("unexpected error for valid image payload: %v", err)
	}
	if err := validateImagePayload("text/plain", []byte{1, 2, 3}); err != ErrUnsupportedImage {
		t.Fatalf("unexpected error for unsupported image: %v", err)
	}
	if err := validateImagePayload("image/jpeg", nil); err != ErrInvalidInput {
		t.Fatalf("unexpected error for empty payload: %v", err)
	}
}
