package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/auth"
)

type AuthHandler struct {
	service *auth.Service
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.service.Register(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if !bindJSON(c, &req) {
		return
	}

	result, err := h.service.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) mapError(c *gin.Context, err error) {
	respondByError(
		c,
		err,
		errorSpec{target: auth.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid login or password"},
		errorSpec{target: auth.ErrLoginAlreadyExists, status: http.StatusConflict, code: "login_exists", message: "login already exists"},
		errorSpec{target: auth.ErrInvalidCredentials, status: http.StatusUnauthorized, code: "invalid_credentials", message: "invalid credentials"},
	)
}
