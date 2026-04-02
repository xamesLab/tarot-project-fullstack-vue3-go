package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	adminService *services.AdminService
	validate     *validator.Validate
}

func NewAdminHandler(adminService *services.AdminService, validate *validator.Validate) *AdminHandler {
	return &AdminHandler{adminService: adminService, validate: validate}
}

func (h *AdminHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/access-check", h.accessCheck)
	router.GET("/dashboard", h.getDashboard)
	router.POST("/decks", h.createDeck)
	router.POST("/card-backs", h.createCardBack)
	router.GET("/users", h.listUsers)
	router.GET("/decks", h.listDecks)
	router.GET("/card-backs", h.listCardBacks)
	router.GET("/decks/:deckId", h.getDeck)
	router.PATCH("/decks/:deckId", h.updateDeck)
	router.PATCH("/users/:userId/role", h.updateUserRole)
	router.PATCH("/users/:userId/status", h.updateUserStatus)
	router.GET("/roles", h.listRolePolicies)
	router.PATCH("/roles/:role", h.updateRolePolicy)
	router.GET("/settings", h.getSettings)
	router.PATCH("/settings", h.updateSettings)
}

func (h *AdminHandler) accessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"has_access": true,
	})
}

func (h *AdminHandler) getDashboard(c *gin.Context) {
	metrics, err := h.adminService.GetDashboardMetrics(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"metrics": metrics,
	})
}

func (h *AdminHandler) listUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, pagination, err := h.adminService.ListUsers(
		c.Request.Context(),
		c.Query("search"),
		c.Query("role"),
		c.Query("status"),
		page,
		limit,
	)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"users":      users,
		"pagination": pagination,
	})
}

func (h *AdminHandler) createDeck(c *gin.Context) {
	nameEN := strings.TrimSpace(c.PostForm("name_en"))
	nameRU := strings.TrimSpace(c.PostForm("name_ru"))
	descriptionEN := strings.TrimSpace(c.PostForm("description_en"))
	descriptionRU := strings.TrimSpace(c.PostForm("description_ru"))
	author := strings.TrimSpace(c.PostForm("author"))
	yearValue := strings.TrimSpace(c.PostForm("year"))
	tags := parseCSVValues(c.PostForm("tags"))

	year := 0
	if yearValue != "" {
		parsedYear, err := strconv.Atoi(yearValue)
		if err != nil || parsedYear < 0 {
			respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Year must be a positive number")
			return
		}
		year = parsedYear
	}

	zipFile, err := c.FormFile("cards_zip")
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "cards_zip is required")
		return
	}

	tempRoot := filepath.Join("uploads", "temp")
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to prepare upload directory")
		return
	}

	tempDir, err := os.MkdirTemp(tempRoot, "deck-upload-*")
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to prepare upload directory")
		return
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, "cards.zip")
	if err := c.SaveUploadedFile(zipFile, zipPath); err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to save uploaded ZIP")
		return
	}

	deck, err := h.adminService.CreateDeck(c.Request.Context(), models.AdminCreateDeckRequest{
		Name: models.LocalizedText{
			EN: nameEN,
			RU: nameRU,
		},
		Description: models.LocalizedText{
			EN: descriptionEN,
			RU: descriptionRU,
		},
		Author:  author,
		Year:    year,
		Tags:    tags,
		ZipPath: zipPath,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"deck":    deck,
		"message": "Deck created and processed successfully",
	})
}

func (h *AdminHandler) createCardBack(c *gin.Context) {
	nameEN := strings.TrimSpace(c.PostForm("name_en"))
	nameRU := strings.TrimSpace(c.PostForm("name_ru"))
	descriptionEN := strings.TrimSpace(c.PostForm("description_en"))
	descriptionRU := strings.TrimSpace(c.PostForm("description_ru"))
	tags := parseCSVValues(c.PostForm("tags"))

	imageFile, err := c.FormFile("back_image")
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "back_image is required")
		return
	}

	tempRoot := filepath.Join("uploads", "temp")
	if err := os.MkdirAll(tempRoot, 0o755); err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to prepare upload directory")
		return
	}

	tempDir, err := os.MkdirTemp(tempRoot, "card-back-upload-*")
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to prepare upload directory")
		return
	}
	defer os.RemoveAll(tempDir)

	imagePath := filepath.Join(tempDir, "back"+strings.ToLower(filepath.Ext(imageFile.Filename)))
	if err := c.SaveUploadedFile(imageFile, imagePath); err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", "Unable to save uploaded back image")
		return
	}

	cardBack, err := h.adminService.CreateCardBack(c.Request.Context(), models.AdminCreateCardBackRequest{
		Name: models.LocalizedText{
			EN: nameEN,
			RU: nameRU,
		},
		Description: models.LocalizedText{
			EN: descriptionEN,
			RU: descriptionRU,
		},
		Tags:      tags,
		ImagePath: imagePath,
	})
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":   true,
		"card_back": cardBack,
		"message":   "Card back created successfully",
	})
}

func (h *AdminHandler) listDecks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	decks, pagination, err := h.adminService.ListDecks(c.Request.Context(), page, limit)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"decks":      decks,
		"pagination": pagination,
	})
}

func (h *AdminHandler) listCardBacks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	cardBacks, pagination, err := h.adminService.ListCardBacks(c.Request.Context(), page, limit)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"card_backs": cardBacks,
		"pagination": pagination,
	})
}

func (h *AdminHandler) getDeck(c *gin.Context) {
	deck, err := h.adminService.GetDeck(c.Request.Context(), c.Param("deckId"))
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deck":    deck,
	})
}

func (h *AdminHandler) updateDeck(c *gin.Context) {
	var req models.AdminUpdateDeckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	deck, err := h.adminService.UpdateDeckMetadata(c.Request.Context(), c.Param("deckId"), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deck":    deck,
	})
}

func (h *AdminHandler) updateUserRole(c *gin.Context) {
	actorID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	targetID, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid user id")
		return
	}

	var req models.AdminUpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, err := h.adminService.UpdateUserRole(c.Request.Context(), actorID, targetID, req.Role)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(true),
	})
}

func (h *AdminHandler) updateUserStatus(c *gin.Context) {
	actorID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	targetID, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid user id")
		return
	}

	var req models.AdminUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, err := h.adminService.UpdateUserStatus(c.Request.Context(), actorID, targetID, req.Status)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(true),
	})
}

func (h *AdminHandler) listRolePolicies(c *gin.Context) {
	policies, err := h.adminService.ListRolePolicies(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"policies": policies,
	})
}

func (h *AdminHandler) updateRolePolicy(c *gin.Context) {
	role := strings.TrimSpace(c.Param("role"))
	if role == "" {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Role is required")
		return
	}

	var req models.AdminUpdateRolePolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	policy, err := h.adminService.UpdateRolePolicy(c.Request.Context(), role, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"policy":  policy,
	})
}

func (h *AdminHandler) getSettings(c *gin.Context) {
	settings, err := h.adminService.GetSettings(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"settings": settings,
	})
}

func (h *AdminHandler) updateSettings(c *gin.Context) {
	var req models.AdminUpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	settings, err := h.adminService.UpdateSettings(c.Request.Context(), req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"settings": settings,
	})
}

func parseCSVValues(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		result = append(result, value)
	}

	return result
}
