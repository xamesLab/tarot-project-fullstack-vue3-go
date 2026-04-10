package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AskRequestStatusPending    = "pending"
	AskRequestStatusProcessing = "processing"
	AskRequestStatusCompleted  = "completed"
	AskRequestStatusFailed     = "failed"

	AskProviderOpenRouter = "openrouter"
)

type AskRequestCardInput struct {
	Key      string `bson:"key" json:"key"`
	Position int    `bson:"position,omitempty" json:"position,omitempty"`
}

type AskRequestPayload struct {
	ReadingType string                `bson:"reading_type" json:"reading_type"`
	Language    string                `bson:"language" json:"language"`
	Cards       []AskRequestCardInput `bson:"cards" json:"cards"`
	MaxWords    int                   `bson:"max_words" json:"max_words"`
}

type AskUsage struct {
	PromptTokens     int `bson:"prompt_tokens,omitempty" json:"prompt_tokens,omitempty"`
	CompletionTokens int `bson:"completion_tokens,omitempty" json:"completion_tokens,omitempty"`
	TotalTokens      int `bson:"total_tokens,omitempty" json:"total_tokens,omitempty"`
}

type AskRequest struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID          primitive.ObjectID `bson:"user_id" json:"-"`
	ReadingID       primitive.ObjectID `bson:"reading_id" json:"-"`
	Provider        string             `bson:"provider" json:"provider"`
	Model           string             `bson:"model" json:"model"`
	PromptVersion   string             `bson:"prompt_version" json:"prompt_version"`
	Status          string             `bson:"status" json:"status"`
	RequestPayload  AskRequestPayload  `bson:"request_payload" json:"request_payload"`
	ResponseData    map[string]any     `bson:"response_data,omitempty" json:"response_data,omitempty"`
	ResponseExcerpt string             `bson:"response_excerpt,omitempty" json:"response_excerpt,omitempty"`
	Usage           *AskUsage          `bson:"usage,omitempty" json:"usage,omitempty"`
	ErrorCode       string             `bson:"error_code,omitempty" json:"error_code,omitempty"`
	ErrorMessage    string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	RequestedAt     time.Time          `bson:"requested_at" json:"requested_at"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
	CompletedAt     *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}

type CreateAskRequest struct {
	ReadingID string `json:"reading_id" validate:"required"`
}

type AskRequestResponse struct {
	ID              string            `json:"id"`
	ReadingID       string            `json:"reading_id"`
	Provider        string            `json:"provider"`
	Model           string            `json:"model"`
	PromptVersion   string            `json:"prompt_version"`
	Status          string            `json:"status"`
	ResponseExcerpt string            `json:"response_excerpt,omitempty"`
	Usage           *AskUsage         `json:"usage,omitempty"`
	ErrorCode       string            `json:"error_code,omitempty"`
	ErrorMessage    string            `json:"error_message,omitempty"`
	RequestedAt     time.Time         `json:"requested_at"`
	CompletedAt     *time.Time        `json:"completed_at,omitempty"`
	RequestPayload  AskRequestPayload `json:"request_payload"`
}

func (r *AskRequest) ToResponse() AskRequestResponse {
	return AskRequestResponse{
		ID:              r.ID.Hex(),
		ReadingID:       r.ReadingID.Hex(),
		Provider:        r.Provider,
		Model:           r.Model,
		PromptVersion:   r.PromptVersion,
		Status:          r.Status,
		ResponseExcerpt: r.ResponseExcerpt,
		Usage:           r.Usage,
		ErrorCode:       r.ErrorCode,
		ErrorMessage:    r.ErrorMessage,
		RequestedAt:     r.RequestedAt,
		CompletedAt:     r.CompletedAt,
		RequestPayload:  r.RequestPayload,
	}
}
