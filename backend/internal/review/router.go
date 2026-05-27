package review

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, verifiedMiddleware gin.HandlerFunc) {
	// 公开接口
	r.GET("/products/:id/reviews", handler.ListByProduct)
	r.GET("/sellers/:id/reviews", handler.ListBySeller)
	r.GET("/reviews/:id", handler.GetByID)

	// 需要登录+认证
	group := r.Group("/reviews")
	group.Use(authMiddleware, verifiedMiddleware)
	{
		group.POST("", handler.Create)
	}
}
