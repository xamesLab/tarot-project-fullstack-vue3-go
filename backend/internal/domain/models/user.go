package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusBanned UserStatus = "banned"
	UserStatusFreeze UserStatus = "freeze"
)

type User struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty"`
	Login           string              `bson:"login"`
	PasswordHash    string              `bson:"passwordHash"`
	Role            UserRole            `bson:"role"`
	Status          UserStatus          `bson:"status"`
	DisplayName     string              `bson:"displayName,omitempty"`
	Timezone        string              `bson:"timezone"`
	PreferredDeckID *primitive.ObjectID `bson:"preferredDeckId,omitempty"`
	// DeckThemeID is kept for MVP compatibility and will be replaced by preferredDeckId.
	DeckThemeID *primitive.ObjectID `bson:"deckThemeId,omitempty"`
	CreatedAt   time.Time           `bson:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt"`
}
