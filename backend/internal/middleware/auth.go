package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwtutil "cau-used-goods-app/backend/pkg/jwt"
	"cau-used-goods-app/backend/pkg/response"
)

const (
	ContextUserID = "userID"
	ContextRole   = "role"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid authorization header")
			c.Abort()
			return
		}

		claims, err := jwtutil.Parse(secret, parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func CurrentUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get(ContextUserID)
	if !ok {
		return 0, false
	}
	userID, ok := value.(uint64)
	return userID, ok
}

func CurrentRole(c *gin.Context) (string, bool) {
	value, ok := c.Get(ContextRole)
	if !ok {
		return "", false
	}
	role, ok := value.(string)
	return role, ok
}
