package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *Handler, authMiddleware gin.HandlerFunc) {
	r.GET("/categories", handler.ListCategories)
	r.GET("/products", handler.ListProducts)
	r.GET("/products/:id", handler.GetProductByID)

	products := r.Group("/products")
	products.Use(authMiddleware)
	{
		products.POST("", handler.CreateProduct)
		products.GET("/my", handler.ListMyProducts)
		products.PUT("/:id", handler.UpdateProduct)
		products.PUT("/:id/status", handler.UpdateProductStatus)
		products.POST("/:id/images", handler.AddProductImages)
		products.DELETE("/:id", handler.DeleteProduct)
	}
}
