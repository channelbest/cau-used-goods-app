package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type wechatLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type devLoginRequest struct {
	OpenID string `json:"openid"`
	Role   string `json:"role"`
}

func (h *Handler) DevLogin(c *gin.Context) {
	var req devLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	result, err := h.service.DevLogin(c.Request.Context(), DevLoginInput{
		OpenID: req.OpenID,
		Role:   req.Role,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *Handler) WechatLogin(c *gin.Context) {
	var req wechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "code is required")
		return
	}

	result, err := h.service.WechatLogin(c.Request.Context(), req.Code)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}
	response.Success(c, result)
}
