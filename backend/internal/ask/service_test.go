package ask

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"tarotdaily/backend/internal/domain/models"
	"tarotdaily/backend/internal/gatewayclient"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
)

type fakeUsers struct {
	user *models.User
}

func (f *fakeUsers) FindByID(_ context.Context, id primitive.ObjectID) (*models.User, error) {
	if f.user == nil || f.user.ID != id {
		return nil, mongorepo.ErrNotFound
	}
	return f.user, nil
}

type fakeAskStore struct {
	mu    sync.Mutex
	items map[string]*models.AskRequest
}

func newFakeAskStore() *fakeAskStore {
	return &fakeAskStore{items: map[string]*models.AskRequest{}}
}

func (f *fakeAskStore) Create(_ context.Context, item *models.AskRequest) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.items[item.RequestID]; exists {
		return errors.New("duplicate")
	}
	copyItem := *item
	f.items[item.RequestID] = &copyItem
	return nil
}

func (f *fakeAskStore) FindByRequestID(_ context.Context, requestID string) (*models.AskRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.items[requestID]
	if !ok {
		return nil, mongorepo.ErrNotFound
	}
	copyItem := *item
	return &copyItem, nil
}

func (f *fakeAskStore) UpdateStatus(_ context.Context, requestID string, status models.AskRequestStatus) (*models.AskRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.items[requestID]
	if !ok {
		return nil, mongorepo.ErrNotFound
	}
	item.Status = status
	copyItem := *item
	return &copyItem, nil
}

func (f *fakeAskStore) UpdateResult(_ context.Context, requestID string, resultText string) (*models.AskRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.items[requestID]
	if !ok {
		return nil, mongorepo.ErrNotFound
	}
	now := time.Now().UTC()
	item.Status = models.AskRequestStatusSuccess
	item.ResultText = resultText
	item.ErrorCode = ""
	item.ErrorMessage = ""
	item.CompletedAt = &now
	copyItem := *item
	return &copyItem, nil
}

func (f *fakeAskStore) UpdateError(_ context.Context, requestID string, errorCode string, errorMessage string) (*models.AskRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	item, ok := f.items[requestID]
	if !ok {
		return nil, mongorepo.ErrNotFound
	}
	now := time.Now().UTC()
	item.Status = models.AskRequestStatusError
	item.ErrorCode = errorCode
	item.ErrorMessage = errorMessage
	item.ResultText = ""
	item.CompletedAt = &now
	copyItem := *item
	return &copyItem, nil
}

type fakeGateway struct {
	resp        *gatewayclient.RequestResponse
	err         error
	lastPayload gatewayclient.RequestPayload
}

func (f *fakeGateway) Request(_ context.Context, payload gatewayclient.RequestPayload) (*gatewayclient.RequestResponse, error) {
	f.lastPayload = payload
	if f.err != nil {
		return nil, f.err
	}
	return f.resp, nil
}

func TestServiceCreateSuccess(t *testing.T) {
	userID := primitive.NewObjectID()
	users := &fakeUsers{
		user: &models.User{
			ID: userID,
		},
	}
	store := newFakeAskStore()
	gateway := &fakeGateway{
		resp: &gatewayclient.RequestResponse{
			Status: "queued",
		},
	}

	svc := NewService(users, store, gateway, "http://127.0.0.1:8080")
	got, err := svc.Create(context.Background(), userID.Hex(), "hello gateway")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.Status != "queued" {
		t.Fatalf("Create() status = %q, want queued", got.Status)
	}
	if got.RequestID == "" {
		t.Fatalf("Create() request id is empty")
	}
	if gateway.lastPayload.CallbackURL != "http://127.0.0.1:8080/api/gateway/callback" {
		t.Fatalf("CallbackURL = %q", gateway.lastPayload.CallbackURL)
	}
}

func TestServiceCreateGatewayTimeoutMapping(t *testing.T) {
	userID := primitive.NewObjectID()
	users := &fakeUsers{
		user: &models.User{
			ID: userID,
		},
	}
	store := newFakeAskStore()
	gateway := &fakeGateway{
		err: context.DeadlineExceeded,
	}

	svc := NewService(users, store, gateway, "http://127.0.0.1:8080")
	_, err := svc.Create(context.Background(), userID.Hex(), "hello timeout")
	if !errors.Is(err, ErrGatewayTimeout) {
		t.Fatalf("Create() error = %v, want ErrGatewayTimeout", err)
	}

	if len(store.items) != 1 {
		t.Fatalf("store items len = %d, want 1", len(store.items))
	}
	for _, item := range store.items {
		if item.Status != models.AskRequestStatusError {
			t.Fatalf("stored status = %q, want error", item.Status)
		}
		if item.ErrorCode != "gateway_timeout" {
			t.Fatalf("stored errorCode = %q, want gateway_timeout", item.ErrorCode)
		}
	}
}

func TestServiceCompleteSuccessIdempotency(t *testing.T) {
	store := newFakeAskStore()
	requestID := "req-1"
	store.items[requestID] = &models.AskRequest{
		RequestID: requestID,
		Status:    models.AskRequestStatusQueued,
	}

	svc := NewService(nil, store, nil, "")
	first, err := svc.CompleteSuccess(context.Background(), requestID, "first result")
	if err != nil {
		t.Fatalf("CompleteSuccess(first) error = %v", err)
	}
	if first.ResultText != "first result" || first.Status != "success" {
		t.Fatalf("first completion = %+v", first)
	}

	second, err := svc.CompleteSuccess(context.Background(), requestID, "second result")
	if err != nil {
		t.Fatalf("CompleteSuccess(second) error = %v", err)
	}
	if second.ResultText != "first result" {
		t.Fatalf("second completion result = %q, want first result", second.ResultText)
	}
}

func TestServiceCompleteErrorIdempotency(t *testing.T) {
	store := newFakeAskStore()
	requestID := "req-2"
	store.items[requestID] = &models.AskRequest{
		RequestID: requestID,
		Status:    models.AskRequestStatusQueued,
	}

	svc := NewService(nil, store, nil, "")
	first, err := svc.CompleteError(context.Background(), requestID, "agent_timeout", "first message")
	if err != nil {
		t.Fatalf("CompleteError(first) error = %v", err)
	}
	if first.ErrorCode != "agent_timeout" || first.Status != "error" {
		t.Fatalf("first completion = %+v", first)
	}

	second, err := svc.CompleteError(context.Background(), requestID, "agent_failure", "second message")
	if err != nil {
		t.Fatalf("CompleteError(second) error = %v", err)
	}
	if second.ErrorCode != "agent_timeout" {
		t.Fatalf("second completion errorCode = %q, want agent_timeout", second.ErrorCode)
	}
}

func TestMapGatewayDispatchError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode string
		wantErr  error
	}{
		{
			name:     "deadline exceeded",
			err:      context.DeadlineExceeded,
			wantCode: "gateway_timeout",
			wantErr:  ErrGatewayTimeout,
		},
		{
			name:     "bad response",
			err:      fmtWrap(gatewayclient.ErrInvalidResponse),
			wantCode: "gateway_invalid_response",
			wantErr:  ErrGatewayBadResponse,
		},
		{
			name:     "invalid request config",
			err:      fmtWrap(gatewayclient.ErrInvalidInput),
			wantCode: "gateway_invalid_request",
			wantErr:  ErrInvalidConfig,
		},
		{
			name:     "gateway rejected",
			err:      &gatewayclient.Error{StatusCode: http.StatusBadRequest, Err: gatewayclient.ErrUnexpectedStatus},
			wantCode: "gateway_rejected",
			wantErr:  ErrGatewayRejected,
		},
		{
			name:     "gateway unavailable",
			err:      &gatewayclient.Error{StatusCode: http.StatusBadGateway, Err: gatewayclient.ErrUnexpectedStatus},
			wantCode: "gateway_unavailable",
			wantErr:  ErrGatewayUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCode, gotErr := mapGatewayDispatchError(tt.err)
			if gotCode != tt.wantCode {
				t.Fatalf("code = %q, want %q", gotCode, tt.wantCode)
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf("error = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func fmtWrap(err error) error {
	return errors.Join(err)
}
