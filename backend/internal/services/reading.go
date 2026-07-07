package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardOfDayAlreadyGeneratedError struct {
	Existing        models.ReadingResponse
	NextAvailableAt time.Time
}

func (e *CardOfDayAlreadyGeneratedError) Error() string {
	return ErrCardOfDayExists.Error()
}

type DailyLimitReachedError struct {
	NextResetAt time.Time
	Count       int64
	Limit       int64
}

func (e *DailyLimitReachedError) Error() string {
	return ErrDailyLimitReached.Error()
}

type ReadingHistoryResult struct {
	Readings      []models.ReadingResponse
	Pagination    models.PaginationResponse
	RetentionInfo models.RetentionInfoResponse
}

type ReadingService struct {
	users     *repositories.UserRepository
	cards     *repositories.CardRepository
	decks     *repositories.DeckRepository
	cardBacks *repositories.CardBackRepository
	readings  *repositories.ReadingRepository
	policies  *repositories.PolicyRepository
	settings  *repositories.AdminSettingsRepository
	jwtSecret string
	now       func() time.Time
}

func NewReadingService(
	users *repositories.UserRepository,
	cards *repositories.CardRepository,
	decks *repositories.DeckRepository,
	cardBacks *repositories.CardBackRepository,
	readings *repositories.ReadingRepository,
	policies *repositories.PolicyRepository,
	settings *repositories.AdminSettingsRepository,
	jwtSecret string,
) *ReadingService {
	return &ReadingService{
		users:     users,
		cards:     cards,
		decks:     decks,
		cardBacks: cardBacks,
		readings:  readings,
		policies:  policies,
		settings:  settings,
		jwtSecret: jwtSecret,
		now:       time.Now,
	}
}

func (s *ReadingService) GenerateGuestCardOfDay(ctx context.Context, req models.GenerateGuestCardOfDayRequest) (*models.GuestCardOfDayResult, error) {
	timezone := strings.TrimSpace(req.Timezone)
	now := s.now()
	accountingDate, err := utils.GetAccountingDateAt(now, timezone)
	if err != nil {
		return nil, wrapValidation("Invalid timezone")
	}
	nextResetAt, err := utils.GetNextResetTimeAt(now, timezone)
	if err != nil {
		return nil, wrapValidation("Invalid timezone")
	}

	keys, err := s.cards.ListAllKeys(ctx)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, wrapValidation("No tarot cards are available")
	}

	deckID, err := s.resolveEffectiveDeckID(ctx, "")
	if err != nil {
		return nil, err
	}
	backID, err := s.resolveEffectiveBackID(ctx, "")
	if err != nil {
		return nil, err
	}
	seed := buildDeterministicSeed("guest:"+req.GuestID, accountingDate, timezone)
	cardKey := selectDeterministicCardKey(keys, seed)
	reading := &models.Reading{
		Type:           models.ReadingTypeCardOfDay,
		AccountingDate: accountingDate,
		Timezone:       timezone,
		Cards:          []models.ReadingCard{{Key: cardKey, Position: 0, Reversed: false, Viewed: false}},
		DeckID:         deckID,
		BackID:         backID,
		GeneratedAt:    now.UTC(),
	}
	response, err := s.buildReadingResponse(ctx, reading, nil)
	if err != nil {
		return nil, err
	}
	response.ID = ""
	response.NextAvailableAt = &nextResetAt
	claimToken, err := utils.GenerateGuestCardToken(req.GuestID, cardKey, accountingDate, timezone, s.jwtSecret, now, nextResetAt)
	if err != nil {
		return nil, err
	}

	return &models.GuestCardOfDayResult{Reading: *response, ClaimToken: claimToken}, nil
}

func (s *ReadingService) ResolveCardOfDay(ctx context.Context, userID primitive.ObjectID, req models.ResolveCardOfDayRequest, requestIP string) (*models.ReadingResponse, string, error) {
	claimToken := strings.TrimSpace(req.ClaimToken)
	if claimToken == "" {
		reading, err := s.GenerateCardOfDay(ctx, userID, models.GenerateReadingRequest{Timezone: req.Timezone}, requestIP)
		if err != nil {
			var alreadyExists *CardOfDayAlreadyGeneratedError
			if errors.As(err, &alreadyExists) {
				return &alreadyExists.Existing, "existing", nil
			}
			return nil, "", err
		}
		return reading, "generated", nil
	}

	claims, err := utils.VerifyGuestCardToken(claimToken, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("%w: token verification failed", ErrInvalidGuestClaim)
	}
	user, policy, timezone, accountingDate, nextResetAt, err := s.resolveReadingContext(ctx, userID, claims.Timezone)
	if err != nil {
		return nil, "", err
	}
	if timezone != claims.Timezone || accountingDate != claims.AccountingDate {
		return nil, "", fmt.Errorf("%w: claim is outside its accounting day", ErrInvalidGuestClaim)
	}

	existing, err := s.readings.FindByUserTypeAndAccountingDate(ctx, userID, models.ReadingTypeCardOfDay, accountingDate)
	if err == nil {
		response, buildErr := s.buildReadingResponse(ctx, existing, nil)
		if buildErr != nil {
			return nil, "", buildErr
		}
		response.NextAvailableAt = &nextResetAt
		return response, "existing", nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, "", err
	}
	if _, err := s.cards.FindByKey(ctx, claims.CardKey); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, "", fmt.Errorf("%w: card no longer exists", ErrInvalidGuestClaim)
		}
		return nil, "", err
	}

	deckID, err := s.resolveEffectiveDeckID(ctx, user.SelectedDeckID)
	if err != nil {
		return nil, "", err
	}
	backID, err := s.resolveEffectiveBackID(ctx, user.SelectedBackID)
	if err != nil {
		return nil, "", err
	}
	now := s.now().UTC()
	reading := &models.Reading{
		UserID:         user.ID,
		Type:           models.ReadingTypeCardOfDay,
		AccountingDate: accountingDate,
		Timezone:       timezone,
		Cards:          []models.ReadingCard{{Key: claims.CardKey, Position: 0, Reversed: false, Viewed: false}},
		DeckID:         deckID,
		BackID:         backID,
		GeneratedAt:    now,
		ExpiresAt:      now.AddDate(0, 0, policy.HistoryRetentionDays),
		Metadata: models.ReadingMetadata{
			Seed:      buildDeterministicSeed("guest:"+claims.GuestID, accountingDate, timezone),
			RequestIP: requestIP,
		},
	}
	if err := s.readings.Create(ctx, reading); err != nil {
		return nil, "", err
	}
	response, err := s.buildReadingResponse(ctx, reading, nil)
	if err != nil {
		return nil, "", err
	}
	response.NextAvailableAt = &nextResetAt
	return response, "guest_claimed", nil
}

func (s *ReadingService) GenerateCardOfDay(ctx context.Context, userID primitive.ObjectID, req models.GenerateReadingRequest, requestIP string) (*models.ReadingResponse, error) {
	user, policy, timezone, accountingDate, nextResetAt, err := s.resolveReadingContext(ctx, userID, req.Timezone)
	if err != nil {
		return nil, err
	}
	deckID, err := s.resolveEffectiveDeckID(ctx, user.SelectedDeckID)
	if err != nil {
		return nil, err
	}
	backID, err := s.resolveEffectiveBackID(ctx, user.SelectedBackID)
	if err != nil {
		return nil, err
	}

	existing, err := s.readings.FindByUserTypeAndAccountingDate(ctx, userID, models.ReadingTypeCardOfDay, accountingDate)
	if err == nil {
		response, buildErr := s.buildReadingResponse(ctx, existing, nil)
		if buildErr != nil {
			return nil, buildErr
		}
		response.NextAvailableAt = &nextResetAt
		return nil, &CardOfDayAlreadyGeneratedError{
			Existing:        *response,
			NextAvailableAt: nextResetAt,
		}
	}
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	keys, err := s.cards.ListAllKeys(ctx)
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, wrapValidation("No tarot cards are available")
	}

	seed := buildDeterministicSeed(user.ID.Hex(), accountingDate, timezone)
	selectedKey := selectDeterministicCardKey(keys, seed)
	now := s.now().UTC()
	reading := &models.Reading{
		UserID:         user.ID,
		Type:           models.ReadingTypeCardOfDay,
		AccountingDate: accountingDate,
		Timezone:       timezone,
		Cards: []models.ReadingCard{
			{Key: selectedKey, Position: 0, Reversed: false, Viewed: false},
		},
		DeckID:      deckID,
		BackID:      backID,
		GeneratedAt: now,
		ExpiresAt:   now.AddDate(0, 0, policy.HistoryRetentionDays),
		Metadata: models.ReadingMetadata{
			Seed:      seed,
			RequestIP: requestIP,
		},
	}

	if err := s.readings.Create(ctx, reading); err != nil {
		return nil, err
	}

	response, err := s.buildReadingResponse(ctx, reading, nil)
	if err != nil {
		return nil, err
	}
	response.NextAvailableAt = &nextResetAt
	return response, nil
}

func (s *ReadingService) GenerateThreeCardSpread(ctx context.Context, userID primitive.ObjectID, req models.GenerateReadingRequest, requestIP string) (*models.ReadingResponse, error) {
	user, policy, timezone, accountingDate, nextResetAt, err := s.resolveReadingContext(ctx, userID, req.Timezone)
	if err != nil {
		return nil, err
	}
	deckID, err := s.resolveEffectiveDeckID(ctx, user.SelectedDeckID)
	if err != nil {
		return nil, err
	}
	backID, err := s.resolveEffectiveBackID(ctx, user.SelectedBackID)
	if err != nil {
		return nil, err
	}

	currentCount, err := s.readings.CountByUserTypeAndAccountingDate(ctx, userID, models.ReadingTypeThreeCard, accountingDate)
	if err != nil {
		return nil, err
	}

	if currentCount >= policy.ThreeCardLimit {
		return nil, &DailyLimitReachedError{
			NextResetAt: nextResetAt,
			Count:       currentCount,
			Limit:       policy.ThreeCardLimit,
		}
	}

	keys, err := s.cards.ListAllKeys(ctx)
	if err != nil {
		return nil, err
	}
	selectedKeys, err := selectRandomUniqueCardKeys(keys, 3)
	if err != nil {
		return nil, err
	}

	readingCards := make([]models.ReadingCard, 0, len(selectedKeys))
	for index, key := range selectedKeys {
		readingCards = append(readingCards, models.ReadingCard{
			Key:      key,
			Position: index,
			Reversed: false,
			Viewed:   false,
		})
	}

	now := s.now().UTC()
	reading := &models.Reading{
		UserID:         user.ID,
		Type:           models.ReadingTypeThreeCard,
		AccountingDate: accountingDate,
		Timezone:       timezone,
		Cards:          readingCards,
		DeckID:         deckID,
		BackID:         backID,
		GeneratedAt:    now,
		ExpiresAt:      now.AddDate(0, 0, policy.HistoryRetentionDays),
		Metadata: models.ReadingMetadata{
			RequestIP: requestIP,
		},
	}

	if err := s.readings.Create(ctx, reading); err != nil {
		return nil, err
	}

	response, err := s.buildReadingResponse(ctx, reading, &models.RemainingToday{
		Count: calculateRemainingCount(policy.ThreeCardLimit, currentCount),
		Limit: policy.ThreeCardLimit,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ReadingService) GetReadingHistory(ctx context.Context, userID primitive.ObjectID, readingType string, page, limit int) (*ReadingHistoryResult, error) {
	if readingType == "" {
		readingType = "all"
	}
	if readingType != "all" && readingType != models.ReadingTypeCardOfDay && readingType != models.ReadingTypeThreeCard {
		return nil, wrapValidation("Unsupported reading type filter")
	}
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	policy, err := s.policies.FindByRole(ctx, user.Role)
	if err != nil {
		return nil, err
	}

	total, err := s.readings.CountHistory(ctx, userID, readingType)
	if err != nil {
		return nil, err
	}

	skip := int64((page - 1) * limit)
	readings, err := s.readings.ListHistory(ctx, userID, readingType, skip, int64(limit))
	if err != nil {
		return nil, err
	}

	responses, err := s.buildReadingResponses(ctx, readings, nil)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	oldestDate, err := buildOldestRetentionDate(s.now(), user.Timezone, policy.HistoryRetentionDays)
	if err != nil {
		return nil, err
	}

	return &ReadingHistoryResult{
		Readings: responses,
		Pagination: models.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		RetentionInfo: models.RetentionInfoResponse{
			Role:          user.Role,
			RetentionDays: policy.HistoryRetentionDays,
			OldestDate:    oldestDate,
		},
	}, nil
}

func (s *ReadingService) RevealReadingCard(ctx context.Context, userID primitive.ObjectID, readingID string, position int) (*models.ReadingResponse, error) {
	readingObjectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(readingID))
	if err != nil {
		return nil, wrapValidation("Invalid reading id")
	}
	if position < 0 {
		return nil, wrapValidation("Card position must be zero or greater")
	}

	reading, err := s.readings.FindByIDAndUser(ctx, readingObjectID, userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, wrapValidation("Reading not found")
		}
		return nil, err
	}

	cardFound := false
	updatedCards := make([]models.ReadingCard, len(reading.Cards))
	copy(updatedCards, reading.Cards)

	for index := range updatedCards {
		if updatedCards[index].Position != position {
			continue
		}

		cardFound = true
		if updatedCards[index].Viewed {
			return s.buildReadingResponse(ctx, reading, nil)
		}

		updatedCards[index].Viewed = true
		break
	}

	if !cardFound {
		return nil, wrapValidation("Card position not found in reading")
	}

	if err := s.readings.UpdateCards(ctx, reading.ID, userID, updatedCards); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, wrapValidation("Reading not found")
		}
		return nil, err
	}

	reading.Cards = updatedCards
	return s.buildReadingResponse(ctx, reading, nil)
}

func (s *ReadingService) resolveReadingContext(ctx context.Context, userID primitive.ObjectID, requestedTimezone string) (*models.User, *models.RolePolicy, string, string, time.Time, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, "", "", time.Time{}, err
	}
	policy, err := s.policies.FindByRole(ctx, user.Role)
	if err != nil {
		return nil, nil, "", "", time.Time{}, err
	}

	timezone := resolveTimezone(requestedTimezone, user.Timezone)
	accountingDate, err := utils.GetAccountingDateAt(s.now(), timezone)
	if err != nil {
		return nil, nil, "", "", time.Time{}, err
	}
	nextResetAt, err := utils.GetNextResetTimeAt(s.now(), timezone)
	if err != nil {
		return nil, nil, "", "", time.Time{}, err
	}

	return user, policy, timezone, accountingDate, nextResetAt, nil
}

func (s *ReadingService) buildReadingResponse(ctx context.Context, reading *models.Reading, remaining *models.RemainingToday) (*models.ReadingResponse, error) {
	responses, err := s.buildReadingResponses(ctx, []models.Reading{*reading}, remaining)
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return &responses[0], nil
}

func (s *ReadingService) buildReadingResponses(ctx context.Context, readings []models.Reading, remaining *models.RemainingToday) ([]models.ReadingResponse, error) {
	if len(readings) == 0 {
		return []models.ReadingResponse{}, nil
	}

	keys := uniqueCardKeys(readings)
	cards, err := s.cards.FindByKeys(ctx, keys)
	if err != nil {
		return nil, err
	}

	cardMap := make(map[string]models.Card, len(cards))
	for _, card := range cards {
		cardMap[card.Key] = card
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return nil, err
	}
	defaultBackID := resolveBackID(settings.DefaultBackID)
	backImageURLs, err := s.buildBackImageURLMap(ctx, readings, defaultBackID)
	if err != nil {
		return nil, err
	}

	responses := make([]models.ReadingResponse, 0, len(readings))
	for _, reading := range readings {
		effectiveBackID := resolveReadingBackID(reading.BackID, defaultBackID)
		responseCards := make([]models.ReadingCardResponse, 0, len(reading.Cards))
		for _, readingCard := range reading.Cards {
			card, ok := cardMap[readingCard.Key]
			if !ok {
				return nil, fmt.Errorf("card %s not found in catalog", readingCard.Key)
			}

			responseCards = append(responseCards, buildReadingCardResponse(
				card,
				readingCard,
				buildCardImageURL(reading.DeckID, card.Key),
				backImageURLs[effectiveBackID],
			))
		}

		sort.Slice(responseCards, func(i, j int) bool {
			return responseCards[i].Position < responseCards[j].Position
		})

		response := models.ReadingResponse{
			ID:             reading.ID.Hex(),
			Type:           reading.Type,
			AccountingDate: reading.AccountingDate,
			GeneratedAt:    reading.GeneratedAt,
			Cards:          responseCards,
			DeckID:         reading.DeckID,
			BackID:         effectiveBackID,
			AI:             reading.AI.ToResponse(),
		}
		if remaining != nil && reading.Type == models.ReadingTypeThreeCard {
			copyRemaining := *remaining
			response.RemainingToday = &copyRemaining
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func buildReadingCardResponse(card models.Card, readingCard models.ReadingCard, imageURL, backImageURL string) models.ReadingCardResponse {
	return models.ReadingCardResponse{
		Key:              card.Key,
		Arcana:           card.Arcana,
		Type:             card.Type.EN,
		TypeRU:           card.Type.RU,
		Name:             card.Names.EN,
		NameRU:           card.Names.RU,
		Interpretation:   card.Interpretations.EN,
		InterpretationRU: card.Interpretations.RU,
		Tags:             card.Tags.EN,
		TagsRU:           card.Tags.RU,
		ImageURL:         imageURL,
		BackImageURL:     backImageURL,
		Position:         readingCard.Position,
		Reversed:         readingCard.Reversed,
		Viewed:           readingCard.Viewed,
	}
}

func resolveTimezone(requestedTimezone, fallback string) string {
	timezone := strings.TrimSpace(requestedTimezone)
	if timezone == "" {
		timezone = strings.TrimSpace(fallback)
	}
	if timezone == "" {
		return "UTC"
	}
	return timezone
}

func resolveDeckID(deckID string) string {
	deckID = strings.TrimSpace(deckID)
	if deckID == "" {
		return "default"
	}
	return deckID
}

func resolveBackID(backID string) string {
	backID = strings.TrimSpace(backID)
	if backID == "" {
		return "default"
	}
	return backID
}

func (s *ReadingService) resolveEffectiveDeckID(ctx context.Context, selectedDeckID string) (string, error) {
	selectedDeckID = strings.TrimSpace(selectedDeckID)
	if selectedDeckID != "" {
		deck, err := s.decksFallbackLookup(ctx, selectedDeckID)
		if err != nil {
			return "", err
		}
		if deck != nil {
			return selectedDeckID, nil
		}
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return "", err
	}

	return resolveDeckID(settings.DefaultDeckID), nil
}

func (s *ReadingService) resolveEffectiveBackID(ctx context.Context, selectedBackID string) (string, error) {
	selectedBackID = strings.TrimSpace(selectedBackID)
	if selectedBackID != "" {
		cardBack, err := s.cardBackFallbackLookup(ctx, selectedBackID)
		if err != nil {
			return "", err
		}
		if cardBack != nil {
			return selectedBackID, nil
		}
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return "", err
	}

	return resolveBackID(settings.DefaultBackID), nil
}

func (s *ReadingService) decksFallbackLookup(ctx context.Context, deckID string) (*models.Deck, error) {
	deck, err := s.decks.FindByID(ctx, deckID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if deck.Status != "done" {
		return nil, nil
	}
	return deck, nil
}

func (s *ReadingService) cardBackFallbackLookup(ctx context.Context, cardBackID string) (*models.CardBack, error) {
	cardBack, err := s.cardBacks.FindByID(ctx, cardBackID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if cardBack.Status != "done" {
		return nil, nil
	}
	return cardBack, nil
}

func buildDeterministicSeed(userID, accountingDate, timezone string) string {
	sum := sha256.Sum256([]byte(userID + "|" + accountingDate + "|" + timezone))
	return fmt.Sprintf("%x", sum[:])
}

func selectDeterministicCardKey(keys []string, seed string) string {
	if len(keys) == 0 {
		return ""
	}
	if len(keys) == 1 {
		return keys[0]
	}

	hash := sha256.Sum256([]byte(seed))
	index := binary.BigEndian.Uint64(hash[:8]) % uint64(len(keys))
	return keys[index]
}

func selectRandomUniqueCardKeys(keys []string, count int) ([]string, error) {
	if count <= 0 {
		return []string{}, nil
	}
	if len(keys) < count {
		return nil, wrapValidation("Not enough tarot cards are available")
	}

	available := append([]string(nil), keys...)
	selected := make([]string, 0, count)

	for i := 0; i < count; i++ {
		max := big.NewInt(int64(len(available)))
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return nil, err
		}
		index := int(n.Int64())
		selected = append(selected, available[index])
		available = append(available[:index], available[index+1:]...)
	}

	return selected, nil
}

func uniqueCardKeys(readings []models.Reading) []string {
	set := make(map[string]struct{})
	for _, reading := range readings {
		for _, card := range reading.Cards {
			set[card.Key] = struct{}{}
		}
	}

	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func uniqueBackIDs(readings []models.Reading, defaultBackID string) []string {
	set := make(map[string]struct{})
	for _, reading := range readings {
		backID := resolveReadingBackID(reading.BackID, defaultBackID)
		set[backID] = struct{}{}
	}

	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func resolveReadingBackID(backID, defaultBackID string) string {
	backID = strings.TrimSpace(backID)
	if backID != "" {
		return backID
	}
	return resolveBackID(defaultBackID)
}

func buildCardImageURL(deckID, cardKey string) string {
	deckID = resolveDeckID(deckID)
	relativePath := fmt.Sprintf("/uploads/decks/%s/cards/%s.png", deckID, cardKey)
	absolutePath := filepath.Join("uploads", "decks", deckID, "cards", cardKey+".png")

	if _, err := os.Stat(absolutePath); err != nil {
		return ""
	}

	return relativePath
}

func (s *ReadingService) buildBackImageURLMap(ctx context.Context, readings []models.Reading, defaultBackID string) (map[string]string, error) {
	backIDs := uniqueBackIDs(readings, defaultBackID)
	backImageURLs := make(map[string]string, len(backIDs))

	for _, backID := range backIDs {
		cardBack, err := s.cardBacks.FindByID(ctx, backID)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				backImageURLs[backID] = ""
				continue
			}
			return nil, err
		}

		backImageURLs[backID] = buildStoredAssetURL(cardBack.ImagePath)
	}

	return backImageURLs, nil
}

func buildStoredAssetURL(relativePath string) string {
	relativePath = strings.TrimSpace(relativePath)
	if relativePath == "" {
		return ""
	}

	absolutePath := filepath.Clean(strings.TrimPrefix(relativePath, "/"))
	if _, err := os.Stat(absolutePath); err != nil {
		return ""
	}

	return relativePath
}

func calculateRemainingCount(limit, currentCount int64) int64 {
	remaining := limit - (currentCount + 1)
	if remaining < 0 {
		return 0
	}
	return remaining
}

func buildOldestRetentionDate(now time.Time, timezone string, retentionDays int) (string, error) {
	loc, err := time.LoadLocation(resolveTimezone(timezone, "UTC"))
	if err != nil {
		return "", err
	}

	return now.In(loc).AddDate(0, 0, -retentionDays).Format("2006-01-02"), nil
}
