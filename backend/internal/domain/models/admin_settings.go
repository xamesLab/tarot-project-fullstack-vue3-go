package models

import "time"

const GlobalAdminSettingsID = "global"

type AdminSettings struct {
	ID              string    `bson:"_id"`
	ProjectTitle    string    `bson:"projectTitle"`
	SupportEmail    string    `bson:"supportEmail"`
	MaintenanceMode bool      `bson:"maintenanceMode"`
	DayResetHour    int       `bson:"dayResetHour"`
	MaxUploadMB     int       `bson:"maxUploadMb"`
	CreatedAt       time.Time `bson:"createdAt"`
	UpdatedAt       time.Time `bson:"updatedAt"`
}
