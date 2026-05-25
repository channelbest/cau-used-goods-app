package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/pkg/response"
)

func Verified(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := CurrentUserID(c)
		if !ok {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		var authStatus, accountStatus string
		err := db.QueryRowContext(
			c.Request.Context(),
			`SELECT auth_status, account_status FROM users WHERE id = ? AND is_deleted = 0 LIMIT 1`,
			userID,
		).Scan(&authStatus, &accountStatus)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not found")
			} else {
				response.Error(c, http.StatusInternalServerError, response.CodeInternal, "check user status failed")
			}
			c.Abort()
			return
		}

		if authStatus != "VERIFIED" || accountStatus != "NORMAL" {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, "student verification required")
			c.Abort()
			return
		}
		c.Next()
	}
}
