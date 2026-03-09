package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/roles"
)

type AdminRolesHandler struct {
	service *roles.Service
}

type upsertRolePolicyRequest struct {
	Role             string `json:"role"`
	DailySpreadLimit int    `json:"dailySpreadLimit"`
}

type assignUserRoleRequest struct {
	Login string `json:"login"`
	Role  string `json:"role"`
}

func NewAdminRolesHandler(service *roles.Service) *AdminRolesHandler {
	return &AdminRolesHandler{service: service}
}

func (h *AdminRolesHandler) List(c *gin.Context) {
	items, err := h.service.ListPolicies(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *AdminRolesHandler) Upsert(c *gin.Context) {
	var req upsertRolePolicyRequest
	if !bindJSON(c, &req) {
		return
	}

	item, err := h.service.UpsertPolicy(c.Request.Context(), req.Role, req.DailySpreadLimit)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: roles.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid role or dailySpreadLimit"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *AdminRolesHandler) Assign(c *gin.Context) {
	var req assignUserRoleRequest
	if !bindJSON(c, &req) {
		return
	}

	item, err := h.service.AssignRoleByLogin(c.Request.Context(), req.Login, req.Role)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: roles.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid login or role"},
			errorSpec{target: roles.ErrRolePolicyNotFound, status: http.StatusBadRequest, code: "invalid_role", message: "role policy not found"},
			errorSpec{target: roles.ErrUserNotFound, status: http.StatusNotFound, code: "user_not_found", message: "user not found"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}
