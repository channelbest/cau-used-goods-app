package report

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

type createReportRequest struct {
	TargetType  string   `json:"targetType" binding:"required,oneof=PRODUCT USER ORDER"`
	TargetID    uint64   `json:"targetId" binding:"required"`
	ReasonType  string   `json:"reasonType" binding:"required"`
	Description *string  `json:"description"`
	Images      []string `json:"images"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	report, err := h.service.Create(c.Request.Context(), CreateReportInput{
		ReporterID:  userID,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		ReasonType:  req.ReasonType,
		Description: req.Description,
		Images:      req.Images,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, report)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid report id")
		return
	}

	report, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	if report == nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "report not found")
		return
	}
	if report.ReporterID != userID {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "permission denied")
		return
	}
	response.Success(c, report)
}

func (h *Handler) ListMyReports(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.service.ListByReporter(c.Request.Context(), userID, page, pageSize)
	if err != nil {
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

func (h *Handler) ListAll(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.service.ListAll(c.Request.Context(), status, page, pageSize)
	if err != nil {
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

type handleReportRequest struct {
	Status        string  `json:"status" binding:"required,oneof=APPROVED REJECTED"`
	HandleResult  *string `json:"handleResult"`
	AccountStatus *string `json:"accountStatus"`
}

type closeReportRequest struct {
	CloseReason *string `json:"closeReason"`
}

func (h *Handler) Handle(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	reportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || reportID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid report id")
		return
	}

	var req handleReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}
	if req.AccountStatus != nil && strings.TrimSpace(*req.AccountStatus) != "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "accountStatus is deprecated; use the user status API with relatedType and relatedId")
		return
	}

	ipAddress := c.ClientIP()
	report, err := h.service.Handle(c.Request.Context(), HandleReportInput{
		ReportID:     reportID,
		HandlerID:    adminID,
		Status:       req.Status,
		HandleResult: req.HandleResult,
		IPAddress:    &ipAddress,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, report)
}

func (h *Handler) AdminGetByID(c *gin.Context) {
	if _, ok := middleware.CurrentUserID(c); !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid report id")
		return
	}

	report, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	if report == nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "report not found")
		return
	}
	response.Success(c, report)
}

func (h *Handler) Close(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	reportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || reportID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid report id")
		return
	}

	var req closeReportRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	report, err := h.service.Close(c.Request.Context(), CloseReportInput{
		ReportID:    reportID,
		ReporterID:  userID,
		CloseReason: req.CloseReason,
	})
	if err != nil {
		switch err.Error() {
		case "permission denied":
			response.Error(c, http.StatusForbidden, response.CodeForbidden, err.Error())
		case "report not found":
			response.Error(c, http.StatusNotFound, response.CodeNotFound, err.Error())
		default:
			response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		}
		return
	}
	response.Success(c, report)
}

func (h *Handler) MarkProcessing(c *gin.Context) {
	adminID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	reportID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || reportID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid report id")
		return
	}

	ipAddress := c.ClientIP()
	report, err := h.service.MarkProcessing(c.Request.Context(), MarkReportProcessingInput{
		ReportID:  reportID,
		HandlerID: adminID,
		IPAddress: &ipAddress,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, report)
}
