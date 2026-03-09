package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/transport/http/response"
)

func RoleRequired(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		role, _ := c.Get("auth_role")
		roleStr, _ := role.(string)
		if roleStr == "" {
			response.JSONError(c, http.StatusForbidden, "forbidden", "missing user role")
			return
		}

		if _, ok := allowed[roleStr]; !ok {
			response.JSONError(c, http.StatusForbidden, "forbidden", "insufficient role")
			return
		}

		c.Next()
	}
}
