package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc, enableDevLogin bool) {
	group := r.Group("/auth")
	if enableDevLogin {
		group.POST("/dev-login", handler.DevLogin)
	}
	group.POST("/wechat-login", handler.WechatLogin)
}
