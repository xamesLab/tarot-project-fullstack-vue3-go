package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/transport/http/response"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		response.JSONError(c, http.StatusInternalServerError, "internal_error", "internal server error")
	})
}
