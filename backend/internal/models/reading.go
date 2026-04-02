package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ReadingTypeCardOfDay = "card_of_day"
	ReadingTypeThreeCard = "three_card"
)

type GenerateReadingRequest struct {
	Timezone string `json:"timezone" validate:"omitempty,timezone"`
}

type ReadingCard struct {
	Key      string `bson:"key" json:"key"`
	Position int    `bson:"position" json:"position"`
	Reversed bool   `bson:"reversed" json:"reversed"`
}

type ReadingMetadata struct {
	Seed      string `bson:"seed,omitempty" json:"seed,omitempty"`
	RequestIP string `bson:"request_ip,omitempty" json:"request_ip,omitempty"`
}

type Reading struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID         primitive.ObjectID `bson:"user_id" json:"-"`
	Type           string             `bson:"type" json:"type"`
	AccountingDate string             `bson:"accounting_date" json:"accounting_date"`
	Timezone       string             `bson:"timezone" json:"timezone"`
	Cards          []ReadingCard      `bson:"cards" json:"cards"`
	DeckID         string             `bson:"deck_id" json:"deck_id"`
	GeneratedAt    time.Time          `bson:"generated_at" json:"generated_at"`
	ExpiresAt      time.Time          `bson:"expires_at" json:"expires_at"`
	Metadata       ReadingMetadata    `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type ReadingCardResponse struct {
	Key              string `json:"key"`
	Name             string `json:"name"`
	NameRU           string `json:"name_ru"`
	Interpretation   string `json:"interpretation"`
	InterpretationRU string `json:"interpretation_ru"`
	Tags             string `json:"tags"`
	TagsRU           string `json:"tags_ru"`
	ImageURL         string `json:"image_url"`
	Position         int    `json:"position"`
	Reversed         bool   `json:"reversed"`
}

type RemainingToday struct {
	Count int64 `json:"count"`
	Limit int64 `json:"limit"`
}

type ReadingResponse struct {
	ID              string                `json:"id"`
	Type            string                `json:"type"`
	AccountingDate  string                `json:"accounting_date"`
	GeneratedAt     time.Time             `json:"generated_at"`
	Cards           []ReadingCardResponse `json:"cards"`
	DeckID          string                `json:"deck_id,omitempty"`
	NextAvailableAt *time.Time            `json:"next_available_at,omitempty"`
	RemainingToday  *RemainingToday       `json:"remaining_today,omitempty"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type RetentionInfoResponse struct {
	Role          string `json:"role"`
	RetentionDays int    `json:"retention_days"`
	OldestDate    string `json:"oldest_date"`
}
