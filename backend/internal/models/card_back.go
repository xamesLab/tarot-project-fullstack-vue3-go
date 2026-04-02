package models

import "time"

type CardBack struct {
	ID          string        `bson:"_id" json:"id"`
	Name        LocalizedText `bson:"name" json:"name"`
	Description LocalizedText `bson:"description" json:"description"`
	Tags        []string      `bson:"tags" json:"tags"`
	Status      string        `bson:"status" json:"status"`
	ImagePath   string        `bson:"image_path" json:"image_path"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
}

type AdminCreateCardBackRequest struct {
	Name        LocalizedText
	Description LocalizedText
	Tags        []string
	ImagePath   string
}
