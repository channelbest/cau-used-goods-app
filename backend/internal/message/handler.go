package message

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

func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	readStatus := c.Query("readStatus")

	items, total, err := h.service.ListByReceiver(c.Request.Context(), userID, readStatus, page, pageSize)
	if err != nil {
		if errors.Is(err, ErrInvalidReadStatus) {
			response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":    items,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || messageID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid message id")
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), userID, messageID)
	if err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			response.Error(c, http.StatusNotFound, response.CodeNotFound, "message not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, item)
}

func (h *Handler) UnreadCount(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	count, err := h.service.CountUnread(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{"count": count})
}

func (h *Handler) MarkRead(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	messageID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || messageID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid message id")
		return
	}

	if err := h.service.MarkRead(c.Request.Context(), userID, messageID); err != nil {
		if errors.Is(err, ErrMessageNotFound) {
			response.Error(c, http.StatusNotFound, response.CodeNotFound, "message not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{"read": true})
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	count, err := h.service.MarkAllRead(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"read":  true,
		"count": count,
	})
}
