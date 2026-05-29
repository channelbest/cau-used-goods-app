package admin

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

type createAnnouncementRequest struct {
	Title    string  `json:"title" binding:"required"`
	Content  *string `json:"content"`
	CoverURL *string `json:"coverUrl"`
	Status   string  `json:"status"`
}

type updateAnnouncementRequest struct {
	Title    string  `json:"title" binding:"required"`
	Content  *string `json:"content"`
	CoverURL *string `json:"coverUrl"`
}

type updateAnnouncementStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) CreateAnnouncement(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ip := c.ClientIP()
	id, err := h.service.CreateAnnouncement(c.Request.Context(), CreateAnnouncementInput{
		AdminID:   adminID,
		Title:     req.Title,
		Content:   req.Content,
		CoverURL:  req.CoverURL,
		Status:    req.Status,
		IPAddress: &ip,
	})
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	response.Success(c, gin.H{"id": id})
}

func (h *Handler) ListAnnouncements(c *gin.Context) {
	if _, ok := middleware.CurrentUserID(c); !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	query := AnnouncementQuery{
		Status:  c.Query("status"),
		Keyword: c.Query("keyword"),
	}
	query.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	query.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	items, total, err := h.service.ListAnnouncements(c.Request.Context(), query)
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	response.Success(c, gin.H{
		"items":    items,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

func (h *Handler) UpdateAnnouncement(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid announcement id")
		return
	}

	var req updateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ip := c.ClientIP()
	err = h.service.UpdateAnnouncement(c.Request.Context(), UpdateAnnouncementInput{
		AdminID:   adminID,
		ID:        id,
		Title:     req.Title,
		Content:   req.Content,
		CoverURL:  req.CoverURL,
		IPAddress: &ip,
	})
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) UpdateAnnouncementStatus(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid announcement id")
		return
	}

	var req updateAnnouncementStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ip := c.ClientIP()
	err = h.service.UpdateAnnouncementStatus(c.Request.Context(), UpdateAnnouncementStatusInput{
		AdminID:   adminID,
		ID:        id,
		Status:    req.Status,
		IPAddress: &ip,
	})
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	response.Success(c, gin.H{"updated": true})
}

func (h *Handler) DeleteAnnouncement(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid announcement id")
		return
	}

	ip := c.ClientIP()
	if err := h.service.DeleteAnnouncement(c.Request.Context(), adminID, id, &ip); err != nil {
		handleAnnouncementError(c, err)
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) ListLogs(c *gin.Context) {
	if _, ok := middleware.CurrentUserID(c); !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	query := LogQuery{
		OperationType: c.Query("operationType"),
		TargetType:    c.Query("targetType"),
		StartTime:     c.Query("startTime"),
		EndTime:       c.Query("endTime"),
	}
	query.AdminID, _ = strconv.ParseUint(c.Query("adminId"), 10, 64)
	query.TargetID, _ = strconv.ParseUint(c.Query("targetId"), 10, 64)
	query.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	query.PageSize, _ = strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	items, total, err := h.service.ListLogs(c.Request.Context(), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":    items,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

func handleAnnouncementError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidAnnouncementInput):
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
	case errors.Is(err, ErrAnnouncementNotFound):
		response.Error(c, http.StatusNotFound, response.CodeNotFound, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
	}
}
