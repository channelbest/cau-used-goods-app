package upload

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := r.Group("/upload")
	group.Use(authMiddleware)
	{
		group.POST("/image", handler.UploadImage)
	}
}
