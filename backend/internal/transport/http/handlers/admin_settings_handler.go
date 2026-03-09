package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/adminsettings"
)

type AdminSettingsHandler struct {
	service *adminsettings.Service
}

type patchAdminSettingsRequest struct {
	ProjectTitle    *string `json:"projectTitle"`
	SupportEmail    *string `json:"supportEmail"`
	MaintenanceMode *bool   `json:"maintenanceMode"`
	DayResetHour    *int    `json:"dayResetHour"`
	MaxUploadMB     *int    `json:"maxUploadMb"`
}

func NewAdminSettingsHandler(service *adminsettings.Service) *AdminSettingsHandler {
	return &AdminSettingsHandler{service: service}
}

func (h *AdminSettingsHandler) Get(c *gin.Context) {
	item, err := h.service.Get(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *AdminSettingsHandler) Patch(c *gin.Context) {
	var req patchAdminSettingsRequest
	if !bindJSON(c, &req) {
		return
	}

	item, err := h.service.Update(c.Request.Context(), adminsettings.UpdateInput{
		ProjectTitle:    req.ProjectTitle,
		SupportEmail:    req.SupportEmail,
		MaintenanceMode: req.MaintenanceMode,
		DayResetHour:    req.DayResetHour,
		MaxUploadMB:     req.MaxUploadMB,
	})
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: adminsettings.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid settings payload"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}
