package middleware

import (
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(frontendURL string) gin.HandlerFunc {
	allowOrigins := allowedOrigins(frontendURL)

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept-Language"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func allowedOrigins(frontendURL string) []string {
	if frontendURL == "" {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}

	parsed, err := url.Parse(frontendURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return []string{frontendURL}
	}

	origins := map[string]struct{}{
		frontendURL: {},
	}

	switch parsed.Hostname() {
	case "localhost":
		parsed.Host = "127.0.0.1:" + parsed.Port()
		origins[parsed.String()] = struct{}{}
	case "127.0.0.1":
		parsed.Host = "localhost:" + parsed.Port()
		origins[parsed.String()] = struct{}{}
	}

	result := make([]string, 0, len(origins))
	for origin := range origins {
		result = append(result, origin)
	}

	return result
}
