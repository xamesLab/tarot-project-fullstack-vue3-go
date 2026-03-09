package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CardStateOpen   = "open"
	CardStateClosed = "closed"
)

type CardDraw struct {
	CardID primitive.ObjectID `bson:"cardId"`
	State  string             `bson:"state"`
}

type DailyReading struct {
	ID                     primitive.ObjectID  `bson:"_id,omitempty"`
	UserID                 primitive.ObjectID  `bson:"userId"`
	Date                   string              `bson:"date"`
	TimezoneUsed           string              `bson:"timezoneUsed"`
	DeckThemeID            *primitive.ObjectID `bson:"deckThemeId,omitempty"`
	CardOfDay              CardDraw            `bson:"cardOfDay"`
	SpreadA                []CardDraw          `bson:"spreadA"`
	SpreadB                []CardDraw          `bson:"spreadB"`
	SpreadGenerationsCount int                 `bson:"spreadGenerationsCount"`
	GeneratedAt            time.Time           `bson:"generatedAt"`
	LastSpreadGeneratedAt  *time.Time          `bson:"lastSpreadGeneratedAt,omitempty"`
	UpdatedAt              time.Time           `bson:"updatedAt"`
}
