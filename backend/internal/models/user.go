package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LocalizedText struct {
	EN string `bson:"en" json:"en"`
	RU string `bson:"ru" json:"ru"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Username       string             `bson:"username" json:"username"`
	Email          string             `bson:"email" json:"email"`
	PasswordHash   string             `bson:"password_hash" json:"-"`
	Role           string             `bson:"role" json:"role"`
	Status         string             `bson:"status" json:"status"`
	Timezone       string             `bson:"timezone" json:"timezone"`
	Language       string             `bson:"language" json:"language"`
	Theme          string             `bson:"theme" json:"theme"`
	SelectedDeckID string             `bson:"selected_deck_id" json:"selected_deck_id"`
	SelectedBackID string             `bson:"selected_back_id" json:"selected_back_id"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	LastLoginAt    *time.Time         `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
}

type UpdateProfileRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=20,alphanum"`
	Timezone string `json:"timezone" validate:"omitempty,timezone"`
}

type UpdateSettingsRequest struct {
	Language       string `json:"language" validate:"omitempty,oneof=en ru"`
	Theme          string `json:"theme" validate:"omitempty,oneof=light dark"`
	SelectedDeckID string `json:"selected_deck_id" validate:"omitempty"`
	SelectedBackID string `json:"selected_back_id" validate:"omitempty"`
}

type UserResponse struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	Timezone       string     `json:"timezone"`
	Language       string     `json:"language"`
	Theme          string     `json:"theme"`
	SelectedDeckID string     `json:"selected_deck_id,omitempty"`
	SelectedBackID string     `json:"selected_back_id,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
}

func (u User) ToResponse(includeMeta bool) UserResponse {
	response := UserResponse{
		ID:             u.ID.Hex(),
		Username:       u.Username,
		Email:          u.Email,
		Role:           u.Role,
		Timezone:       u.Timezone,
		Language:       u.Language,
		Theme:          u.Theme,
		SelectedDeckID: u.SelectedDeckID,
		SelectedBackID: u.SelectedBackID,
	}

	if includeMeta {
		response.CreatedAt = &u.CreatedAt
		response.LastLoginAt = u.LastLoginAt
	}

	return response
}
