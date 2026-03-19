package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/transport/http/response"
)

const callbackSignatureHeader = "X-Signature"

func GatewayCallbackVerify(hmacSecret string, allowlist []string) gin.HandlerFunc {
	secret := strings.TrimSpace(hmacSecret)
	allowlistMap := make(map[string]struct{}, len(allowlist))
	for _, item := range allowlist {
		normalized := strings.TrimSpace(item)
		if normalized != "" {
			allowlistMap[normalized] = struct{}{}
		}
	}

	hasHMAC := secret != ""
	hasAllowlist := len(allowlistMap) > 0

	return func(c *gin.Context) {
		if !hasHMAC && !hasAllowlist {
			response.JSONError(c, http.StatusInternalServerError, "invalid_config", "callback verification is not configured")
			return
		}

		if hasAllowlist {
			clientIP := strings.TrimSpace(c.ClientIP())
			if _, ok := allowlistMap[clientIP]; !ok {
				response.JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid callback signature or source")
				return
			}
		}

		if hasHMAC {
			rawBody, err := c.GetRawData()
			if err != nil {
				response.JSONError(c, http.StatusBadRequest, "invalid_request", "invalid request body")
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

			signature := strings.TrimSpace(c.GetHeader(callbackSignatureHeader))
			if signature == "" || !verifyHMAC(rawBody, signature, secret) {
				response.JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid callback signature or source")
				return
			}
		}

		c.Next()
	}
}

func verifyHMAC(rawBody []byte, signatureHeader, secret string) bool {
	signature := strings.TrimSpace(signatureHeader)
	signature = strings.TrimPrefix(signature, "sha256=")
	signature = strings.TrimSpace(signature)

	received, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write(rawBody); err != nil {
		return false
	}
	expected := mac.Sum(nil)

	return hmac.Equal(received, expected)
}
