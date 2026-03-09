package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwtplatform "tarotdaily/backend/internal/platform/jwt"
	"tarotdaily/backend/internal/transport/http/response"
)

func AuthRequired(jwtManager *jwtplatform.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			response.JSONError(c, http.StatusUnauthorized, "unauthorized", "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid authorization format")
			return
		}

		claims, err := jwtManager.Parse(strings.TrimSpace(parts[1]))
		if err != nil {
			response.JSONError(c, http.StatusUnauthorized, "unauthorized", "invalid or expired token")
			return
		}

		c.Set("auth_user_id", claims.UserID)
		c.Set("auth_login", claims.Login)
		c.Set("auth_role", claims.Role)
		c.Next()
	}
}
