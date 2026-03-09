package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UnlimitedDailySpreadLimit = -1

type RolePolicy struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Role             string             `bson:"role"`
	DailySpreadLimit int                `bson:"dailySpreadLimit"`
	CreatedAt        time.Time          `bson:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt"`
}
