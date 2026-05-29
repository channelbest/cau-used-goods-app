package sensitive

import "github.com/gin-gonic/gin"

func RegisterAdminRoutes(r *gin.Engine, handler *Handler, authMiddleware, adminMiddleware gin.HandlerFunc) {
	group := r.Group("/admin/sensitive-words")
	group.Use(authMiddleware, adminMiddleware)
	{
		group.GET("", handler.ListWords)
		group.POST("", handler.CreateWord)
		group.PUT("/:id", handler.UpdateWord)
		group.DELETE("/:id", handler.DeleteWord)
	}
}
