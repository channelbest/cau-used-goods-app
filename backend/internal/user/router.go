package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	group := r.Group("/users")
	group.Use(authMiddleware)

	group.GET("/me", handler.Me)
	group.PUT("/profile", handler.UpdateProfile)
	group.POST("/avatar", handler.UploadAvatar)
	group.POST("/student-verify", handler.SubmitStudentVerification)
	group.GET("/student-verify", handler.StudentVerification)
}

func RegisterAdminRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {
	group := r.Group("/admin/users")
	group.Use(authMiddleware, adminMiddleware)

	group.GET("/student-verifications", handler.ListStudentVerifications)
	group.PUT("/:id/student-verify", handler.ReviewStudentVerification)
}
