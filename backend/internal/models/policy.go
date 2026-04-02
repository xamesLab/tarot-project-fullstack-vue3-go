package models

import "time"

type RolePolicy struct {
	ID                   string        `bson:"_id" json:"id"`
	Role                 string        `bson:"role" json:"role"`
	DisplayName          LocalizedText `bson:"display_name" json:"display_name"`
	ThreeCardLimit       int64         `bson:"three_card_limit" json:"three_card_limit"`
	HistoryRetentionDays int           `bson:"history_retention_days" json:"history_retention_days"`
	CreatedAt            time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time     `bson:"updated_at" json:"updated_at"`
}
