package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeckTheme struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	IsPlaceholder bool               `bson:"isPlaceholder"`
	Assets        map[string]string  `bson:"assets"`
	CreatedBy     primitive.ObjectID `bson:"createdBy"`
	CreatedAt     time.Time          `bson:"createdAt"`
}
