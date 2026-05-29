package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, adminMiddleware gin.HandlerFunc) {
	group := r.Group("/admin")
	group.Use(authMiddleware, adminMiddleware)
	{
		group.GET("/announcements", handler.ListAnnouncements)
		group.POST("/announcements", handler.CreateAnnouncement)
		group.PUT("/announcements/:id", handler.UpdateAnnouncement)
		group.PUT("/announcements/:id/status", handler.UpdateAnnouncementStatus)
		group.DELETE("/announcements/:id", handler.DeleteAnnouncement)

		group.GET("/logs", handler.ListLogs)
	}
}
