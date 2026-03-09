package roles

import (
	"context"
	"errors"
	"sort"
	"strings"

	"tarotdaily/backend/internal/domain/models"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

const (
	DefaultUserDailySpreadLimit = 10
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user not found")
	ErrRolePolicyNotFound = errors.New("role policy not found")
)

type PolicyView struct {
	Role             string `json:"role"`
	DailySpreadLimit int    `json:"dailySpreadLimit"`
	IsUnlimited      bool   `json:"isUnlimited"`
}

type RoleAssignmentView struct {
	Login            string `json:"login"`
	Role             string `json:"role"`
	DailySpreadLimit int    `json:"dailySpreadLimit"`
	IsUnlimited      bool   `json:"isUnlimited"`
}

type Service struct {
	rolePolicies *mongorepo.RolePolicyRepository
	users        *mongorepo.UserRepository
}

func NewService(rolePolicies *mongorepo.RolePolicyRepository, users *mongorepo.UserRepository) *Service {
	return &Service{
		rolePolicies: rolePolicies,
		users:        users,
	}
}

func (s *Service) ResolveDailySpreadLimit(ctx context.Context, role string) (int, error) {
	normalizedRole := models.NormalizeRole(role)
	if normalizedRole == "" {
		return 0, ErrInvalidInput
	}

	policy, err := s.rolePolicies.FindByRole(ctx, normalizedRole)
	if err == nil {
		return policy.DailySpreadLimit, nil
	}
	if !errors.Is(err, mongorepo.ErrNotFound) {
		return 0, err
	}

	switch normalizedRole {
	case string(models.UserRoleAdmin):
		return models.UnlimitedDailySpreadLimit, nil
	case string(models.UserRoleUser):
		return DefaultUserDailySpreadLimit, nil
	default:
		return 0, ErrRolePolicyNotFound
	}
}

func (s *Service) ListPolicies(ctx context.Context) ([]PolicyView, error) {
	items, err := s.rolePolicies.List(ctx)
	if err != nil {
		return nil, err
	}

	index := map[string]PolicyView{}
	for _, item := range items {
		index[item.Role] = toPolicyView(item.Role, item.DailySpreadLimit)
	}

	if _, ok := index[string(models.UserRoleUser)]; !ok {
		index[string(models.UserRoleUser)] = toPolicyView(string(models.UserRoleUser), DefaultUserDailySpreadLimit)
	}
	if _, ok := index[string(models.UserRoleAdmin)]; !ok {
		index[string(models.UserRoleAdmin)] = toPolicyView(string(models.UserRoleAdmin), models.UnlimitedDailySpreadLimit)
	}

	out := make([]PolicyView, 0, len(index))
	rolesOrder := make([]string, 0, len(index))
	for role := range index {
		rolesOrder = append(rolesOrder, role)
	}
	sort.Strings(rolesOrder)
	for _, role := range rolesOrder {
		out = append(out, index[role])
	}
	return out, nil
}

func (s *Service) UpsertPolicy(ctx context.Context, role string, dailySpreadLimit int) (*PolicyView, error) {
	normalizedRole := models.NormalizeRole(role)
	if !models.IsValidRole(normalizedRole) {
		return nil, ErrInvalidInput
	}
	if dailySpreadLimit < models.UnlimitedDailySpreadLimit {
		return nil, ErrInvalidInput
	}

	item, err := s.rolePolicies.Upsert(ctx, normalizedRole, dailySpreadLimit)
	if err != nil {
		return nil, err
	}

	view := toPolicyView(item.Role, item.DailySpreadLimit)
	return &view, nil
}

func (s *Service) AssignRoleByLogin(ctx context.Context, login, role string) (*RoleAssignmentView, error) {
	normalizedLogin := strings.TrimSpace(login)
	normalizedRole := models.NormalizeRole(role)
	if normalizedLogin == "" || !models.IsValidRole(normalizedRole) {
		return nil, ErrInvalidInput
	}

	limit, err := s.ResolveDailySpreadLimit(ctx, normalizedRole)
	if err != nil {
		return nil, err
	}

	user, err := s.users.UpdateRoleByLogin(ctx, normalizedLogin, models.UserRole(normalizedRole))
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &RoleAssignmentView{
		Login:            user.Login,
		Role:             string(user.Role),
		DailySpreadLimit: limit,
		IsUnlimited:      limit == models.UnlimitedDailySpreadLimit,
	}, nil
}

func toPolicyView(role string, dailySpreadLimit int) PolicyView {
	return PolicyView{
		Role:             role,
		DailySpreadLimit: dailySpreadLimit,
		IsUnlimited:      dailySpreadLimit == models.UnlimitedDailySpreadLimit,
	}
}
