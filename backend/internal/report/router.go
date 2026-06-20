package report

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, verifiedMiddleware, adminMiddleware gin.HandlerFunc) {
	// 用户端
	group := r.Group("/reports")
	group.Use(authMiddleware, verifiedMiddleware)
	{
		group.POST("", handler.Create)
		group.GET("/my", handler.ListMyReports)
		group.GET("/:id", handler.GetByID)
		group.POST("/:id/close", handler.Close)
	}

	// 管理员端
	adminGroup := r.Group("/admin/reports")
	adminGroup.Use(authMiddleware, adminMiddleware)
	{
		adminGroup.GET("", handler.ListAll)
		adminGroup.GET("/:id", handler.AdminGetByID)
		adminGroup.POST("/:id/processing", handler.MarkProcessing)
		adminGroup.POST("/:id/handle", handler.Handle)
	}
}
