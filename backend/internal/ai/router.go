package ai

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := r.Group("/ai")
	group.Use(authMiddleware)
	{
		group.POST("/optimize-product", handler.OptimizeProduct)
	}
}
