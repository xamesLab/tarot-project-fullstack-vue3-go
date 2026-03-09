package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccessHandler struct{}

func NewAccessHandler() *AccessHandler {
	return &AccessHandler{}
}

func (h *AccessHandler) AdminAccess(c *gin.Context) {
	login, _ := c.Get("auth_login")
	role, _ := c.Get("auth_role")
	userID, _ := c.Get("auth_user_id")

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"scope":  "admin",
		"user": gin.H{
			"id":    userID,
			"login": login,
			"role":  role,
		},
	})
}
