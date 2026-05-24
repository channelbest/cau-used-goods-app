package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/pkg/response"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := CurrentRole(c)
		if !ok || role != "ADMIN" {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, "admin permission required")
			c.Abort()
			return
		}
		c.Next()
	}
}
