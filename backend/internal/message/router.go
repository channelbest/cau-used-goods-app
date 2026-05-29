package message

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := r.Group("/messages")
	group.Use(authMiddleware)
	{
		group.GET("", handler.List)
		group.GET("/unread-count", handler.UnreadCount)
		group.GET("/:id", handler.GetByID)
		group.PUT("/read-all", handler.MarkAllRead)
		group.PUT("/:id/read", handler.MarkRead)
	}
}
