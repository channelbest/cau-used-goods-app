package upload

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

func (h *Handler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "file is required")
		return
	}

	result, err := h.service.PrepareImage(file)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	if err := c.SaveUploadedFile(file, result.SavePath); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "save file failed")
		return
	}

	response.Success(c, gin.H{
		"url": result.URL,
	})
}
