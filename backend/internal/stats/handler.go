package stats

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ProductOverview(c *gin.Context) {
	result, err := h.service.ProductOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) CategoryDistribution(c *gin.Context) {
	result, err := h.service.CategoryDistribution(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) StatusDistribution(c *gin.Context) {
	result, err := h.service.StatusDistribution(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) ProductTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	result, err := h.service.ProductTrend(c.Request.Context(), days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"days": days,
		"list": result,
	})
}

func (h *Handler) OrderOverview(c *gin.Context) {
	result, err := h.service.OrderOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) UserOverview(c *gin.Context) {
	result, err := h.service.UserOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) ReportOverview(c *gin.Context) {
	result, err := h.service.ReportOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}
