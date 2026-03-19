package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGatewayCallbackVerifyMissingConfig(t *testing.T) {
	status := runGatewayVerifyRequest(
		t,
		GatewayCallbackVerify("", nil),
		`{"request_id":"req-1","status":"success","text":"ok"}`,
		map[string]string{},
		"127.0.0.1:10001",
	)

	if status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", status)
	}
}

func TestGatewayCallbackVerifyAllowlist(t *testing.T) {
	t.Run("allowlist ok", func(t *testing.T) {
		status := runGatewayVerifyRequest(
			t,
			GatewayCallbackVerify("", []string{"127.0.0.1", "::1"}),
			`{"request_id":"req-1","status":"success","text":"ok"}`,
			map[string]string{},
			"127.0.0.1:10001",
		)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want 200", status)
		}
	})

	t.Run("allowlist deny", func(t *testing.T) {
		status := runGatewayVerifyRequest(
			t,
			GatewayCallbackVerify("", []string{"127.0.0.1"}),
			`{"request_id":"req-1","status":"success","text":"ok"}`,
			map[string]string{},
			"10.10.10.10:10001",
		)
		if status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", status)
		}
	})
}

func TestGatewayCallbackVerifyHMAC(t *testing.T) {
	body := `{"request_id":"req-1","status":"success","text":"ok"}`
	secret := "test-secret"
	signature := buildHMAC(body, secret)

	t.Run("hmac ok", func(t *testing.T) {
		status := runGatewayVerifyRequest(
			t,
			GatewayCallbackVerify(secret, nil),
			body,
			map[string]string{
				callbackSignatureHeader: "sha256=" + signature,
			},
			"127.0.0.1:10001",
		)
		if status != http.StatusOK {
			t.Fatalf("status = %d, want 200", status)
		}
	})

	t.Run("hmac invalid", func(t *testing.T) {
		status := runGatewayVerifyRequest(
			t,
			GatewayCallbackVerify(secret, nil),
			body,
			map[string]string{
				callbackSignatureHeader: "sha256=deadbeef",
			},
			"127.0.0.1:10001",
		)
		if status != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", status)
		}
	})
}

func runGatewayVerifyRequest(
	t *testing.T,
	mw gin.HandlerFunc,
	body string,
	headers map[string]string,
	remoteAddr string,
) int {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/gateway/callback", mw, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/api/gateway/callback", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec.Code
}

func buildHMAC(body string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}
