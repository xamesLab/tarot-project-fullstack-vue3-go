package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/transport/http/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}

		if len(c.Errors) == 0 {
			return
		}

		response.JSONError(c, http.StatusInternalServerError, "internal_error", c.Errors.Last().Error())
	}
}
