package handlers

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/decks"
	"tarotdaily/backend/internal/transport/http/response"
)

type AdminDecksHandler struct {
	service *decks.Service
}

type createDeckRequest struct {
	Name string `json:"name"`
}

func NewAdminDecksHandler(service *decks.Service) *AdminDecksHandler {
	return &AdminDecksHandler{service: service}
}

func (h *AdminDecksHandler) Create(c *gin.Context) {
	var req createDeckRequest
	if !bindJSON(c, &req) {
		return
	}

	created, err := h.service.Create(c.Request.Context(), currentUserID(c), req.Name)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid deck name"},
			errorSpec{target: decks.ErrDeckNameConflict, status: http.StatusConflict, code: "deck_exists", message: "deck name already exists"},
		)
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *AdminDecksHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *AdminDecksHandler) GetByID(c *gin.Context) {
	deckID := c.Param("deckId")
	item, err := h.service.GetByID(c.Request.Context(), deckID)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid deck id"},
			errorSpec{target: decks.ErrDeckNotFound, status: http.StatusNotFound, code: "deck_not_found", message: "deck not found"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *AdminDecksHandler) UploadBack(c *gin.Context) {
	deckID := c.Param("deckId")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "invalid_input", "file is required")
		return
	}

	payload, contentType, err := readAndValidateUploadedImage(fileHeader)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrUploadTooLarge, status: http.StatusRequestEntityTooLarge, code: "file_too_large", message: "file too large"},
			errorSpec{target: decks.ErrUnsupportedImage, status: http.StatusBadRequest, code: "invalid_file_type", message: "unsupported image format"},
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid file"},
		)
		return
	}

	item, err := h.service.UploadBack(c.Request.Context(), deckID, contentType, payload)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid deck id or file"},
			errorSpec{target: decks.ErrUploadTooLarge, status: http.StatusRequestEntityTooLarge, code: "file_too_large", message: "file too large"},
			errorSpec{target: decks.ErrUnsupportedImage, status: http.StatusBadRequest, code: "invalid_file_type", message: "unsupported image format"},
			errorSpec{target: decks.ErrDeckNotFound, status: http.StatusNotFound, code: "deck_not_found", message: "deck not found"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *AdminDecksHandler) UploadCard(c *gin.Context) {
	deckID := c.Param("deckId")
	cardKey := c.Param("cardKey")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.JSONError(c, http.StatusBadRequest, "invalid_input", "file is required")
		return
	}

	payload, contentType, err := readAndValidateUploadedImage(fileHeader)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrUploadTooLarge, status: http.StatusRequestEntityTooLarge, code: "file_too_large", message: "file too large"},
			errorSpec{target: decks.ErrUnsupportedImage, status: http.StatusBadRequest, code: "invalid_file_type", message: "unsupported image format"},
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid file"},
		)
		return
	}

	item, err := h.service.UploadCard(c.Request.Context(), deckID, cardKey, contentType, payload)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: decks.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid deck id, card key or file"},
			errorSpec{target: decks.ErrUploadTooLarge, status: http.StatusRequestEntityTooLarge, code: "file_too_large", message: "file too large"},
			errorSpec{target: decks.ErrUnsupportedImage, status: http.StatusBadRequest, code: "invalid_file_type", message: "unsupported image format"},
			errorSpec{target: decks.ErrCardAlreadyExists, status: http.StatusConflict, code: "card_exists", message: "card image already uploaded"},
			errorSpec{target: decks.ErrDeckAlreadyFull, status: http.StatusConflict, code: "deck_full", message: "deck already has all required card images"},
			errorSpec{target: decks.ErrDeckNotFound, status: http.StatusNotFound, code: "deck_not_found", message: "deck not found"},
		)
		return
	}

	c.JSON(http.StatusOK, item)
}

func readAndValidateUploadedImage(fileHeader *multipart.FileHeader) ([]byte, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	payload, err := io.ReadAll(io.LimitReader(file, decks.MaxUploadBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(payload) == 0 {
		return nil, "", decks.ErrInvalidInput
	}
	if len(payload) > decks.MaxUploadBytes {
		return nil, "", decks.ErrUploadTooLarge
	}

	contentType := http.DetectContentType(payload)
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return payload, contentType, nil
	default:
		return nil, "", decks.ErrUnsupportedImage
	}
}
