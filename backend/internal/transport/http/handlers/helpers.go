package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/transport/http/response"
)

type errorSpec struct {
	target  error
	status  int
	code    string
	message string
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		response.JSONError(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return false
	}
	return true
}

func currentUserID(c *gin.Context) string {
	userID, _ := c.Get("auth_user_id")
	userIDStr, _ := userID.(string)
	return userIDStr
}

func parseInt64Query(c *gin.Context, key string) (int64, error) {
	raw := c.Query(key)
	if raw == "" {
		return 0, nil
	}
	return strconv.ParseInt(raw, 10, 64)
}

func respondInternalError(c *gin.Context) {
	response.JSONError(c, http.StatusInternalServerError, "internal_error", "internal server error")
}

func respondByError(c *gin.Context, err error, specs ...errorSpec) {
	for _, spec := range specs {
		if errors.Is(err, spec.target) {
			response.JSONError(c, spec.status, spec.code, spec.message)
			return
		}
	}

	respondInternalError(c)
}
