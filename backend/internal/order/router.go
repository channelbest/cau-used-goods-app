package order

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, verifiedMiddleware gin.HandlerFunc) {
	group := r.Group("/orders")
	group.Use(authMiddleware, verifiedMiddleware)
	{
		group.POST("", handler.Create)
		group.GET("", handler.ListMyOrders)
		group.GET("/:id", handler.GetByID)
		group.POST("/:id/confirm", handler.Confirm)
		group.POST("/:id/cancel", handler.Cancel)
		group.POST("/:id/complete", handler.Complete)
		group.POST("/:id/exception-close", handler.ExceptionClose)
	}

	// 管理员清理超时订单
	adminGroup := r.Group("/admin/orders")
	adminGroup.Use(authMiddleware)
	{
		adminGroup.POST("/cleanup-expired", handler.CancelExpired)
	}
}
