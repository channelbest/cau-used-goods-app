package favorite

import (
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

type addFavoriteRequest struct {
	ProductID uint64 `json:"productId" binding:"required"`
}

func (h *Handler) Add(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req addFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "productId is required")
		return
	}

	if err := h.service.Add(c.Request.Context(), userID, req.ProductID); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"favorited": true})
}

func (h *Handler) Remove(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	productID, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	if err := h.service.Remove(c.Request.Context(), userID, productID); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"favorited": false})
}

func (h *Handler) Check(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	productID, err := strconv.ParseUint(c.Query("productId"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	favorited, err := h.service.IsFavorited(c.Request.Context(), userID, productID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"favorited": favorited})
}

func (h *Handler) List(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.service.ListByUser(c.Request.Context(), userID, page, pageSize)
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
