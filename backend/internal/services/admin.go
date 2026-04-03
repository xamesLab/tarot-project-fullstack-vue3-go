package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	users     *repositories.UserRepository
	cards     *repositories.CardRepository
	readings  *repositories.ReadingRepository
	policies  *repositories.PolicyRepository
	settings  *repositories.AdminSettingsRepository
	decks     *repositories.DeckRepository
	cardBacks *repositories.CardBackRepository
}

func NewAdminService(
	users *repositories.UserRepository,
	cards *repositories.CardRepository,
	readings *repositories.ReadingRepository,
	policies *repositories.PolicyRepository,
	settings *repositories.AdminSettingsRepository,
	decks *repositories.DeckRepository,
	cardBacks *repositories.CardBackRepository,
) *AdminService {
	return &AdminService{
		users:     users,
		cards:     cards,
		readings:  readings,
		policies:  policies,
		settings:  settings,
		decks:     decks,
		cardBacks: cardBacks,
	}
}

func (s *AdminService) GetDashboardMetrics(ctx context.Context) (*models.AdminDashboardMetrics, error) {
	var metrics models.AdminDashboardMetrics
	var err error

	metrics.Users.Total, err = s.users.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	metrics.Users.ByRole.Seeker, err = s.users.CountByRole(ctx, "seeker")
	if err != nil {
		return nil, err
	}
	metrics.Users.ByRole.Adept, err = s.users.CountByRole(ctx, "adept")
	if err != nil {
		return nil, err
	}
	metrics.Users.ByRole.Master, err = s.users.CountByRole(ctx, "master")
	if err != nil {
		return nil, err
	}
	metrics.Users.ByRole.Keeper, err = s.users.CountByRole(ctx, "keeper")
	if err != nil {
		return nil, err
	}

	metrics.Users.ByStatus.Active, err = s.users.CountByStatus(ctx, "active")
	if err != nil {
		return nil, err
	}
	metrics.Users.ByStatus.Disabled, err = s.users.CountByStatus(ctx, "disabled")
	if err != nil {
		return nil, err
	}

	startOfToday := time.Now().UTC().Truncate(24 * time.Hour)
	metrics.Readings.Today.Total, metrics.Readings.Today.CardOfDay, metrics.Readings.Today.ThreeCard, err = s.readings.CountSince(ctx, startOfToday)
	if err != nil {
		return nil, err
	}

	startOfWeek := startOfToday.AddDate(0, 0, -6)
	metrics.Readings.ThisWeek.Total, metrics.Readings.ThisWeek.CardOfDay, metrics.Readings.ThisWeek.ThreeCard, err = s.readings.CountSince(ctx, startOfWeek)
	if err != nil {
		return nil, err
	}

	metrics.Decks.Total, err = s.decks.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	metrics.Decks.ByStatus.Done, err = s.decks.CountByStatus(ctx, "done")
	if err != nil {
		return nil, err
	}
	metrics.Decks.ByStatus.Processing, err = s.decks.CountByStatus(ctx, "processing")
	if err != nil {
		return nil, err
	}
	metrics.Decks.ByStatus.Disabled, err = s.decks.CountByStatus(ctx, "disabled")
	if err != nil {
		return nil, err
	}
	metrics.Decks.ByStatus.Archived, err = s.decks.CountByStatus(ctx, "archived")
	if err != nil {
		return nil, err
	}

	return &metrics, nil
}

func (s *AdminService) ListUsers(ctx context.Context, search, role, status string, page, limit int) ([]models.UserResponse, models.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	role = strings.TrimSpace(role)
	status = strings.TrimSpace(status)
	skip := int64((page - 1) * limit)

	users, err := s.users.ListAdminUsers(ctx, search, role, status, skip, int64(limit))
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	total, err := s.users.CountAdminUsers(ctx, search, role, status)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	responseUsers := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		responseUsers = append(responseUsers, user.ToResponse(true))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return responseUsers, models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) UpdateUserRole(ctx context.Context, actorID, targetID primitive.ObjectID, role string) (*models.User, error) {
	if actorID == targetID && role != "keeper" {
		return nil, errors.Join(ErrValidationFailed, errors.New("You cannot remove your own keeper role"))
	}

	return s.users.UpdateRole(ctx, targetID, role)
}

func (s *AdminService) UpdateUserStatus(ctx context.Context, actorID, targetID primitive.ObjectID, status string) (*models.User, error) {
	if actorID == targetID && status != "active" {
		return nil, errors.Join(ErrValidationFailed, errors.New("You cannot disable your own account"))
	}

	return s.users.UpdateStatus(ctx, targetID, status)
}

func (s *AdminService) ListRolePolicies(ctx context.Context) ([]models.RolePolicy, error) {
	return s.policies.ListAll(ctx)
}

func (s *AdminService) UpdateRolePolicy(ctx context.Context, role string, req models.AdminUpdateRolePolicyRequest) (*models.RolePolicy, error) {
	if req.ThreeCardLimit == nil || req.HistoryRetentionDays == nil {
		return nil, errors.Join(ErrValidationFailed, errors.New("Both policy fields are required"))
	}

	return s.policies.UpdateByRole(ctx, role, *req.ThreeCardLimit, *req.HistoryRetentionDays)
}

func (s *AdminService) GetSettings(ctx context.Context) (*models.AdminSettings, error) {
	return s.settings.Get(ctx)
}

func (s *AdminService) UpdateSettings(ctx context.Context, req models.AdminUpdateSettingsRequest) (*models.AdminSettings, error) {
	if req.DefaultRole != nil {
		role := strings.TrimSpace(*req.DefaultRole)
		req.DefaultRole = &role
	}
	if req.DefaultDeckID != nil {
		deckID := strings.TrimSpace(*req.DefaultDeckID)
		if deckID == "" {
			return nil, wrapValidation("Default deck is required")
		}
		deck, err := s.decks.FindByID(ctx, deckID)
		if err != nil {
			return nil, err
		}
		if deck.Status != "done" {
			return nil, wrapValidation("Default deck must be available")
		}
		req.DefaultDeckID = &deckID
	}
	if req.DefaultBackID != nil {
		backID := strings.TrimSpace(*req.DefaultBackID)
		if backID == "" {
			return nil, wrapValidation("Default card back is required")
		}
		cardBack, err := s.cardBacks.FindByID(ctx, backID)
		if err != nil {
			return nil, err
		}
		if cardBack.Status != "done" {
			return nil, wrapValidation("Default card back must be available")
		}
		req.DefaultBackID = &backID
	}

	return s.settings.Update(ctx, req)
}

func (s *AdminService) ListDecks(ctx context.Context, scope string, page, limit int) ([]models.Deck, models.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	skip := int64((page - 1) * limit)
	statuses, err := resolveAssetScopeStatuses(scope)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	decks, err := s.decks.ListByStatuses(ctx, statuses, skip, int64(limit))
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	total, err := s.decks.CountByStatuses(ctx, statuses)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return decks, models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) CreateDeck(ctx context.Context, req models.AdminCreateDeckRequest) (*models.Deck, error) {
	nameEN := strings.TrimSpace(req.Name.EN)
	nameRU := strings.TrimSpace(req.Name.RU)
	if nameEN == "" && nameRU == "" {
		return nil, wrapValidation("Deck name is required")
	}
	if nameEN == "" {
		nameEN = nameRU
	}
	if nameRU == "" {
		nameRU = nameEN
	}

	cardKeys, err := s.cards.ListAllKeys(ctx)
	if err != nil {
		return nil, err
	}
	if len(cardKeys) != 78 {
		return nil, wrapValidation("Card catalog is incomplete")
	}

	deckID := primitive.NewObjectID().Hex()
	now := time.Now().UTC()
	deck := &models.Deck{
		ID: deckID,
		Name: models.LocalizedText{
			EN: nameEN,
			RU: nameRU,
		},
		Description: models.LocalizedText{
			EN: strings.TrimSpace(req.Description.EN),
			RU: strings.TrimSpace(req.Description.RU),
		},
		Author:            strings.TrimSpace(req.Author),
		Year:              req.Year,
		Tags:              sanitizeTags(req.Tags),
		Status:            "processing",
		CardCount:         0,
		BackImagePath:     "",
		CardImageBasePath: fmt.Sprintf("/uploads/decks/%s/cards/", deckID),
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.decks.Create(ctx, deck); err != nil {
		return nil, err
	}

	destPath := filepath.Join("uploads", "decks", deckID)
	manifest, err := utils.ValidateAndExtractDeckZIP(req.ZipPath, destPath, cardKeys)
	if err != nil {
		_ = s.decks.Delete(ctx, deckID)
		_ = os.RemoveAll(destPath)
		return nil, wrapValidation(err.Error())
	}

	assets := make([]models.DeckCardAsset, 0, len(manifest.CardKeys))
	for _, key := range manifest.CardKeys {
		assets = append(assets, models.DeckCardAsset{
			Key:       key,
			ImagePath: fmt.Sprintf("/uploads/decks/%s/cards/%s.png", deckID, key),
		})
	}

	updatedDeck, err := s.decks.Update(ctx, deckID, bson.M{
		"status":               "done",
		"card_count":           len(assets),
		"card_image_base_path": fmt.Sprintf("/uploads/decks/%s/cards/", deckID),
		"cards":                assets,
	})
	if err != nil {
		_ = os.RemoveAll(destPath)
		_ = s.decks.Delete(ctx, deckID)
		return nil, err
	}

	return updatedDeck, nil
}

func (s *AdminService) ListCardBacks(ctx context.Context, scope string, page, limit int) ([]models.CardBack, models.PaginationResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	skip := int64((page - 1) * limit)
	statuses, err := resolveAssetScopeStatuses(scope)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	cardBacks, err := s.cardBacks.ListByStatuses(ctx, statuses, skip, int64(limit))
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	total, err := s.cardBacks.CountByStatuses(ctx, statuses)
	if err != nil {
		return nil, models.PaginationResponse{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return cardBacks, models.PaginationResponse{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *AdminService) CreateCardBack(ctx context.Context, req models.AdminCreateCardBackRequest) (*models.CardBack, error) {
	nameEN := strings.TrimSpace(req.Name.EN)
	nameRU := strings.TrimSpace(req.Name.RU)
	if nameEN == "" && nameRU == "" {
		return nil, wrapValidation("Card back name is required")
	}
	if nameEN == "" {
		nameEN = nameRU
	}
	if nameRU == "" {
		nameRU = nameEN
	}

	cardBackID := primitive.NewObjectID().Hex()
	now := time.Now().UTC()
	cardBack := &models.CardBack{
		ID: cardBackID,
		Name: models.LocalizedText{
			EN: nameEN,
			RU: nameRU,
		},
		Description: models.LocalizedText{
			EN: strings.TrimSpace(req.Description.EN),
			RU: strings.TrimSpace(req.Description.RU),
		},
		Tags:      sanitizeTags(req.Tags),
		Status:    "processing",
		ImagePath: fmt.Sprintf("/uploads/card-backs/%s/back.png", cardBackID),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.cardBacks.Create(ctx, cardBack); err != nil {
		return nil, err
	}

	destDir := filepath.Join("uploads", "card-backs", cardBackID)
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		_ = s.cardBacks.Delete(ctx, cardBackID)
		return nil, err
	}

	if err := utils.ValidateAndSaveBackImage(req.ImagePath, filepath.Join(destDir, "back.png")); err != nil {
		_ = s.cardBacks.Delete(ctx, cardBackID)
		_ = os.RemoveAll(destDir)
		return nil, wrapValidation(err.Error())
	}

	updatedCardBack, err := s.cardBacks.Update(ctx, cardBackID, bson.M{
		"status":     "done",
		"image_path": fmt.Sprintf("/uploads/card-backs/%s/back.png", cardBackID),
	})
	if err != nil {
		_ = s.cardBacks.Delete(ctx, cardBackID)
		_ = os.RemoveAll(destDir)
		return nil, err
	}

	return updatedCardBack, nil
}

func (s *AdminService) GetDeck(ctx context.Context, deckID string) (*models.Deck, error) {
	return s.decks.FindByID(ctx, strings.TrimSpace(deckID))
}

func (s *AdminService) UpdateDeckMetadata(ctx context.Context, deckID string, req models.AdminUpdateDeckRequest) (*models.Deck, error) {
	deckID = strings.TrimSpace(deckID)
	if deckID == "" {
		return nil, wrapValidation("Deck id is required")
	}

	updates := bson.M{}
	if req.Name != nil {
		nameEN := strings.TrimSpace(req.Name.EN)
		nameRU := strings.TrimSpace(req.Name.RU)
		if nameEN == "" && nameRU == "" {
			return nil, wrapValidation("Deck name is required")
		}
		if nameEN == "" {
			nameEN = nameRU
		}
		if nameRU == "" {
			nameRU = nameEN
		}
		updates["name"] = models.LocalizedText{EN: nameEN, RU: nameRU}
	}
	if req.Description != nil {
		updates["description"] = models.LocalizedText{
			EN: strings.TrimSpace(req.Description.EN),
			RU: strings.TrimSpace(req.Description.RU),
		}
	}
	if req.Author != nil {
		updates["author"] = strings.TrimSpace(*req.Author)
	}
	if req.Year != nil {
		if *req.Year < 0 {
			return nil, wrapValidation("Year must be a positive number")
		}
		updates["year"] = *req.Year
	}
	if req.Tags != nil {
		updates["tags"] = sanitizeTags(*req.Tags)
	}
	if req.Status != nil {
		updates["status"] = strings.TrimSpace(*req.Status)
	}
	if len(updates) == 0 {
		return s.decks.FindByID(ctx, deckID)
	}

	return s.decks.Update(ctx, deckID, updates)
}

func (s *AdminService) ArchiveDeck(ctx context.Context, deckID string) (*models.Deck, error) {
	deckID = strings.TrimSpace(deckID)
	if deckID == "" {
		return nil, wrapValidation("Deck id is required")
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	if settings.DefaultDeckID == deckID {
		return nil, wrapValidation("Cannot archive the current default deck")
	}

	return s.decks.Update(ctx, deckID, bson.M{"status": "archived"})
}

func (s *AdminService) RestoreDeck(ctx context.Context, deckID string) (*models.Deck, error) {
	deckID = strings.TrimSpace(deckID)
	if deckID == "" {
		return nil, wrapValidation("Deck id is required")
	}

	deck, err := s.decks.FindByID(ctx, deckID)
	if err != nil {
		return nil, err
	}
	if deck.Status != "archived" {
		return deck, nil
	}

	return s.decks.Update(ctx, deckID, bson.M{"status": "done"})
}

func (s *AdminService) DeleteDeck(ctx context.Context, deckID string) error {
	deckID = strings.TrimSpace(deckID)
	if deckID == "" {
		return wrapValidation("Deck id is required")
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return err
	}
	if settings.DefaultDeckID == deckID {
		return wrapValidation("Cannot delete the current default deck")
	}

	deck, err := s.decks.FindByID(ctx, deckID)
	if err != nil {
		return err
	}
	if deck.Status != "archived" {
		return wrapValidation("Only archived decks can be deleted permanently")
	}

	if err := s.decks.Delete(ctx, deckID); err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join("uploads", "decks", deckID))
}

func (s *AdminService) ArchiveCardBack(ctx context.Context, cardBackID string) (*models.CardBack, error) {
	cardBackID = strings.TrimSpace(cardBackID)
	if cardBackID == "" {
		return nil, wrapValidation("Card back id is required")
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	if settings.DefaultBackID == cardBackID {
		return nil, wrapValidation("Cannot archive the current default card back")
	}

	return s.cardBacks.Update(ctx, cardBackID, bson.M{"status": "archived"})
}

func (s *AdminService) RestoreCardBack(ctx context.Context, cardBackID string) (*models.CardBack, error) {
	cardBackID = strings.TrimSpace(cardBackID)
	if cardBackID == "" {
		return nil, wrapValidation("Card back id is required")
	}

	cardBack, err := s.cardBacks.FindByID(ctx, cardBackID)
	if err != nil {
		return nil, err
	}
	if cardBack.Status != "archived" {
		return cardBack, nil
	}

	return s.cardBacks.Update(ctx, cardBackID, bson.M{"status": "done"})
}

func (s *AdminService) DeleteCardBack(ctx context.Context, cardBackID string) error {
	cardBackID = strings.TrimSpace(cardBackID)
	if cardBackID == "" {
		return wrapValidation("Card back id is required")
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return err
	}
	if settings.DefaultBackID == cardBackID {
		return wrapValidation("Cannot delete the current default card back")
	}

	cardBack, err := s.cardBacks.FindByID(ctx, cardBackID)
	if err != nil {
		return err
	}
	if cardBack.Status != "archived" {
		return wrapValidation("Only archived card backs can be deleted permanently")
	}

	if err := s.cardBacks.Delete(ctx, cardBackID); err != nil {
		return err
	}

	return os.RemoveAll(filepath.Join("uploads", "card-backs", cardBackID))
}

func resolveAssetScopeStatuses(scope string) ([]string, error) {
	switch strings.TrimSpace(scope) {
	case "", "active":
		return []string{"done", "processing", "disabled"}, nil
	case "archived":
		return []string{"archived"}, nil
	case "all":
		return nil, nil
	default:
		return nil, wrapValidation("Unsupported asset scope")
	}
}

func sanitizeTags(raw []string) []string {
	result := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, tag := range raw {
		value := strings.ToLower(strings.TrimSpace(tag))
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
