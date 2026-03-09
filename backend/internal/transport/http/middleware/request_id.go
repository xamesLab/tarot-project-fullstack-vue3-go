package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const requestIDHeader = "X-Request-ID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(requestIDHeader)
		if reqID == "" {
			reqID = generateRequestID()
		}

		c.Set("request_id", reqID)
		c.Writer.Header().Set(requestIDHeader, reqID)
		c.Next()
	}
}

func generateRequestID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "req_fallback"
	}
	return hex.EncodeToString(buf)
}
