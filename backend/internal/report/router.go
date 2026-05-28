package report

import (
	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, verifiedMiddleware gin.HandlerFunc) {
	// 用户端
	group := r.Group("/reports")
	group.Use(authMiddleware, verifiedMiddleware)
	{
		group.POST("", handler.Create)
		group.GET("/my", handler.ListMyReports)
		group.GET("/:id", handler.GetByID)
	}

	// 管理员端
	adminGroup := r.Group("/admin/reports")
	adminGroup.Use(authMiddleware, middleware.Admin())
	{
		adminGroup.GET("", handler.ListAll)
		adminGroup.POST("/:id/handle", handler.Handle)
	}
}
