package models

import "time"

type AdminSettings struct {
	ID                  string        `bson:"_id" json:"id"`
	MaintenanceMode     bool          `bson:"maintenance_mode" json:"maintenance_mode"`
	MaintenanceMessage  LocalizedText `bson:"maintenance_message" json:"maintenance_message"`
	DefaultRole         string        `bson:"default_role" json:"default_role"`
	DefaultDeckID       string        `bson:"default_deck_id" json:"default_deck_id"`
	DefaultBackID       string        `bson:"default_back_id" json:"default_back_id"`
	RegistrationEnabled bool          `bson:"registration_enabled" json:"registration_enabled"`
	CreatedAt           time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time     `bson:"updated_at" json:"updated_at"`
}

type AdminUpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=seeker adept master keeper"`
}

type AdminUpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active disabled"`
}

type AdminUpdateRolePolicyRequest struct {
	ThreeCardLimit       *int64 `json:"three_card_limit" validate:"required,gte=0"`
	HistoryRetentionDays *int   `json:"history_retention_days" validate:"required,gte=0"`
}

type AdminUpdateSettingsRequest struct {
	MaintenanceMode     *bool          `json:"maintenance_mode"`
	MaintenanceMessage  *LocalizedText `json:"maintenance_message"`
	DefaultRole         *string        `json:"default_role" validate:"omitempty,oneof=seeker adept master keeper"`
	DefaultDeckID       *string        `json:"default_deck_id"`
	DefaultBackID       *string        `json:"default_back_id"`
	RegistrationEnabled *bool          `json:"registration_enabled"`
}

type AdminUserRoleCounts struct {
	Seeker int64 `json:"seeker"`
	Adept  int64 `json:"adept"`
	Master int64 `json:"master"`
	Keeper int64 `json:"keeper"`
}

type AdminUserStatusCounts struct {
	Active   int64 `json:"active"`
	Disabled int64 `json:"disabled"`
}

type AdminUserMetrics struct {
	Total    int64                 `json:"total"`
	ByRole   AdminUserRoleCounts   `json:"by_role"`
	ByStatus AdminUserStatusCounts `json:"by_status"`
}

type AdminReadingPeriodMetrics struct {
	Total     int64 `json:"total"`
	CardOfDay int64 `json:"card_of_day"`
	ThreeCard int64 `json:"three_card"`
}

type AdminReadingMetrics struct {
	Today    AdminReadingPeriodMetrics `json:"today"`
	ThisWeek AdminReadingPeriodMetrics `json:"this_week"`
}

type AdminDeckStatusCounts struct {
	Done       int64 `json:"done"`
	Processing int64 `json:"processing"`
	Disabled   int64 `json:"disabled"`
	Archived   int64 `json:"archived"`
}

type AdminDeckMetrics struct {
	Total    int64                 `json:"total"`
	ByStatus AdminDeckStatusCounts `json:"by_status"`
}

type AdminDashboardMetrics struct {
	Users    AdminUserMetrics    `json:"users"`
	Readings AdminReadingMetrics `json:"readings"`
	Decks    AdminDeckMetrics    `json:"decks"`
}
