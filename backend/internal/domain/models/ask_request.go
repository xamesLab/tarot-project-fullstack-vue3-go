package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AskRequestStatus string

const (
	AskRequestStatusAccepted   AskRequestStatus = "accepted"
	AskRequestStatusQueued     AskRequestStatus = "queued"
	AskRequestStatusProcessing AskRequestStatus = "processing"
	AskRequestStatusSuccess    AskRequestStatus = "success"
	AskRequestStatusError      AskRequestStatus = "error"
)

type AskRequest struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	RequestID    string             `bson:"requestId"`
	UserID       primitive.ObjectID `bson:"userId"`
	Text         string             `bson:"text"`
	Status       AskRequestStatus   `bson:"status"`
	ResultText   string             `bson:"resultText,omitempty"`
	ErrorCode    string             `bson:"errorCode,omitempty"`
	ErrorMessage string             `bson:"errorMessage,omitempty"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
	CompletedAt  *time.Time         `bson:"completedAt,omitempty"`
}
