package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/ask"
	"tarotdaily/backend/internal/transport/http/response"
)

type AskHandler struct {
	service *ask.Service
}

type createAskRequest struct {
	Text string `json:"text"`
}

type gatewayCallbackRequest struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Text      string `json:"text"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}

func NewAskHandler(service *ask.Service) *AskHandler {
	return &AskHandler{
		service: service,
	}
}

func (h *AskHandler) Create(c *gin.Context) {
	var req createAskRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.service.Create(c.Request.Context(), currentUserID(c), req.Text)
	if err != nil {
		h.mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, result)
}

func (h *AskHandler) Get(c *gin.Context) {
	result, err := h.service.Get(c.Request.Context(), currentUserID(c), c.Param("request_id"))
	if err != nil {
		h.mapServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AskHandler) GatewayCallback(c *gin.Context) {
	var req gatewayCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSONError(c, http.StatusBadRequest, "invalid_request", "invalid request body")
		return
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	switch status {
	case "success":
		if _, err := h.service.CompleteSuccess(c.Request.Context(), req.RequestID, req.Text); err != nil {
			h.mapServiceError(c, err)
			return
		}
	case "error":
		if _, err := h.service.CompleteError(c.Request.Context(), req.RequestID, req.Error, req.Message); err != nil {
			h.mapServiceError(c, err)
			return
		}
	default:
		response.JSONError(c, http.StatusBadRequest, "invalid_input", "invalid callback status")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id": strings.TrimSpace(req.RequestID),
		"status":     "ok",
	})
}

func (h *AskHandler) mapServiceError(c *gin.Context, err error) {
	respondByError(
		c,
		err,
		errorSpec{target: ask.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid ask input"},
		errorSpec{target: ask.ErrRequestNotFound, status: http.StatusNotFound, code: "not_found", message: "request not found"},
		errorSpec{target: ask.ErrGatewayTimeout, status: http.StatusGatewayTimeout, code: "gateway_timeout", message: "gateway request timed out"},
		errorSpec{target: ask.ErrGatewayRejected, status: http.StatusBadGateway, code: "gateway_rejected", message: "gateway rejected request"},
		errorSpec{target: ask.ErrGatewayBadResponse, status: http.StatusBadGateway, code: "gateway_bad_response", message: "gateway returned invalid response"},
		errorSpec{target: ask.ErrGatewayUnavailable, status: http.StatusBadGateway, code: "gateway_unavailable", message: "gateway is unavailable"},
		errorSpec{target: ask.ErrInvalidConfig, status: http.StatusInternalServerError, code: "invalid_config", message: "ask integration is not configured"},
	)
}
