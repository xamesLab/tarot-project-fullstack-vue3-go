package response

import "github.com/gin-gonic/gin"

type ErrorBody struct {
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"requestId,omitempty"`
}

func JSONError(c *gin.Context, status int, code, message string) {
	body := ErrorBody{
		Error:   code,
		Message: message,
	}
	if reqID := c.GetString("request_id"); reqID != "" {
		body.RequestID = reqID
	}

	c.AbortWithStatusJSON(status, body)
}
