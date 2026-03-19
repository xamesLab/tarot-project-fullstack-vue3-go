package gatewayclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidResponse  = errors.New("invalid response")
	ErrUnexpectedStatus = errors.New("unexpected status code")
)

type Error struct {
	StatusCode int
	Body       string
	Err        error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.StatusCode > 0 {
		if strings.TrimSpace(e.Body) != "" {
			return fmt.Sprintf("gateway request failed: status=%d body=%s", e.StatusCode, e.Body)
		}
		return fmt.Sprintf("gateway request failed: status=%d", e.StatusCode)
	}
	if e.Err != nil {
		return fmt.Sprintf("gateway request failed: %v", e.Err)
	}
	return "gateway request failed"
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL     string
	requestPath string
	httpClient  HTTPDoer
}

type HealthResponse struct {
	Status string  `json:"status"`
	Uptime float64 `json:"uptime,omitempty"`
}

type RequestMetadata struct {
	UserID string `json:"user_id"`
}

type RequestPayload struct {
	Text        string          `json:"text"`
	CallbackURL string          `json:"callback_url"`
	RequestID   string          `json:"request_id"`
	Metadata    RequestMetadata `json:"metadata"`
}

type RequestResponse struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

func New(baseURL, requestPath string, timeout time.Duration) (*Client, error) {
	if timeout <= 0 {
		return nil, ErrInvalidInput
	}
	return NewWithHTTPClient(baseURL, requestPath, &http.Client{Timeout: timeout})
}

func NewWithHTTPClient(baseURL, requestPath string, httpClient HTTPDoer) (*Client, error) {
	normalizedBase := strings.TrimSpace(baseURL)
	normalizedPath := normalizePath(requestPath)

	if normalizedBase == "" || normalizedPath == "" || httpClient == nil {
		return nil, ErrInvalidInput
	}

	return &Client{
		baseURL:     strings.TrimRight(normalizedBase, "/"),
		requestPath: normalizedPath,
		httpClient:  httpClient,
	}, nil
}

func (c *Client) Request(ctx context.Context, payload RequestPayload) (*RequestResponse, error) {
	if c == nil || c.httpClient == nil {
		return nil, ErrInvalidInput
	}
	if err := validatePayload(payload); err != nil {
		return nil, err
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal gateway payload: %w", err)
	}

	url := c.baseURL + c.requestPath
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("build gateway request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024))
	if err != nil {
		return nil, fmt.Errorf("read gateway response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(body)),
			Err:        ErrUnexpectedStatus,
		}
	}

	var out RequestResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("%w: decode gateway response: %v", ErrInvalidResponse, err)
	}

	if strings.TrimSpace(out.RequestID) == "" {
		out.RequestID = payload.RequestID
	}
	if strings.TrimSpace(out.RequestID) == "" {
		return nil, fmt.Errorf("%w: request_id is empty", ErrInvalidResponse)
	}
	if strings.TrimSpace(out.Status) == "" {
		out.Status = "queued"
	}

	return &out, nil
}

func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	if c == nil || c.httpClient == nil {
		return nil, ErrInvalidInput
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return nil, fmt.Errorf("build gateway health request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8*1024))
	if err != nil {
		return nil, fmt.Errorf("read gateway health response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, &Error{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(body)),
			Err:        ErrUnexpectedStatus,
		}
	}

	var out HealthResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, fmt.Errorf("%w: decode gateway health response: %v", ErrInvalidResponse, err)
	}

	if strings.TrimSpace(out.Status) == "" {
		return nil, fmt.Errorf("%w: status is empty", ErrInvalidResponse)
	}

	return &out, nil
}

func validatePayload(payload RequestPayload) error {
	if strings.TrimSpace(payload.Text) == "" {
		return ErrInvalidInput
	}
	if strings.TrimSpace(payload.CallbackURL) == "" {
		return ErrInvalidInput
	}
	if strings.TrimSpace(payload.RequestID) == "" {
		return ErrInvalidInput
	}
	if strings.TrimSpace(payload.Metadata.UserID) == "" {
		return ErrInvalidInput
	}
	return nil
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/request"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}
