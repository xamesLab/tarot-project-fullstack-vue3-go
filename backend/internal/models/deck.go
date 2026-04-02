package models

import "time"

type DeckCardAsset struct {
	Key       string `bson:"key" json:"key"`
	ImagePath string `bson:"image_path" json:"image_path"`
}

type Deck struct {
	ID                string          `bson:"_id" json:"id"`
	Name              LocalizedText   `bson:"name" json:"name"`
	Description       LocalizedText   `bson:"description" json:"description"`
	Author            string          `bson:"author" json:"author"`
	Year              int             `bson:"year" json:"year"`
	Tags              []string        `bson:"tags" json:"tags"`
	Status            string          `bson:"status" json:"status"`
	CardCount         int             `bson:"card_count" json:"card_count"`
	BackImagePath     string          `bson:"back_image_path" json:"back_image_path"`
	CardImageBasePath string          `bson:"card_image_base_path" json:"card_image_base_path"`
	Cards             []DeckCardAsset `bson:"cards,omitempty" json:"cards,omitempty"`
	CreatedAt         time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time       `bson:"updated_at" json:"updated_at"`
}

type AdminCreateDeckRequest struct {
	Name        LocalizedText
	Description LocalizedText
	Author      string
	Year        int
	Tags        []string
	ZipPath     string
}

type AdminUpdateDeckRequest struct {
	Name        *LocalizedText `json:"name"`
	Description *LocalizedText `json:"description"`
	Author      *string        `json:"author"`
	Year        *int           `json:"year"`
	Tags        *[]string      `json:"tags"`
	Status      *string        `json:"status" validate:"omitempty,oneof=processing done disabled"`
}
