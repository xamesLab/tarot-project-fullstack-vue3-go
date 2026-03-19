package gatewayclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientRequestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/request" {
			t.Fatalf("path = %s, want /request", r.URL.Path)
		}

		var payload RequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if payload.RequestID != "req-123" {
			t.Fatalf("payload request_id = %q", payload.RequestID)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"req-123","status":"queued"}`))
	}))
	defer server.Close()

	client, err := New(server.URL, "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := client.Request(context.Background(), RequestPayload{
		Text:        "hello",
		CallbackURL: "http://127.0.0.1:8080/api/gateway/callback",
		RequestID:   "req-123",
		Metadata: RequestMetadata{
			UserID: "u-1",
		},
	})
	if err != nil {
		t.Fatalf("Request() error = %v", err)
	}
	if got.RequestID != "req-123" || got.Status != "queued" {
		t.Fatalf("Request() = %+v", got)
	}
}

func TestClientRequestUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway request", http.StatusBadRequest)
	}))
	defer server.Close()

	client, err := New(server.URL, "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Request(context.Background(), RequestPayload{
		Text:        "hello",
		CallbackURL: "http://127.0.0.1:8080/api/gateway/callback",
		RequestID:   "req-123",
		Metadata: RequestMetadata{
			UserID: "u-1",
		},
	})
	if err == nil {
		t.Fatalf("Request() error is nil")
	}
	if !errors.Is(err, ErrUnexpectedStatus) {
		t.Fatalf("error = %v, want ErrUnexpectedStatus", err)
	}
	var reqErr *Error
	if !errors.As(err, &reqErr) {
		t.Fatalf("error = %T, want *Error", err)
	}
	if reqErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status code = %d, want 400", reqErr.StatusCode)
	}
}

func TestClientRequestInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	client, err := New(server.URL, "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Request(context.Background(), RequestPayload{
		Text:        "hello",
		CallbackURL: "http://127.0.0.1:8080/api/gateway/callback",
		RequestID:   "req-123",
		Metadata: RequestMetadata{
			UserID: "u-1",
		},
	})
	if err == nil {
		t.Fatalf("Request() error is nil")
	}
	if !errors.Is(err, ErrInvalidResponse) {
		t.Fatalf("error = %v, want ErrInvalidResponse", err)
	}
}

func TestClientRequestPayloadValidation(t *testing.T) {
	client, err := New("http://127.0.0.1:18081", "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Request(context.Background(), RequestPayload{
		Text:        " ",
		CallbackURL: "http://127.0.0.1:8080/api/gateway/callback",
		RequestID:   "req-123",
		Metadata: RequestMetadata{
			UserID: "u-1",
		},
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("error = %v, want ErrInvalidInput", err)
	}
}

func TestClientHealthSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/health" {
			t.Fatalf("path = %s, want /health", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","uptime":123.45}`))
	}))
	defer server.Close()

	client, err := New(server.URL, "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Health() error = %v", err)
	}
	if got.Status != "ok" {
		t.Fatalf("status = %q, want ok", got.Status)
	}
	if got.Uptime != 123.45 {
		t.Fatalf("uptime = %v, want 123.45", got.Uptime)
	}
}

func TestClientHealthUnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client, err := New(server.URL, "/request", time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	_, err = client.Health(context.Background())
	if err == nil {
		t.Fatalf("Health() error is nil")
	}
	if !errors.Is(err, ErrUnexpectedStatus) {
		t.Fatalf("error = %v, want ErrUnexpectedStatus", err)
	}
}

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{in: "", want: "/request"},
		{in: "request", want: "/request"},
		{in: "/request", want: "/request"},
		{in: "  /api/request  ", want: "/api/request"},
	}

	for _, tt := range tests {
		got := normalizePath(tt.in)
		if got != tt.want {
			t.Fatalf("normalizePath(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if !strings.HasPrefix(got, "/") {
			t.Fatalf("normalizePath(%q) must start with /", tt.in)
		}
	}
}
