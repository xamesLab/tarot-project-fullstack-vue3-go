package ask

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/gatewayclient"
	"tarotdaily/backend/internal/platform/inpututil"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidConfig      = errors.New("invalid config")
	ErrRequestNotFound    = errors.New("request not found")
	ErrGatewayTimeout     = errors.New("gateway timeout")
	ErrGatewayRejected    = errors.New("gateway rejected request")
	ErrGatewayBadResponse = errors.New("gateway bad response")
	ErrGatewayUnavailable = errors.New("gateway unavailable")
)

type RequestView struct {
	RequestID    string `json:"request_id"`
	Status       string `json:"status"`
	ResultText   string `json:"result_text,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type userReader interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error)
}

type askRequestStore interface {
	Create(ctx context.Context, item *models.AskRequest) error
	FindByRequestID(ctx context.Context, requestID string) (*models.AskRequest, error)
	UpdateStatus(ctx context.Context, requestID string, status models.AskRequestStatus) (*models.AskRequest, error)
	UpdateResult(ctx context.Context, requestID string, resultText string) (*models.AskRequest, error)
	UpdateError(ctx context.Context, requestID string, errorCode string, errorMessage string) (*models.AskRequest, error)
}

type gatewayRequester interface {
	Request(ctx context.Context, payload gatewayclient.RequestPayload) (*gatewayclient.RequestResponse, error)
}

type Service struct {
	users                userReader
	requests             askRequestStore
	gateway              gatewayRequester
	backendPublicBaseURL string
}

func NewService(
	users userReader,
	requests askRequestStore,
	gateway gatewayRequester,
	backendPublicBaseURL string,
) *Service {
	return &Service{
		users:                users,
		requests:             requests,
		gateway:              gateway,
		backendPublicBaseURL: strings.TrimSpace(backendPublicBaseURL),
	}
}

func (s *Service) Create(ctx context.Context, userID string, text string) (*RequestView, error) {
	if s == nil || s.users == nil || s.requests == nil || s.gateway == nil {
		return nil, ErrInvalidConfig
	}

	userOID, err := inpututil.ParseObjectID(userID)
	if err != nil {
		return nil, ErrInvalidInput
	}
	if _, err := s.users.FindByID(ctx, userOID); err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrInvalidInput
		}
		return nil, err
	}

	normalizedText, err := inpututil.RequireTrimmed(text)
	if err != nil {
		return nil, ErrInvalidInput
	}

	callbackURL, err := s.callbackURL()
	if err != nil {
		return nil, err
	}

	requestID := buildRequestID()
	item := &models.AskRequest{
		RequestID: requestID,
		UserID:    userOID,
		Text:      normalizedText,
		Status:    models.AskRequestStatusAccepted,
	}

	if err := s.requests.Create(ctx, item); err != nil {
		return nil, err
	}

	resp, err := s.gateway.Request(ctx, gatewayclient.RequestPayload{
		Text:        normalizedText,
		CallbackURL: callbackURL,
		RequestID:   requestID,
		Metadata: gatewayclient.RequestMetadata{
			UserID: userOID.Hex(),
		},
	})
	if err != nil {
		errorCode, mappedErr := mapGatewayDispatchError(err)
		_, _ = s.requests.UpdateError(ctx, requestID, errorCode, err.Error())
		return nil, mappedErr
	}

	status := normalizeStatus(resp.Status)
	updated, err := s.requests.UpdateStatus(ctx, requestID, status)
	if err != nil {
		return nil, err
	}

	view := toRequestView(updated)
	return &view, nil
}

func (s *Service) Get(ctx context.Context, userID, requestID string) (*RequestView, error) {
	if s == nil || s.requests == nil {
		return nil, ErrInvalidConfig
	}

	userOID, err := inpututil.ParseObjectID(userID)
	if err != nil {
		return nil, ErrInvalidInput
	}
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, ErrInvalidInput
	}

	item, err := s.requests.FindByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	if item.UserID != userOID {
		return nil, ErrRequestNotFound
	}

	view := toRequestView(item)
	return &view, nil
}

func (s *Service) CompleteSuccess(ctx context.Context, requestID, resultText string) (*RequestView, error) {
	if s == nil || s.requests == nil {
		return nil, ErrInvalidConfig
	}
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, ErrInvalidInput
	}
	normalizedResult, err := inpututil.RequireTrimmed(resultText)
	if err != nil {
		return nil, ErrInvalidInput
	}

	item, err := s.requests.FindByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	if isFinalStatus(item.Status) {
		view := toRequestView(item)
		return &view, nil
	}

	updated, err := s.requests.UpdateResult(ctx, requestID, normalizedResult)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	view := toRequestView(updated)
	return &view, nil
}

func (s *Service) CompleteError(
	ctx context.Context,
	requestID string,
	errorCode string,
	errorMessage string,
) (*RequestView, error) {
	if s == nil || s.requests == nil {
		return nil, ErrInvalidConfig
	}
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, ErrInvalidInput
	}

	normalizedCode, err := inpututil.RequireTrimmed(errorCode)
	if err != nil {
		return nil, ErrInvalidInput
	}
	normalizedMessage, err := inpututil.RequireTrimmed(errorMessage)
	if err != nil {
		return nil, ErrInvalidInput
	}

	item, err := s.requests.FindByRequestID(ctx, requestID)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	if isFinalStatus(item.Status) {
		view := toRequestView(item)
		return &view, nil
	}

	updated, err := s.requests.UpdateError(ctx, requestID, normalizedCode, normalizedMessage)
	if err != nil {
		if errors.Is(err, mongorepo.ErrNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}

	view := toRequestView(updated)
	return &view, nil
}

func isFinalStatus(status models.AskRequestStatus) bool {
	return status == models.AskRequestStatusSuccess || status == models.AskRequestStatusError
}

func normalizeStatus(raw string) models.AskRequestStatus {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case string(models.AskRequestStatusAccepted):
		return models.AskRequestStatusAccepted
	case string(models.AskRequestStatusProcessing):
		return models.AskRequestStatusProcessing
	case string(models.AskRequestStatusQueued):
		return models.AskRequestStatusQueued
	case string(models.AskRequestStatusSuccess):
		return models.AskRequestStatusSuccess
	case string(models.AskRequestStatusError):
		return models.AskRequestStatusError
	default:
		return models.AskRequestStatusQueued
	}
}

func (s *Service) callbackURL() (string, error) {
	base := strings.TrimSpace(s.backendPublicBaseURL)
	if base == "" {
		return "", ErrInvalidConfig
	}
	return strings.TrimRight(base, "/") + "/api/gateway/callback", nil
}

func buildRequestID() string {
	return "req-" + primitive.NewObjectID().Hex()
}

func toRequestView(item *models.AskRequest) RequestView {
	if item == nil {
		return RequestView{}
	}
	return RequestView{
		RequestID:    item.RequestID,
		Status:       string(item.Status),
		ResultText:   item.ResultText,
		ErrorCode:    item.ErrorCode,
		ErrorMessage: item.ErrorMessage,
	}
}

func mapGatewayDispatchError(err error) (string, error) {
	if err == nil {
		return "gateway_request_failed", ErrGatewayUnavailable
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return "gateway_timeout", ErrGatewayTimeout
	}

	if errors.Is(err, gatewayclient.ErrInvalidResponse) {
		return "gateway_invalid_response", ErrGatewayBadResponse
	}

	if errors.Is(err, gatewayclient.ErrInvalidInput) {
		return "gateway_invalid_request", ErrInvalidConfig
	}

	var reqErr *gatewayclient.Error
	if errors.As(err, &reqErr) {
		switch {
		case reqErr.StatusCode >= http.StatusBadRequest && reqErr.StatusCode < http.StatusInternalServerError:
			return "gateway_rejected", ErrGatewayRejected
		case reqErr.StatusCode >= http.StatusInternalServerError:
			return "gateway_unavailable", ErrGatewayUnavailable
		}
	}

	return "gateway_request_failed", ErrGatewayUnavailable
}
