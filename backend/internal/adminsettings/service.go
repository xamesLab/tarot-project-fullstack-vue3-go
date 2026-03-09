package adminsettings

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type SettingsView struct {
	ProjectTitle    string `json:"projectTitle"`
	SupportEmail    string `json:"supportEmail"`
	MaintenanceMode bool   `json:"maintenanceMode"`
	DayResetHour    int    `json:"dayResetHour"`
	MaxUploadMB     int    `json:"maxUploadMb"`
	UpdatedAt       string `json:"updatedAt"`
}

type UpdateInput struct {
	ProjectTitle    *string
	SupportEmail    *string
	MaintenanceMode *bool
	DayResetHour    *int
	MaxUploadMB     *int
}

type Service struct {
	repo *mongorepo.AdminSettingsRepository
}

func NewService(repo *mongorepo.AdminSettingsRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context) (*SettingsView, error) {
	item, err := s.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	view := toSettingsView(item.ProjectTitle, item.SupportEmail, item.MaintenanceMode, item.DayResetHour, item.MaxUploadMB, item.UpdatedAt)
	return &view, nil
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (*SettingsView, error) {
	if input.ProjectTitle == nil && input.SupportEmail == nil && input.MaintenanceMode == nil && input.DayResetHour == nil && input.MaxUploadMB == nil {
		return nil, ErrInvalidInput
	}

	setFields := bson.M{}

	if input.ProjectTitle != nil {
		title := strings.TrimSpace(*input.ProjectTitle)
		if len(title) < 2 || len([]rune(title)) > 120 {
			return nil, ErrInvalidInput
		}
		setFields["projectTitle"] = title
	}

	if input.SupportEmail != nil {
		email := strings.TrimSpace(*input.SupportEmail)
		if email != "" {
			if _, err := mail.ParseAddress(email); err != nil {
				return nil, ErrInvalidInput
			}
		}
		setFields["supportEmail"] = email
	}

	if input.MaintenanceMode != nil {
		setFields["maintenanceMode"] = *input.MaintenanceMode
	}

	if input.DayResetHour != nil {
		if *input.DayResetHour < 0 || *input.DayResetHour > 23 {
			return nil, ErrInvalidInput
		}
		setFields["dayResetHour"] = *input.DayResetHour
	}

	if input.MaxUploadMB != nil {
		if *input.MaxUploadMB < 1 || *input.MaxUploadMB > 100 {
			return nil, ErrInvalidInput
		}
		setFields["maxUploadMb"] = *input.MaxUploadMB
	}

	updated, err := s.repo.Update(ctx, setFields)
	if err != nil {
		return nil, err
	}

	view := toSettingsView(updated.ProjectTitle, updated.SupportEmail, updated.MaintenanceMode, updated.DayResetHour, updated.MaxUploadMB, updated.UpdatedAt)
	return &view, nil
}

func toSettingsView(projectTitle, supportEmail string, maintenanceMode bool, dayResetHour, maxUploadMB int, updatedAt time.Time) SettingsView {
	return SettingsView{
		ProjectTitle:    projectTitle,
		SupportEmail:    supportEmail,
		MaintenanceMode: maintenanceMode,
		DayResetHour:    dayResetHour,
		MaxUploadMB:     maxUploadMB,
		UpdatedAt:       updatedAt.UTC().Format(time.RFC3339),
	}
}
