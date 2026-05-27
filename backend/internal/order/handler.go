package order

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

type createOrderRequest struct {
	ProductID    uint64  `json:"productId" binding:"required"`
	Remark       *string `json:"remark"`
	MeetTime     *string `json:"meetTime"`
	MeetLocation *string `json:"meetLocation"`
}

type cancelOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	order, err := h.service.Create(c.Request.Context(), CreateOrderInput{
		ProductID:    req.ProductID,
		BuyerID:      userID,
		Remark:       req.Remark,
		MeetTime:     req.MeetTime,
		MeetLocation: req.MeetLocation,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *Handler) Confirm(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid order id")
		return
	}

	order, err := h.service.Confirm(c.Request.Context(), ConfirmOrderInput{
		OrderID:  orderID,
		SellerID: userID,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *Handler) Cancel(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid order id")
		return
	}

	var req cancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "cancel reason is required")
		return
	}

	order, err := h.service.Cancel(c.Request.Context(), CancelOrderInput{
		OrderID: orderID,
		UserID:  userID,
		Reason:  req.Reason,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *Handler) Complete(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid order id")
		return
	}

	order, err := h.service.Complete(c.Request.Context(), CompleteOrderInput{
		OrderID:  orderID,
		SellerID: userID,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *Handler) ExceptionClose(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid order id")
		return
	}

	var req cancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "reason is required")
		return
	}

	order, err := h.service.ExceptionClose(c.Request.Context(), ExceptionCloseOrderInput{
		OrderID: orderID,
		UserID:  userID,
		Reason:  req.Reason,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *Handler) GetByID(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || orderID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid order id")
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), orderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	if order == nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "order not found")
		return
	}
	if order.BuyerID != userID && order.SellerID != userID {
		response.Error(c, http.StatusForbidden, response.CodeForbidden, "permission denied")
		return
	}
	response.Success(c, order)
}

func (h *Handler) ListMyOrders(c *gin.Context) {
	userID, ok := middleware.CurrentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	role := c.DefaultQuery("role", "buyer")
	status := c.DefaultQuery("status", "")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	var orders []OrderDetail
	var total int
	var err error

	if role == "seller" {
		orders, total, err = h.service.ListBySeller(c.Request.Context(), userID, status, page, pageSize)
	} else {
		orders, total, err = h.service.ListByBuyer(c.Request.Context(), userID, status, page, pageSize)
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":    orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) CancelExpired(c *gin.Context) {
	count, err := h.service.CancelExpiredOrders(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, gin.H{"cancelledCount": count})
}
