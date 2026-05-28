package ai

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

func (h *Handler) OptimizeProduct(c *gin.Context) {
	var req OptimizeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if req.Title == "" && req.Description == "" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "title or description is required")
		return
	}

	result, err := h.service.OptimizeProduct(c.Request.Context(), req)
	if err != nil {
		// 调试阶段先返回真实错误，方便排查 API key、模型名、JSON 解析等问题。
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}
