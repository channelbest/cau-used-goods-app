package review

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

type createReviewRequest struct {
	OrderID uint64  `json:"orderId" binding:"required"`
	Rating  int     `json:"rating" binding:"required,min=1,max=5"`
	Content *string `json:"content"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	review, err := h.service.Create(c.Request.Context(), CreateReviewInput{
		OrderID:    req.OrderID,
		ReviewerID: userID,
		Rating:     req.Rating,
		Content:    req.Content,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, review)
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid review id")
		return
	}

	review, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	if review == nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "review not found")
		return
	}
	response.Success(c, review)
}

func (h *Handler) ListByProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.service.ListByProduct(c.Request.Context(), productID, page, pageSize)
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

func (h *Handler) ListBySeller(c *gin.Context) {
	sellerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || sellerID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid seller id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	items, total, err := h.service.ListBySeller(c.Request.Context(), sellerID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	avgRating, reviewCount, _ := h.service.GetSellerRating(c.Request.Context(), sellerID)

	response.Success(c, gin.H{
		"items":       items,
		"total":       total,
		"page":        page,
		"pageSize":    pageSize,
		"avgRating":   avgRating,
		"reviewCount": reviewCount,
	})
}
