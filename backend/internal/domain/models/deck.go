package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeckStatus string

const (
	DeckStatusProcess DeckStatus = "process"
	DeckStatusDone    DeckStatus = "done"
)

const DefaultRequiredCardsCount = 78

type Deck struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Name               string             `bson:"name"`
	Status             DeckStatus         `bson:"status"`
	CreatedByUserID    primitive.ObjectID `bson:"createdByUserId"`
	BackImageURL       string             `bson:"backImageUrl,omitempty"`
	Cards              map[string]string  `bson:"cards,omitempty"`
	CardsUploadedCount int                `bson:"cardsUploadedCount"`
	RequiredCardsCount int                `bson:"requiredCardsCount"`
	CreatedAt          time.Time          `bson:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt"`
}
