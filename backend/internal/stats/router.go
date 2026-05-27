package stats

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := r.Group("/stats")
	group.Use(authMiddleware)
	{
		group.GET("/products/overview", handler.ProductOverview)
		group.GET("/products/category-distribution", handler.CategoryDistribution)
		group.GET("/products/status-distribution", handler.StatusDistribution)
		group.GET("/products/trend", handler.ProductTrend)
	}
}
