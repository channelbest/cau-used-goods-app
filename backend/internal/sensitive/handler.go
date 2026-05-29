package sensitive

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type createWordRequest struct {
	Word     string `json:"word" binding:"required"`
	WordType string `json:"wordType" binding:"required"`
	Status   string `json:"status"`
}

type updateWordRequest struct {
	Word     string `json:"word" binding:"required"`
	WordType string `json:"wordType" binding:"required"`
	Status   string `json:"status" binding:"required"`
}

func (h *Handler) ListWords(c *gin.Context) {
	if _, ok := middleware.CurrentUserID(c); !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	query := WordQuery{
		Status:   c.Query("status"),
		WordType: c.Query("wordType"),
		Keyword:  c.Query("keyword"),
	}
	query.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	query.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	items, total, err := h.service.ListWords(c.Request.Context(), query)
	if err != nil {
		handleSensitiveWordError(c, err)
		return
	}

	response.Success(c, gin.H{
		"items":    items,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

func (h *Handler) CreateWord(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ip := c.ClientIP()
	id, err := h.service.CreateWord(c.Request.Context(), adminID, CreateWordInput{
		Word:     req.Word,
		WordType: req.WordType,
		Status:   req.Status,
	}, &ip)
	if err != nil {
		handleSensitiveWordError(c, err)
		return
	}

	response.Success(c, gin.H{"id": id})
}

func (h *Handler) UpdateWord(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid sensitive word id")
		return
	}

	var req updateWordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ip := c.ClientIP()
	err = h.service.UpdateWord(c.Request.Context(), adminID, UpdateWordInput{
		ID:       id,
		Word:     req.Word,
		WordType: req.WordType,
		Status:   req.Status,
	}, &ip)
	if err != nil {
		handleSensitiveWordError(c, err)
		return
	}

	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) DeleteWord(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid sensitive word id")
		return
	}

	ip := c.ClientIP()
	if err := h.service.DeleteWord(c.Request.Context(), adminID, id, &ip); err != nil {
		handleSensitiveWordError(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func handleSensitiveWordError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidSensitiveWordInput):
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
	case errors.Is(err, ErrSensitiveWordNotFound):
		response.Error(c, http.StatusNotFound, response.CodeNotFound, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
	}
}
