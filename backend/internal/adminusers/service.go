package adminusers

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"tarotdaily/backend/internal/domain/models"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type ListUsersInput struct {
	Page     int64
	PageSize int64
	Role     string
	Status   string
	Sort     string
	Name     string
}

type UserListItem struct {
	ID           string    `json:"id"`
	Login        string    `json:"login"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	RegisteredAt time.Time `json:"registeredAt"`
}

type ListUsersResult struct {
	Items      []UserListItem `json:"items"`
	Page       int64          `json:"page"`
	PageSize   int64          `json:"pageSize"`
	Total      int64          `json:"total"`
	TotalPages int64          `json:"totalPages"`
}

type Service struct {
	users *mongorepo.UserRepository
}

func NewService(users *mongorepo.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersResult, error) {
	page, pageSize, err := normalizePagination(input.Page, input.PageSize)
	if err != nil {
		return nil, ErrInvalidInput
	}

	role := models.NormalizeRole(input.Role)
	if role != "" && !models.IsValidRole(role) {
		return nil, ErrInvalidInput
	}

	status := strings.TrimSpace(strings.ToLower(input.Status))
	if status != "" && !isSupportedStatus(status) {
		return nil, ErrInvalidInput
	}

	sortAsc, err := normalizeSort(input.Sort)
	if err != nil {
		return nil, ErrInvalidInput
	}

	users, total, err := s.users.ListForAdmin(ctx, mongorepo.AdminUserListQuery{
		Page:     page,
		PageSize: pageSize,
		Role:     role,
		Status:   status,
		Name:     strings.TrimSpace(input.Name),
		SortAsc:  sortAsc,
	})
	if err != nil {
		return nil, err
	}

	items := make([]UserListItem, 0, len(users))
	for _, user := range users {
		items = append(items, UserListItem{
			ID:           user.ID.Hex(),
			Login:        user.Login,
			Name:         models.DisplayNameOrLogin(user.DisplayName, user.Login),
			Role:         string(user.Role),
			Status:       string(models.NormalizeUserStatus(user.Status)),
			RegisteredAt: user.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(pageSize)))
	}

	return &ListUsersResult{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func normalizePagination(page, pageSize int64) (int64, int64, error) {
	if page == 0 {
		page = defaultPage
	}
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	if page < 1 {
		return 0, 0, ErrInvalidInput
	}
	if pageSize < 1 || pageSize > maxPageSize {
		return 0, 0, ErrInvalidInput
	}
	return page, pageSize, nil
}

func normalizeSort(sortRaw string) (bool, error) {
	sort := strings.TrimSpace(strings.ToLower(sortRaw))
	switch sort {
	case "", "registeredat:desc", "desc":
		return false, nil
	case "registeredat:asc", "asc":
		return true, nil
	default:
		return false, ErrInvalidInput
	}
}

func isSupportedStatus(status string) bool {
	switch status {
	case string(models.UserStatusActive), string(models.UserStatusBanned), string(models.UserStatusFreeze):
		return true
	default:
		return false
	}
}
