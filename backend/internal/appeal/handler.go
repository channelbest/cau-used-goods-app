package appeal

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

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

type createAppealRequest struct {
	TargetType   string   `json:"targetType" binding:"required"`
	TargetID     uint64   `json:"targetId" binding:"required"`
	Reason       string   `json:"reason" binding:"required"`
	EvidenceURLs []string `json:"evidenceUrls"`
}

type handleAppealRequest struct {
	Status        string `json:"status" binding:"required"`
	HandleResult  string `json:"handleResult" binding:"required"`
	AccountStatus string `json:"accountStatus"`
}

type closeAppealRequest struct {
	CloseReason string `json:"closeReason"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	item, err := h.service.Create(c.Request.Context(), CreateAppealInput{
		AppellantID:  userID,
		TargetType:   req.TargetType,
		TargetID:     req.TargetID,
		Reason:       req.Reason,
		EvidenceURLs: req.EvidenceURLs,
	})
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) ListMy(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	query := parseQuery(c)
	items, total, err := h.service.ListMy(c.Request.Context(), userID, query)
	if err != nil {
		writeAppealError(c, err)
		return
	}
	writeList(c, items, total, query)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	appealID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid appeal id")
		return
	}
	item, err := h.service.GetByID(c.Request.Context(), appealID, userID, false)
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) Close(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	appealID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid appeal id")
		return
	}

	var req closeAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	item, err := h.service.Close(c.Request.Context(), CloseAppealInput{
		AppealID:    appealID,
		AppellantID: userID,
		CloseReason: req.CloseReason,
	})
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) AdminList(c *gin.Context) {
	query := parseQuery(c)
	items, total, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		writeAppealError(c, err)
		return
	}
	writeList(c, items, total, query)
}

func (h *Handler) AdminGetByID(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	appealID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid appeal id")
		return
	}
	item, err := h.service.GetByID(c.Request.Context(), appealID, adminID, true)
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) AdminHandle(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	appealID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid appeal id")
		return
	}

	var req handleAppealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.AccountStatus) != "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "accountStatus is deprecated; use the user status API with relatedType and relatedId")
		return
	}
	ipAddress := c.ClientIP()
	item, err := h.service.Handle(c.Request.Context(), HandleAppealInput{
		AppealID:     appealID,
		AdminID:      adminID,
		Status:       req.Status,
		HandleResult: req.HandleResult,
		IPAddress:    &ipAddress,
	})
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) AdminMarkProcessing(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	appealID, err := parseID(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid appeal id")
		return
	}

	ipAddress := c.ClientIP()
	item, err := h.service.MarkProcessing(c.Request.Context(), appealID, adminID, &ipAddress)
	if err != nil {
		writeAppealError(c, err)
		return
	}
	response.Success(c, item)
}

func parseID(c *gin.Context) (uint64, error) {
	return strconv.ParseUint(c.Param("id"), 10, 64)
}

func parseQuery(c *gin.Context) AppealQuery {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	targetID, _ := strconv.ParseUint(c.DefaultQuery("targetId", "0"), 10, 64)
	return AppealQuery{
		TargetType: c.Query("targetType"),
		TargetID:   targetID,
		Status:     c.Query("status"),
		Page:       page,
		PageSize:   pageSize,
	}
}

func writeList(c *gin.Context, items []AppealDetail, total int, query AppealQuery) {
	response.Success(c, gin.H{
		"items":    items,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

func writeAppealError(c *gin.Context, err error) {
	message := err.Error()
	switch message {
	case "appeal not found", "appeal target not found":
		response.Error(c, http.StatusNotFound, response.CodeNotFound, message)
	case "permission denied":
		response.Error(c, http.StatusForbidden, response.CodeForbidden, message)
	case "you already have an active appeal for this target", "appeal already handled":
		response.Error(c, http.StatusConflict, response.CodeConflict, message)
	case "appeal must be marked processing before handle", "user appeal target must be restored before approval",
		"report result can only be appealed after handled":
		response.Error(c, http.StatusConflict, response.CodeConflict, message)
	case "invalid targetType", "invalid status", "appellantId is required", "targetType must be PRODUCT, USER, ORDER or REPORT",
		"targetId is required", "reason is required", "reason cannot exceed 500 characters",
		"evidenceUrls cannot exceed 9", "appealId is required", "adminId is required",
		"status must be APPROVED or REJECTED", "handleResult is required",
		"handleResult cannot exceed 500 characters", "appeal cannot be marked processing", "appeal cannot be closed",
		"closeReason cannot exceed 500 characters":
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, message)
	default:
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, message)
	}
}
