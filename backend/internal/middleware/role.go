package middleware

import (
	"net/http"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   services.ErrForbidden.Error(),
			"message": "Admin access required",
		})
		c.Abort()
	}
}
