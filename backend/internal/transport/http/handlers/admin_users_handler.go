package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/adminusers"
)

type AdminUsersHandler struct {
	service *adminusers.Service
}

func NewAdminUsersHandler(service *adminusers.Service) *AdminUsersHandler {
	return &AdminUsersHandler{service: service}
}

func (h *AdminUsersHandler) List(c *gin.Context) {
	page, err := parseInt64Query(c, "page")
	if err != nil {
		respondByError(c, adminusers.ErrInvalidInput, errorSpec{
			target:  adminusers.ErrInvalidInput,
			status:  http.StatusBadRequest,
			code:    "invalid_input",
			message: "invalid page",
		})
		return
	}

	pageSize, err := parseInt64Query(c, "pageSize")
	if err != nil {
		respondByError(c, adminusers.ErrInvalidInput, errorSpec{
			target:  adminusers.ErrInvalidInput,
			status:  http.StatusBadRequest,
			code:    "invalid_input",
			message: "invalid pageSize",
		})
		return
	}

	result, err := h.service.ListUsers(c.Request.Context(), adminusers.ListUsersInput{
		Page:     page,
		PageSize: pageSize,
		Role:     c.Query("role"),
		Status:   c.Query("status"),
		Sort:     c.Query("sort"),
		Name:     c.Query("name"),
	})
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: adminusers.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid query parameters"},
		)
		return
	}

	c.JSON(http.StatusOK, result)
}
