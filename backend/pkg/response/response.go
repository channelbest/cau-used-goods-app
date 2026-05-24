package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess      = 0
	CodeBadRequest   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeConflict     = 409
	CodeInternal     = 500
)

type Body struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
	RequestID string      `json:"requestId,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{
		Code:      CodeSuccess,
		Message:   "success",
		Data:      data,
		Timestamp: now(),
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: now(),
	})
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
