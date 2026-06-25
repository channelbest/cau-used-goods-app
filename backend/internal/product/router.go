package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware, optionalAuthMiddleware, verifiedMiddleware, adminMiddleware gin.HandlerFunc) {
	r.GET("/categories", handler.ListCategories)
	r.GET("/products", handler.ListProducts)
	r.GET("/products/:id", optionalAuthMiddleware, handler.GetProductByID)

	adminCategories := r.Group("/admin/categories")
	adminCategories.Use(authMiddleware, adminMiddleware)
	{
		adminCategories.GET("", handler.AdminListCategories)
		adminCategories.POST("", handler.AdminCreateCategory)
		adminCategories.PUT("/:id", handler.AdminUpdateCategory)
		adminCategories.PUT("/:id/status", handler.AdminUpdateCategoryStatus)
		adminCategories.DELETE("/:id", handler.AdminDeleteCategory)
	}

	adminProducts := r.Group("/admin/products")
	adminProducts.Use(authMiddleware, adminMiddleware)
	{
		adminProducts.GET("", handler.AdminListProducts)
		adminProducts.GET("/:id", handler.AdminGetProductByID)
		adminProducts.PUT("/:id/status", handler.AdminUpdateProductStatus)
	}

	products := r.Group("/products")
	products.Use(authMiddleware)
	{
		products.POST("", verifiedMiddleware, handler.CreateProduct)
		products.GET("/my", handler.ListMyProducts)
		products.POST("/batch-on-sale", verifiedMiddleware, handler.BatchPutOnSaleRestorable)
		products.PUT("/:id", verifiedMiddleware, handler.UpdateProduct)
		products.PUT("/:id/status", verifiedMiddleware, handler.UpdateProductStatus)
		products.POST("/:id/images", verifiedMiddleware, handler.AddProductImages)
		products.DELETE("/:id", verifiedMiddleware, handler.DeleteProduct)
		products.DELETE("/:id/images/:imageId", verifiedMiddleware, handler.DeleteProductImage)
		products.PUT("/:id/images", verifiedMiddleware, handler.ReplaceProductImages)
	}
}
