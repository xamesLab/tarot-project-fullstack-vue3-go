package models

import "time"

type Card struct {
	ID              string        `bson:"_id" json:"id"`
	OriginalID      int           `bson:"id" json:"original_id"`
	Key             string        `bson:"key" json:"key"`
	Arcana          string        `bson:"arcana" json:"arcana"`
	Type            LocalizedText `bson:"type" json:"type"`
	Names           LocalizedText `bson:"names" json:"names"`
	Interpretations LocalizedText `bson:"interpretations" json:"interpretations"`
	Tags            LocalizedText `bson:"tags" json:"tags"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
}
