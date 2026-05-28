package favorite

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, verifiedMiddleware gin.HandlerFunc) {
	group := r.Group("/favorites")
	group.Use(authMiddleware, verifiedMiddleware)
	{
		group.POST("", handler.Add)
		group.GET("", handler.List)
		group.GET("/check", handler.Check)
		group.DELETE("/:productId", handler.Remove)
	}
}
