package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, list)
}

type createProductRequest struct {
	CategoryID     uint64   `json:"categoryId" binding:"required"`
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	OriginalPrice  *float64 `json:"originalPrice"`
	Price          float64  `json:"price" binding:"required"`
	ConditionLevel string   `json:"conditionLevel"`
	MeetLocation   string   `json:"meetLocation"`
}

type createCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  uint64 `json:"parentId"`
	SortOrder int    `json:"sortOrder"`
	Status    string `json:"status"`
}

type updateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  uint64 `json:"parentId"`
	SortOrder int    `json:"sortOrder"`
	Status    string `json:"status"`
}

type updateCategoryStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) AdminListCategories(c *gin.Context) {
	status := c.Query("status")
	list, err := h.service.ListAllCategories(c.Request.Context(), status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *Handler) AdminCreateCategory(c *gin.Context) {
	adminID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	id, err := h.service.CreateCategory(c.Request.Context(), CategoryCreateInput{
		AdminID:   adminID,
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		IPAddress: &ipAddress,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *Handler) AdminUpdateCategory(c *gin.Context) {
	adminID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid category id")
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	if err := h.service.UpdateCategory(c.Request.Context(), CategoryUpdateInput{
		AdminID:   adminID,
		ID:        id,
		Name:      req.Name,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Status:    req.Status,
		IPAddress: &ipAddress,
	}); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *Handler) AdminUpdateCategoryStatus(c *gin.Context) {
	adminID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid category id")
		return
	}

	var req updateCategoryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	if err := h.service.UpdateCategoryStatus(c.Request.Context(), adminID, id, req.Status, &ipAddress, ""); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":     id,
		"status": req.Status,
	})
}

func (h *Handler) AdminDeleteCategory(c *gin.Context) {
	adminID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid category id")
		return
	}

	ipAddress := c.ClientIP()
	if err := h.service.UpdateCategoryStatus(c.Request.Context(), adminID, id, "DISABLED", &ipAddress, "DELETE_CATEGORY"); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":     id,
		"status": "DISABLED",
	})
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req createProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	productID, err := h.service.CreateProduct(c.Request.Context(), ProductCreateInput{
		SellerID:       userID,
		CategoryID:     req.CategoryID,
		Title:          req.Title,
		Description:    req.Description,
		OriginalPrice:  req.OriginalPrice,
		Price:          req.Price,
		ConditionLevel: req.ConditionLevel,
		MeetLocation:   req.MeetLocation,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id": productID,
	})
}

func (h *Handler) ListProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	conditionLevel := c.Query("conditionLevel")
	status := c.DefaultQuery("status", "ON_SALE")
	sort := c.DefaultQuery("sort", "newest")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var categoryID uint64
	if v := c.Query("categoryId"); v != "" {
		parsed, _ := strconv.ParseUint(v, 10, 64)
		categoryID = parsed
	}

	var minPrice *float64
	if v := c.Query("minPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			minPrice = &parsed
		}
	}

	var maxPrice *float64
	if v := c.Query("maxPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			maxPrice = &parsed
		}
	}

	if !isValidProductStatus(status) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid status")
		return
	}

	if !isValidProductSort(sort) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid sort")
		return
	}

	result, err := h.service.ListProducts(c.Request.Context(), ProductListInput{
		Keyword:        keyword,
		CategoryID:     categoryID,
		ConditionLevel: conditionLevel,
		Status:         status,
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		Sort:           sort,
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) AdminListProducts(c *gin.Context) {
	keyword := c.Query("keyword")
	conditionLevel := c.Query("conditionLevel")
	status := c.Query("status")
	sort := c.DefaultQuery("sort", "newest")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	var categoryID uint64
	if v := c.Query("categoryId"); v != "" {
		parsed, _ := strconv.ParseUint(v, 10, 64)
		categoryID = parsed
	}

	var minPrice *float64
	if v := c.Query("minPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			minPrice = &parsed
		}
	}

	var maxPrice *float64
	if v := c.Query("maxPrice"); v != "" {
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			maxPrice = &parsed
		}
	}

	if status != "" && !isValidAdminProductStatus(status) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid status")
		return
	}

	if !isValidProductSort(sort) {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid sort")
		return
	}

	result, err := h.service.ListAdminProducts(c.Request.Context(), ProductListInput{
		Keyword:        keyword,
		CategoryID:     categoryID,
		ConditionLevel: conditionLevel,
		Status:         status,
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		Sort:           sort,
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *Handler) AdminGetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	product, err := h.service.AdminGetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "product not found")
		return
	}

	response.Success(c, product)
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	userID, _ := middleware.CurrentUserID(c)
	role, _ := middleware.CurrentRole(c)
	product, err := h.service.GetProductByID(c.Request.Context(), id, ProductViewer{
		UserID: userID,
		Role:   role,
	})
	if err != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, "product not available")
		return
	}

	response.Success(c, product)
}

func (h *Handler) ListMyProducts(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	list, err := h.service.ListMyProducts(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, err.Error())
		return
	}

	response.Success(c, list)
}

type updateProductRequest struct {
	CategoryID     int64   `json:"categoryId"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	OriginalPrice  float64 `json:"originalPrice"`
	Price          float64 `json:"price"`
	ConditionLevel string  `json:"conditionLevel"`
	MeetLocation   string  `json:"meetLocation"`
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	if req.CategoryID <= 0 || req.Title == "" || req.Price <= 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "categoryId, title and price are required")
		return
	}

	err = h.service.UpdateProduct(c.Request.Context(), ProductUpdateInput{
		ProductID:      productID,
		SellerID:       userID,
		CategoryID:     uint64(req.CategoryID),
		Title:          req.Title,
		Description:    req.Description,
		OriginalPrice:  &req.OriginalPrice,
		Price:          req.Price,
		ConditionLevel: req.ConditionLevel,
		MeetLocation:   req.MeetLocation,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id": productID,
	})
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), productID, userID); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "product not found or cannot delete")
		return
	}

	response.Success(c, gin.H{
		"id": productID,
	})
}

type updateProductStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason"`
}

type adminUpdateProductStatusRequest struct {
	Status      string `json:"status" binding:"required"`
	Reason      string `json:"reason"`
	RelatedType string `json:"relatedType"`
	RelatedID   uint64 `json:"relatedId"`
}

func (h *Handler) UpdateProductStatus(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	var req updateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if req.Status != "ON_SALE" && req.Status != "OFF_SHELF" {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "status must be ON_SALE or OFF_SHELF")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	if err := h.service.UpdateProductStatus(c.Request.Context(), productID, userID, req.Status, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":     productID,
		"status": req.Status,
	})
}

func (h *Handler) AdminUpdateProductStatus(c *gin.Context) {
	adminID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "unauthorized")
		return
	}

	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	var req adminUpdateProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	ipAddress := c.ClientIP()
	if err := h.service.AdminUpdateProductStatus(c.Request.Context(), AdminUpdateProductStatusInput{
		AdminID:     adminID,
		ProductID:   productID,
		Status:      req.Status,
		Reason:      req.Reason,
		IPAddress:   &ipAddress,
		RelatedType: req.RelatedType,
		RelatedID:   req.RelatedID,
	}); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":     productID,
		"status": req.Status,
	})
}

type addProductImagesRequest struct {
	Images []string `json:"images" binding:"required"`
}

func (h *Handler) AddProductImages(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	var req addProductImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if len(req.Images) == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "images is required")
		return
	}

	if len(req.Images) > 9 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "最多只能上传9张图片")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	err = h.service.AddProductImages(c.Request.Context(), ProductImagesInput{
		ProductID: productID,
		SellerID:  userID,
		Images:    req.Images,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":     productID,
		"images": req.Images,
	})
}

func (h *Handler) DeleteProductImage(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)
	if err != nil || imageID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid image id")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	err = h.service.DeleteProductImage(c.Request.Context(), ProductImageDeleteInput{
		ProductID: productID,
		SellerID:  userID,
		ImageID:   imageID,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":      productID,
		"imageId": imageID,
	})
}

func (h *Handler) ReplaceProductImages(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || productID == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid product id")
		return
	}

	var req addProductImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "invalid request body")
		return
	}

	if len(req.Images) > 9 {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, "最多只能上传9张图片")
		return
	}

	userID, ok := currentUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user not login")
		return
	}

	err = h.service.ReplaceProductImages(c.Request.Context(), ProductImagesInput{
		ProductID: productID,
		SellerID:  userID,
		Images:    req.Images,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":     productID,
		"images": req.Images,
	})
}

func currentUserID(c *gin.Context) (uint64, bool) {
	value, ok := c.Get(middleware.ContextUserID)
	if !ok {
		return 0, false
	}

	userID, ok := value.(uint64)
	if !ok {
		return 0, false
	}

	return userID, true
}

func isValidProductStatus(status string) bool {
	switch status {
	case "ON_SALE":
		return true
	default:
		return false
	}
}

func isValidProductSort(sort string) bool {
	switch sort {
	case "newest", "price_asc", "price_desc", "popular":
		return true
	default:
		return false
	}
}
