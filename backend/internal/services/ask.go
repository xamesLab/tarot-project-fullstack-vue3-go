package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"tarot-daily/config"
	"tarot-daily/internal/models"
	"tarot-daily/internal/repositories"
	"tarot-daily/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const askPromptVersion = "v0.1.1-compact-100w"

type AIProvider interface {
	GenerateInterpretation(ctx context.Context, req AIProviderRequest) (*AIProviderResult, error)
}

type AIProviderRequest struct {
	Model     string
	Language  string
	Messages  []AIProviderMessage
	MaxTokens int
}

type AIProviderMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIProviderResult struct {
	ProviderRequestID string
	Content           string
	FinishReason      string
	Usage             *models.AskUsage
}

type AIProviderError struct {
	Code       string
	Message    string
	StatusCode int
	Retryable  bool
}

func (e *AIProviderError) Error() string {
	if e == nil {
		return ""
	}
	return e.Code
}

type AskService struct {
	config   *config.Config
	users    *repositories.UserRepository
	readings *repositories.ReadingRepository
	cards    *repositories.CardRepository
	requests *repositories.AskRequestRepository
	policies *repositories.PolicyRepository
	settings *repositories.AdminSettingsRepository
	provider AIProvider
	now      func() time.Time
}

func NewAskService(
	cfg *config.Config,
	users *repositories.UserRepository,
	readings *repositories.ReadingRepository,
	cards *repositories.CardRepository,
	requests *repositories.AskRequestRepository,
	policies *repositories.PolicyRepository,
	settings *repositories.AdminSettingsRepository,
) *AskService {
	return &AskService{
		config:   cfg,
		users:    users,
		readings: readings,
		cards:    cards,
		requests: requests,
		policies: policies,
		settings: settings,
		provider: NewOpenRouterProvider(cfg),
		now:      func() time.Time { return time.Now().UTC() },
	}
}

type AIDailyLimitReachedError struct {
	NextResetAt time.Time
	Count       int64
	Limit       int64
}

func (e *AIDailyLimitReachedError) Error() string {
	return ErrAIDailyLimitReached.Error()
}

func (s *AskService) CreateOrReuseInterpretation(ctx context.Context, userID primitive.ObjectID, req models.CreateAskRequest) (*models.AskRequestResponse, error) {
	if !s.config.AIDirectEnabled {
		return nil, ErrAIUnavailable
	}
	if strings.TrimSpace(s.config.OpenRouterKey) == "" || strings.TrimSpace(s.config.OpenRouterModelMain) == "" {
		return nil, ErrAIUnavailable
	}

	readingID, err := primitive.ObjectIDFromHex(strings.TrimSpace(req.ReadingID))
	if err != nil {
		return nil, wrapValidation("Invalid reading id")
	}

	reading, err := s.readings.FindByIDAndUser(ctx, readingID, userID)
	if err != nil {
		return nil, err
	}

	if reading.AI != nil && reading.AI.Status == models.AskRequestStatusCompleted && strings.TrimSpace(reading.AI.Content) != "" {
		existing, err := s.findExistingRequest(ctx, userID, reading)
		if err == nil && shouldReuseAskRequest(existing) {
			response := existing.ToResponse()
			return &response, nil
		}
	}

	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureAIAvailableForUser(ctx, user); err != nil {
		return nil, err
	}

	payload, providerRequest, err := s.buildAIRequest(ctx, reading, user)
	if err != nil {
		return nil, err
	}

	now := s.now()
	askRequest := &models.AskRequest{
		UserID:         userID,
		ReadingID:      reading.ID,
		Provider:       models.AskProviderOpenRouter,
		Model:          s.config.OpenRouterModelMain,
		PromptVersion:  askPromptVersion,
		Status:         models.AskRequestStatusPending,
		RequestPayload: payload,
		RequestedAt:    now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := s.requests.Create(ctx, askRequest); err != nil {
		return nil, err
	}

	processingAt := now
	if err := s.requests.UpdateStatus(
		ctx,
		askRequest.ID,
		models.AskRequestStatusProcessing,
		nil,
		"",
		"",
		"",
		nil,
		nil,
	); err != nil {
		return nil, err
	}
	if err := s.readings.UpdateAIInterpretation(ctx, reading.ID, userID, &models.AIInterpretation{
		RequestID:       askRequest.ID,
		Provider:        models.AskProviderOpenRouter,
		Model:           s.config.OpenRouterModelMain,
		PromptVersion:   askPromptVersion,
		Language:        payload.Language,
		Status:          models.AskRequestStatusProcessing,
		LastRequestedAt: &processingAt,
		UpdatedAt:       &processingAt,
	}); err != nil {
		return nil, err
	}

	result, err := s.generateInterpretationWithRetry(ctx, providerRequest)
	if err != nil {
		return nil, s.handleProviderFailure(ctx, userID, reading.ID, askRequest.ID, payload.Language, err, now)
	}

	content := normalizeInterpretation(result.Content)
	if content == "" {
		providerErr := &AIProviderError{
			Code:       "ai_provider_invalid_response",
			Message:    "AI provider returned an empty interpretation",
			StatusCode: http.StatusBadGateway,
		}
		return nil, s.handleProviderFailure(ctx, userID, reading.ID, askRequest.ID, payload.Language, providerErr, now)
	}

	completedAt := s.now()
	if err := s.requests.UpdateStatus(
		ctx,
		askRequest.ID,
		models.AskRequestStatusCompleted,
		&completedAt,
		content,
		"",
		"",
		result.Usage,
		map[string]any{
			"provider_request_id": result.ProviderRequestID,
			"finish_reason":       result.FinishReason,
		},
	); err != nil {
		return nil, err
	}

	if err := s.readings.UpdateAIInterpretation(ctx, reading.ID, userID, &models.AIInterpretation{
		RequestID:       askRequest.ID,
		Provider:        models.AskProviderOpenRouter,
		Model:           s.config.OpenRouterModelMain,
		PromptVersion:   askPromptVersion,
		Language:        payload.Language,
		Status:          models.AskRequestStatusCompleted,
		Content:         content,
		CompletedAt:     &completedAt,
		LastRequestedAt: &now,
		UpdatedAt:       &completedAt,
	}); err != nil {
		return nil, err
	}

	askRequest.Status = models.AskRequestStatusCompleted
	askRequest.ResponseExcerpt = content
	askRequest.Usage = result.Usage
	askRequest.ResponseData = map[string]any{
		"provider_request_id": result.ProviderRequestID,
		"finish_reason":       result.FinishReason,
	}
	askRequest.CompletedAt = &completedAt

	response := askRequest.ToResponse()
	return &response, nil
}

func (s *AskService) generateInterpretationWithRetry(ctx context.Context, req AIProviderRequest) (*AIProviderResult, error) {
	result, err := s.provider.GenerateInterpretation(ctx, req)
	if err != nil {
		return nil, err
	}
	if !isIncompleteFinishReason(result.FinishReason) {
		return result, nil
	}

	retryRequest := req
	if retryRequest.MaxTokens < 900 {
		retryRequest.MaxTokens += 320
		if retryRequest.MaxTokens > 900 {
			retryRequest.MaxTokens = 900
		}
	}

	retryResult, retryErr := s.provider.GenerateInterpretation(ctx, retryRequest)
	if retryErr != nil {
		return nil, retryErr
	}
	if isIncompleteFinishReason(retryResult.FinishReason) {
		return nil, &AIProviderError{
			Code:       "ai_provider_incomplete_response",
			Message:    "AI provider returned an incomplete response",
			StatusCode: http.StatusBadGateway,
			Retryable:  true,
		}
	}

	return retryResult, nil
}

func (s *AskService) ensureAIAvailableForUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return mongo.ErrNoDocuments
	}

	settings, err := s.settings.Get(ctx)
	if err != nil {
		return err
	}
	if !settings.AIEnabled {
		return ErrAIUnavailable
	}

	policy, err := s.policies.FindByRole(ctx, user.Role)
	if err != nil {
		return err
	}
	if !policy.AIEnabled {
		return ErrAIUnavailable
	}
	if policy.AIDailyLimit <= 0 {
		return &AIDailyLimitReachedError{Limit: policy.AIDailyLimit}
	}

	timezone := strings.TrimSpace(user.Timezone)
	if timezone == "" {
		timezone = "UTC"
	}

	windowStart, nextResetAt, err := utils.GetCurrentResetWindowAt(s.now(), timezone)
	if err != nil {
		return err
	}

	currentCount, err := s.requests.CountByUserSince(ctx, user.ID, windowStart)
	if err != nil {
		return err
	}
	if currentCount >= policy.AIDailyLimit {
		return &AIDailyLimitReachedError{
			NextResetAt: nextResetAt,
			Count:       currentCount,
			Limit:       policy.AIDailyLimit,
		}
	}

	return nil
}

func (s *AskService) GetRequest(ctx context.Context, userID primitive.ObjectID, requestID string) (*models.AskRequestResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(requestID))
	if err != nil {
		return nil, wrapValidation("Invalid request id")
	}

	askRequest, err := s.requests.FindByIDAndUser(ctx, objectID, userID)
	if err != nil {
		return nil, err
	}

	response := askRequest.ToResponse()
	return &response, nil
}

func (s *AskService) findExistingRequest(ctx context.Context, userID primitive.ObjectID, reading *models.Reading) (*models.AskRequest, error) {
	if reading == nil || reading.AI == nil {
		return nil, mongo.ErrNoDocuments
	}
	if !reading.AI.RequestID.IsZero() {
		return s.requests.FindByIDAndUser(ctx, reading.AI.RequestID, userID)
	}
	return s.requests.FindLatestByReadingAndUser(ctx, reading.ID, userID)
}

func (s *AskService) buildAIRequest(ctx context.Context, reading *models.Reading, user *models.User) (models.AskRequestPayload, AIProviderRequest, error) {
	if reading == nil {
		return models.AskRequestPayload{}, AIProviderRequest{}, mongo.ErrNoDocuments
	}
	if user == nil {
		return models.AskRequestPayload{}, AIProviderRequest{}, mongo.ErrNoDocuments
	}

	readyCards, err := ensureReadingReadyForAI(reading)
	if err != nil {
		return models.AskRequestPayload{}, AIProviderRequest{}, err
	}

	keys := make([]string, 0, len(readyCards))
	for _, card := range readyCards {
		keys = append(keys, card.Key)
	}

	cardRecords, err := s.cards.FindByKeys(ctx, keys)
	if err != nil {
		return models.AskRequestPayload{}, AIProviderRequest{}, err
	}

	cardMap := make(map[string]models.Card, len(cardRecords))
	for _, card := range cardRecords {
		cardMap[card.Key] = card
	}

	language := normalizeLanguage(user.Language)
	payload := models.AskRequestPayload{
		ReadingType: reading.Type,
		Language:    language,
		Cards:       make([]models.AskRequestCardInput, 0, len(readyCards)),
		MaxWords:    200,
	}
	lines := make([]string, 0, len(readyCards)+4)
	lines = append(lines, fmt.Sprintf("Reading type: %s", reading.Type))
	lines = append(lines, fmt.Sprintf("Response language: %s", language))
	lines = append(lines, "Cards:")

	for _, readingCard := range readyCards {
		card, ok := cardMap[readingCard.Key]
		if !ok {
			return models.AskRequestPayload{}, AIProviderRequest{}, fmt.Errorf("card %s not found in catalog", readingCard.Key)
		}

		payload.Cards = append(payload.Cards, models.AskRequestCardInput{
			Key:      readingCard.Key,
			Position: readingCard.Position,
		})

		name := localizedCardName(card, language)
		positionLabel := localizedPosition(reading.Type, readingCard.Position, language)
		if positionLabel != "" {
			lines = append(lines, fmt.Sprintf("- %s: %s", positionLabel, name))
			continue
		}
		lines = append(lines, fmt.Sprintf("- %s", name))
	}

	systemPrompt := "You are a tarot interpreter. The user's system language is English. Respond only in English. Return only one compact interpretation of the revealed reading. Do not use headings, bullet lists, disclaimers, or meta commentary. Aim for about 200 words. This is a recommendation, not a hard limit."
	if language == "ru" {
		systemPrompt = "Ты таро-интерпретатор. Системный язык пользователя — русский. Отвечай только на русском языке. Верни только одну компактную интерпретацию раскрытого расклада. Не используй заголовки, списки, дисклеймеры и мета-комментарии. Ориентируйся на объём около 200 слов. Это рекомендация, а не жёсткий лимит."
	}

	userPrompt := strings.Join(lines, "\n") + "\nUse the user's system language exactly. Return only the compact interpretation. Aim for about 200 words, but treat this as a recommendation rather than a strict limit."
	if language == "ru" {
		userPrompt = strings.Join(lines, "\n") + "\nИспользуй ровно тот язык, который установлен у пользователя в системе. Верни только компактную интерпретацию расклада. Ориентируйся на объём около 200 слов, но считай это рекомендацией, а не строгим ограничением."
	}

	return payload, AIProviderRequest{
		Model:     s.config.OpenRouterModelMain,
		Language:  language,
		MaxTokens: recommendedMaxTokens(reading.Type),
		Messages: []AIProviderMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}, nil
}

func (s *AskService) handleProviderFailure(ctx context.Context, userID, readingID, requestID primitive.ObjectID, language string, sourceErr error, requestedAt time.Time) error {
	var providerErr *AIProviderError
	if !errors.As(sourceErr, &providerErr) {
		providerErr = &AIProviderError{
			Code:       "ai_provider_failed",
			Message:    "AI provider request failed",
			StatusCode: http.StatusBadGateway,
		}
	}

	completedAt := s.now()
	if err := s.requests.UpdateStatus(
		ctx,
		requestID,
		models.AskRequestStatusFailed,
		&completedAt,
		"",
		providerErr.Code,
		providerErr.Message,
		nil,
		nil,
	); err != nil {
		return err
	}

	if err := s.readings.UpdateAIInterpretation(ctx, readingID, userID, &models.AIInterpretation{
		RequestID:       requestID,
		Provider:        models.AskProviderOpenRouter,
		Model:           s.config.OpenRouterModelMain,
		PromptVersion:   askPromptVersion,
		Language:        language,
		Status:          models.AskRequestStatusFailed,
		ErrorCode:       providerErr.Code,
		ErrorMessage:    providerErr.Message,
		CompletedAt:     &completedAt,
		LastRequestedAt: &requestedAt,
		UpdatedAt:       &completedAt,
	}); err != nil {
		return err
	}

	return providerErr
}

type OpenRouterProvider struct {
	baseURL    string
	apiKey     string
	referer    string
	title      string
	httpClient *http.Client
}

func NewOpenRouterProvider(cfg *config.Config) *OpenRouterProvider {
	return &OpenRouterProvider{
		baseURL: strings.TrimRight(cfg.OpenRouterBaseURL, "/"),
		apiKey:  cfg.OpenRouterKey,
		referer: strings.TrimSpace(cfg.OpenRouterReferer),
		title:   strings.TrimSpace(cfg.OpenRouterTitle),
		httpClient: &http.Client{
			Timeout: cfg.OpenRouterTimeout,
		},
	}
}

func (p *OpenRouterProvider) GenerateInterpretation(ctx context.Context, req AIProviderRequest) (*AIProviderResult, error) {
	payload := openRouterChatCompletionRequest{
		Model:     req.Model,
		Messages:  req.Messages,
		MaxTokens: req.MaxTokens,
	}
	if payload.MaxTokens <= 0 {
		payload.MaxTokens = 512
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	if p.referer != "" {
		httpReq.Header.Set("HTTP-Referer", p.referer)
	}
	if p.title != "" {
		httpReq.Header.Set("X-Title", p.title)
		httpReq.Header.Set("X-OpenRouter-Title", p.title)
	}

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, classifyOpenRouterTransportError(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, &AIProviderError{
			Code:       "ai_provider_read_failed",
			Message:    "Failed to read AI provider response",
			StatusCode: http.StatusBadGateway,
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, classifyOpenRouterHTTPError(resp.StatusCode, responseBody)
	}

	var completion openRouterChatCompletionResponse
	if err := json.Unmarshal(responseBody, &completion); err != nil {
		return nil, &AIProviderError{
			Code:       "ai_provider_invalid_response",
			Message:    "AI provider returned an invalid response",
			StatusCode: http.StatusBadGateway,
		}
	}

	if len(completion.Choices) == 0 || strings.TrimSpace(completion.Choices[0].Message.Content) == "" {
		return nil, &AIProviderError{
			Code:       "ai_provider_empty_response",
			Message:    "AI provider returned an empty response",
			StatusCode: http.StatusBadGateway,
		}
	}

	return &AIProviderResult{
		ProviderRequestID: completion.ID,
		Content:           completion.Choices[0].Message.Content,
		FinishReason:      completion.Choices[0].FinishReason,
		Usage: &models.AskUsage{
			PromptTokens:     completion.Usage.PromptTokens,
			CompletionTokens: completion.Usage.CompletionTokens,
			TotalTokens:      completion.Usage.TotalTokens,
		},
	}, nil
}

type openRouterChatCompletionRequest struct {
	Model     string              `json:"model"`
	Messages  []AIProviderMessage `json:"messages"`
	MaxTokens int                 `json:"max_tokens,omitempty"`
}

type openRouterChatCompletionResponse struct {
	ID      string `json:"id"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func classifyOpenRouterTransportError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return &AIProviderError{
			Code:       "ai_provider_timeout",
			Message:    "AI provider timed out",
			StatusCode: http.StatusGatewayTimeout,
			Retryable:  true,
		}
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &AIProviderError{
			Code:       "ai_provider_timeout",
			Message:    "AI provider timed out",
			StatusCode: http.StatusGatewayTimeout,
			Retryable:  true,
		}
	}

	return &AIProviderError{
		Code:       "ai_provider_network_error",
		Message:    "Unable to reach AI provider",
		StatusCode: http.StatusBadGateway,
		Retryable:  true,
	}
}

func classifyOpenRouterHTTPError(statusCode int, _ []byte) error {
	message := "AI provider request failed"
	switch statusCode {
	case http.StatusUnauthorized:
		message = "AI provider authentication failed"
	case http.StatusTooManyRequests:
		message = "AI provider rate limit reached"
	case http.StatusBadRequest:
		message = "AI provider rejected the request"
	default:
		if statusCode >= 500 {
			message = "AI provider is temporarily unavailable"
		}
	}

	return &AIProviderError{
		Code:       fmt.Sprintf("ai_provider_http_%d", statusCode),
		Message:    message,
		StatusCode: mapProviderStatusCode(statusCode),
		Retryable:  statusCode >= 500 || statusCode == http.StatusTooManyRequests,
	}
}

func mapProviderStatusCode(providerStatus int) int {
	switch providerStatus {
	case http.StatusUnauthorized:
		return http.StatusServiceUnavailable
	case http.StatusTooManyRequests:
		return http.StatusTooManyRequests
	default:
		if providerStatus >= 500 {
			return http.StatusServiceUnavailable
		}
		return http.StatusBadGateway
	}
}

func ensureReadingReadyForAI(reading *models.Reading) ([]models.ReadingCard, error) {
	if reading == nil {
		return nil, mongo.ErrNoDocuments
	}

	cards := append([]models.ReadingCard(nil), reading.Cards...)
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Position < cards[j].Position
	})

	switch reading.Type {
	case models.ReadingTypeCardOfDay:
		if len(cards) == 0 || !cards[0].Viewed {
			return nil, wrapValidation("Reveal the card before requesting AI interpretation")
		}
		return cards[:1], nil
	case models.ReadingTypeThreeCard:
		if len(cards) == 0 {
			return nil, wrapValidation("Reading has no cards")
		}
		for _, card := range cards {
			if !card.Viewed {
				return nil, wrapValidation("Reveal all cards before requesting AI interpretation")
			}
		}
		return cards, nil
	default:
		return nil, wrapValidation("AI interpretation is not supported for this reading type")
	}
}

func normalizeLanguage(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), "ru") {
		return "ru"
	}
	return "en"
}

func localizedCardName(card models.Card, language string) string {
	if language == "ru" && strings.TrimSpace(card.Names.RU) != "" {
		return card.Names.RU
	}
	return card.Names.EN
}

func localizedPosition(readingType string, position int, language string) string {
	if readingType != models.ReadingTypeThreeCard {
		return ""
	}

	if language == "ru" {
		switch position {
		case 0:
			return "Прошлое"
		case 1:
			return "Настоящее"
		case 2:
			return "Будущее"
		}
	} else {
		switch position {
		case 0:
			return "Past"
		case 1:
			return "Present"
		case 2:
			return "Future"
		}
	}

	return ""
}

func normalizeInterpretation(value string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}

func recommendedMaxTokens(readingType string) int {
	switch readingType {
	case models.ReadingTypeThreeCard:
		return 520
	case models.ReadingTypeCardOfDay:
		return 320
	default:
		return 420
	}
}

func isIncompleteFinishReason(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "length")
}

func shouldReuseAskRequest(request *models.AskRequest) bool {
	if request == nil || strings.TrimSpace(request.ResponseExcerpt) == "" {
		return false
	}
	if request.ResponseData == nil {
		return true
	}

	rawFinishReason, ok := request.ResponseData["finish_reason"]
	if !ok {
		return true
	}

	finishReason, ok := rawFinishReason.(string)
	if !ok {
		return true
	}

	return !isIncompleteFinishReason(finishReason)
}
