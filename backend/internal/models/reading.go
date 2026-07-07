package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ReadingTypeCardOfDay = "card_of_day"
	ReadingTypeThreeCard = "three_card"
)

type GenerateReadingRequest struct {
	Timezone string `json:"timezone" validate:"omitempty,timezone"`
}

type GenerateGuestCardOfDayRequest struct {
	GuestID  string `json:"guest_id" validate:"required,uuid4"`
	Timezone string `json:"timezone" validate:"required,timezone"`
}

type ResolveCardOfDayRequest struct {
	Timezone   string `json:"timezone" validate:"omitempty,timezone"`
	ClaimToken string `json:"claim_token"`
}

type GuestCardClaims struct {
	Purpose        string `json:"purpose"`
	GuestID        string `json:"guest_id"`
	CardKey        string `json:"card_key"`
	AccountingDate string `json:"accounting_date"`
	Timezone       string `json:"timezone"`
	jwt.RegisteredClaims
}

type ReadingCard struct {
	Key      string `bson:"key" json:"key"`
	Position int    `bson:"position" json:"position"`
	Reversed bool   `bson:"reversed" json:"reversed"`
	Viewed   bool   `bson:"viewed" json:"viewed"`
}

type ReadingMetadata struct {
	Seed      string `bson:"seed,omitempty" json:"seed,omitempty"`
	RequestIP string `bson:"request_ip,omitempty" json:"request_ip,omitempty"`
}

type AIInterpretation struct {
	RequestID       primitive.ObjectID `bson:"request_id,omitempty" json:"-"`
	Provider        string             `bson:"provider,omitempty" json:"provider,omitempty"`
	Model           string             `bson:"model,omitempty" json:"model,omitempty"`
	PromptVersion   string             `bson:"prompt_version,omitempty" json:"prompt_version,omitempty"`
	Language        string             `bson:"language,omitempty" json:"language,omitempty"`
	Status          string             `bson:"status,omitempty" json:"status,omitempty"`
	Content         string             `bson:"content,omitempty" json:"content,omitempty"`
	ErrorCode       string             `bson:"error_code,omitempty" json:"error_code,omitempty"`
	ErrorMessage    string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	CompletedAt     *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	LastRequestedAt *time.Time         `bson:"last_requested_at,omitempty" json:"last_requested_at,omitempty"`
	UpdatedAt       *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type AIInterpretationResponse struct {
	RequestID       string     `json:"request_id,omitempty"`
	Provider        string     `json:"provider,omitempty"`
	Model           string     `json:"model,omitempty"`
	PromptVersion   string     `json:"prompt_version,omitempty"`
	Language        string     `json:"language,omitempty"`
	Status          string     `json:"status,omitempty"`
	Content         string     `json:"content,omitempty"`
	ErrorCode       string     `json:"error_code,omitempty"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	LastRequestedAt *time.Time `json:"last_requested_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

type Reading struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID         primitive.ObjectID `bson:"user_id" json:"-"`
	Type           string             `bson:"type" json:"type"`
	AccountingDate string             `bson:"accounting_date" json:"accounting_date"`
	Timezone       string             `bson:"timezone" json:"timezone"`
	Cards          []ReadingCard      `bson:"cards" json:"cards"`
	DeckID         string             `bson:"deck_id" json:"deck_id"`
	BackID         string             `bson:"back_id,omitempty" json:"back_id,omitempty"`
	GeneratedAt    time.Time          `bson:"generated_at" json:"generated_at"`
	ExpiresAt      time.Time          `bson:"expires_at" json:"expires_at"`
	Metadata       ReadingMetadata    `bson:"metadata,omitempty" json:"metadata,omitempty"`
	AI             *AIInterpretation  `bson:"ai,omitempty" json:"ai,omitempty"`
}

type ReadingCardResponse struct {
	Key              string `json:"key"`
	Arcana           string `json:"arcana"`
	Type             string `json:"type"`
	TypeRU           string `json:"type_ru"`
	Name             string `json:"name"`
	NameRU           string `json:"name_ru"`
	Interpretation   string `json:"interpretation"`
	InterpretationRU string `json:"interpretation_ru"`
	Tags             string `json:"tags"`
	TagsRU           string `json:"tags_ru"`
	ImageURL         string `json:"image_url"`
	BackImageURL     string `json:"back_image_url,omitempty"`
	Position         int    `json:"position"`
	Reversed         bool   `json:"reversed"`
	Viewed           bool   `json:"viewed"`
}

type RemainingToday struct {
	Count int64 `json:"count"`
	Limit int64 `json:"limit"`
}

type ReadingResponse struct {
	ID              string                    `json:"id,omitempty"`
	Type            string                    `json:"type"`
	AccountingDate  string                    `json:"accounting_date"`
	GeneratedAt     time.Time                 `json:"generated_at"`
	Cards           []ReadingCardResponse     `json:"cards"`
	DeckID          string                    `json:"deck_id,omitempty"`
	BackID          string                    `json:"back_id,omitempty"`
	AI              *AIInterpretationResponse `json:"ai,omitempty"`
	NextAvailableAt *time.Time                `json:"next_available_at,omitempty"`
	RemainingToday  *RemainingToday           `json:"remaining_today,omitempty"`
}

type GuestCardOfDayResult struct {
	Reading    ReadingResponse `json:"reading"`
	ClaimToken string          `json:"claim_token"`
}

func (a *AIInterpretation) ToResponse() *AIInterpretationResponse {
	if a == nil {
		return nil
	}

	response := &AIInterpretationResponse{
		Provider:        a.Provider,
		Model:           a.Model,
		PromptVersion:   a.PromptVersion,
		Language:        a.Language,
		Status:          a.Status,
		Content:         a.Content,
		ErrorCode:       a.ErrorCode,
		ErrorMessage:    a.ErrorMessage,
		CompletedAt:     a.CompletedAt,
		LastRequestedAt: a.LastRequestedAt,
		UpdatedAt:       a.UpdatedAt,
	}

	if !a.RequestID.IsZero() {
		response.RequestID = a.RequestID.Hex()
	}

	return response
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
