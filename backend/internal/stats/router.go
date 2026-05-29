package stats

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, adminMiddleware gin.HandlerFunc) {
	group := r.Group("/stats")
	group.Use(authMiddleware, adminMiddleware)
	{
		group.GET("/products/overview", handler.ProductOverview)
		group.GET("/products/category-distribution", handler.CategoryDistribution)
		group.GET("/products/status-distribution", handler.StatusDistribution)
		group.GET("/products/trend", handler.ProductTrend)
		group.GET("/orders/overview", handler.OrderOverview)
		group.GET("/users/overview", handler.UserOverview)
		group.GET("/reports/overview", handler.ReportOverview)
	}
}
