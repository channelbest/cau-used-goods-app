package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"cau-used-goods-app/backend/internal/admin"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/message"
	"cau-used-goods-app/backend/internal/sensitive"
)

type Service struct {
	repo             *Repository
	sensitiveService *sensitive.Service
	adminLogger      *admin.Service
	messageService   *message.Service
}

func NewService(repo *Repository, sensitiveService *sensitive.Service, adminLogger *admin.Service, messageService *message.Service) *Service {
	return &Service{
		repo:             repo,
		sensitiveService: sensitiveService,
		adminLogger:      adminLogger,
		messageService:   messageService,
	}
}

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) ListAllCategories(ctx context.Context, status string) ([]Category, error) {
	status = strings.TrimSpace(status)
	if status != "" && !isValidCategoryStatus(status) {
		return nil, fmt.Errorf("invalid category status")
	}
	return s.repo.ListAllCategories(ctx, status)
}

type CategoryCreateInput struct {
	AdminID   uint64
	Name      string
	ParentID  uint64
	SortOrder int
	Status    string
	IPAddress *string
}

func (s *Service) CreateCategory(ctx context.Context, input CategoryCreateInput) (uint64, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Status = strings.TrimSpace(input.Status)
	if input.IPAddress != nil {
		trimmed := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &trimmed
	}
	if input.Status == "" {
		input.Status = "ENABLED"
	}
	if err := validateCategoryInput(input.Name, input.Status); err != nil {
		return 0, err
	}
	var id uint64
	description := fmt.Sprintf("create category: %s", input.Name)
	err := db.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		id, err = s.repo.CreateCategoryTx(ctx, tx, CreateCategoryInput{
			Name:      input.Name,
			ParentID:  input.ParentID,
			SortOrder: input.SortOrder,
			Status:    input.Status,
		})
		if err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, input.AdminID, admin.OperationCreateCategory, admin.TargetTypeCategory, id, description, input.IPAddress, nil, nil)
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

type CategoryUpdateInput struct {
	AdminID   uint64
	ID        uint64
	Name      string
	ParentID  uint64
	SortOrder int
	Status    string
	IPAddress *string
}

func (s *Service) UpdateCategory(ctx context.Context, input CategoryUpdateInput) error {
	input.Name = strings.TrimSpace(input.Name)
	input.Status = strings.TrimSpace(input.Status)
	if input.IPAddress != nil {
		trimmed := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &trimmed
	}
	if input.Status == "" {
		input.Status = "ENABLED"
	}
	if input.ID == 0 {
		return fmt.Errorf("category id is required")
	}
	if err := validateCategoryInput(input.Name, input.Status); err != nil {
		return err
	}
	if input.ParentID == input.ID {
		return fmt.Errorf("category parent cannot be itself")
	}
	description := fmt.Sprintf("update category: %s", input.Name)
	return db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.UpdateCategoryTx(ctx, tx, UpdateCategoryInput{
			ID:        input.ID,
			Name:      input.Name,
			ParentID:  input.ParentID,
			SortOrder: input.SortOrder,
			Status:    input.Status,
		}); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, input.AdminID, admin.OperationUpdateCategory, admin.TargetTypeCategory, input.ID, description, input.IPAddress, nil, nil)
	})
}

func (s *Service) UpdateCategoryStatus(ctx context.Context, adminID, id uint64, status string, ipAddress *string, operationType string) error {
	status = strings.TrimSpace(status)
	if ipAddress != nil {
		trimmed := strings.TrimSpace(*ipAddress)
		ipAddress = &trimmed
	}
	if id == 0 {
		return fmt.Errorf("category id is required")
	}
	if !isValidCategoryStatus(status) {
		return fmt.Errorf("invalid category status")
	}
	if operationType == "" {
		operationType = admin.OperationStatusCategory
	}
	description := fmt.Sprintf("update category status to %s", status)
	return db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.UpdateCategoryStatusTx(ctx, tx, id, status); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, adminID, operationType, admin.TargetTypeCategory, id, description, ipAddress, nil, nil)
	})
}

func validateCategoryInput(name string, status string) error {
	if name == "" {
		return fmt.Errorf("category name is required")
	}
	if len([]rune(name)) > 50 {
		return fmt.Errorf("category name cannot exceed 50 characters")
	}
	if !isValidCategoryStatus(status) {
		return fmt.Errorf("invalid category status")
	}
	return nil
}

func isValidCategoryStatus(status string) bool {
	return status == "ENABLED" || status == "DISABLED"
}

type ProductCreateInput struct {
	SellerID       uint64
	CategoryID     uint64
	Title          string
	Description    string
	OriginalPrice  *float64
	Price          float64
	ConditionLevel string
	MeetLocation   string
}

func (s *Service) CreateProduct(ctx context.Context, input ProductCreateInput) (uint64, error) {
	if err := s.checkSensitive(ctx, input.Title, input.Description); err != nil {
		return 0, err
	}

	return s.repo.CreateProduct(ctx, CreateProductInput{
		SellerID:       input.SellerID,
		CategoryID:     input.CategoryID,
		Title:          input.Title,
		Description:    input.Description,
		OriginalPrice:  input.OriginalPrice,
		Price:          input.Price,
		ConditionLevel: input.ConditionLevel,
		MeetLocation:   input.MeetLocation,
	})
}

type ProductListInput struct {
	Keyword        string
	CategoryID     uint64
	ConditionLevel string
	Status         string
	MinPrice       *float64
	MaxPrice       *float64
	Sort           string
	Page           int
	PageSize       int
}

func (s *Service) ListProducts(ctx context.Context, input ProductListInput) (*ProductListResult, error) {
	return s.repo.ListProducts(ctx, ListProductsInput{
		Keyword:        input.Keyword,
		CategoryID:     input.CategoryID,
		ConditionLevel: input.ConditionLevel,
		Status:         input.Status,
		MinPrice:       input.MinPrice,
		MaxPrice:       input.MaxPrice,
		Sort:           input.Sort,
		Page:           input.Page,
		PageSize:       input.PageSize,
	})
}

func (s *Service) ListAdminProducts(ctx context.Context, input ProductListInput) (*ProductListResult, error) {
	return s.repo.ListAdminProducts(ctx, ListProductsInput{
		Keyword:        input.Keyword,
		CategoryID:     input.CategoryID,
		ConditionLevel: input.ConditionLevel,
		Status:         input.Status,
		MinPrice:       input.MinPrice,
		MaxPrice:       input.MaxPrice,
		Sort:           input.Sort,
		Page:           input.Page,
		PageSize:       input.PageSize,
	})
}

type ProductViewer struct {
	UserID uint64
	Role   string
}

func (s *Service) GetProductByID(ctx context.Context, id uint64, viewer ProductViewer) (*Product, error) {
	if err := s.repo.IncrementViewCount(ctx, id, viewer); err != nil {
		return nil, err
	}
	return s.repo.GetProductByID(ctx, id, viewer)
}

func (s *Service) AdminGetProductByID(ctx context.Context, id uint64) (*Product, error) {
	if id == 0 {
		return nil, fmt.Errorf("productId is required")
	}
	return s.repo.AdminGetProductByID(ctx, id)
}

func (s *Service) ListMyProducts(ctx context.Context, sellerID uint64) ([]Product, error) {
	return s.repo.ListMyProducts(ctx, sellerID)
}

func (s *Service) DeleteProduct(ctx context.Context, productID uint64, sellerID uint64) error {
	return s.repo.DeleteProduct(ctx, productID, sellerID)
}

type ProductUpdateInput struct {
	ProductID      uint64
	SellerID       uint64
	CategoryID     uint64
	Title          string
	Description    string
	OriginalPrice  *float64
	Price          float64
	ConditionLevel string
	MeetLocation   string
}

func (s *Service) UpdateProduct(ctx context.Context, input ProductUpdateInput) error {
	if err := s.checkSensitive(ctx, input.Title, input.Description); err != nil {
		return err
	}

	return s.repo.UpdateProduct(ctx, UpdateProductInput{
		ProductID:      input.ProductID,
		SellerID:       input.SellerID,
		CategoryID:     input.CategoryID,
		Title:          input.Title,
		Description:    input.Description,
		OriginalPrice:  input.OriginalPrice,
		Price:          input.Price,
		ConditionLevel: input.ConditionLevel,
		MeetLocation:   input.MeetLocation,
	})
}

func (s *Service) UpdateProductStatus(ctx context.Context, productID uint64, sellerID uint64, status string, reason string) error {
	return s.repo.UpdateProductStatus(ctx, productID, sellerID, status, reason)
}

func (s *Service) BatchPutOnSaleRestorable(ctx context.Context, sellerID uint64) (int64, error) {
	if sellerID == 0 {
		return 0, fmt.Errorf("user not login")
	}
	return s.repo.BatchPutOnSaleRestorable(ctx, sellerID)
}

type AdminUpdateProductStatusInput struct {
	AdminID     uint64
	ProductID   uint64
	Status      string
	Reason      string
	IPAddress   *string
	RelatedType string
	RelatedID   uint64
}

func (s *Service) AdminUpdateProductStatus(ctx context.Context, input AdminUpdateProductStatusInput) error {
	if input.AdminID == 0 {
		return fmt.Errorf("adminId is required")
	}
	if input.ProductID == 0 {
		return fmt.Errorf("productId is required")
	}
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))
	input.Reason = strings.TrimSpace(input.Reason)
	relatedType, relatedID, err := normalizeAdminRelated(input.RelatedType, input.RelatedID)
	if err != nil {
		return err
	}
	if input.IPAddress != nil {
		trimmed := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &trimmed
	}
	if !isValidAdminProductStatus(input.Status) {
		return fmt.Errorf("status must be ON_SALE, OFF_SHELF, LOCKED or DELETED")
	}
	if len([]rune(input.Reason)) > 500 {
		return fmt.Errorf("reason cannot exceed 500 characters")
	}
	description := fmt.Sprintf("update product status to %s", input.Status)
	if input.Reason != "" {
		description = fmt.Sprintf("%s: %s", description, input.Reason)
	}
	productInfo, err := s.repo.GetProductNoticeInfo(ctx, input.ProductID)
	if err != nil {
		return fmt.Errorf("product not found")
	}
	if err := db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.ValidateRelatedRecordTx(ctx, tx, relatedType, relatedID, input.ProductID); err != nil {
			return err
		}
		if err := s.repo.AdminUpdateProductStatusTx(ctx, tx, input); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, input.AdminID, admin.OperationUpdateProductStatus, admin.TargetTypeProduct, input.ProductID, description, input.IPAddress, relatedTypePtr(relatedType), relatedIDPtr(relatedID))
	}); err != nil {
		return err
	}
	s.notifyAdminProductStatusChanged(ctx, input, productInfo)
	return nil
}

func (s *Service) notifyAdminProductStatusChanged(ctx context.Context, input AdminUpdateProductStatusInput, productInfo *ProductNoticeInfo) {
	if s.messageService == nil || productInfo == nil || productInfo.SellerID == 0 {
		return
	}

	var title string
	var content string
	switch input.Status {
	case "OFF_SHELF":
		title = "商品下架通知"
		content = fmt.Sprintf("你的商品「%s」已被管理员下架。", productInfo.Title)
	case "ON_SALE":
		title = "商品上架通知"
		content = fmt.Sprintf("你的商品「%s」已被管理员上架。", productInfo.Title)
	default:
		return
	}
	if input.Reason != "" {
		content = fmt.Sprintf("%s原因：%s", content, input.Reason)
	}

	relatedType := message.RelatedTypeProduct
	relatedID := productInfo.ID
	senderID := input.AdminID
	_, _ = s.messageService.Create(ctx, message.CreateMessageInput{
		ReceiverID:  productInfo.SellerID,
		SenderID:    &senderID,
		MessageType: message.MessageTypeSystemNotice,
		Title:       title,
		Content:     content,
		RelatedType: &relatedType,
		RelatedID:   &relatedID,
	})
}

func (s *Service) logAdminActionTx(ctx context.Context, tx *sql.Tx, adminID uint64, operationType, targetType string, targetID uint64, description string, ipAddress *string, relatedType *string, relatedID *uint64) error {
	if s.adminLogger == nil {
		return fmt.Errorf("admin logger is not configured")
	}
	_, err := s.adminLogger.LogActionTx(ctx, tx, admin.LogActionInput{
		AdminID:       adminID,
		OperationType: operationType,
		TargetType:    targetType,
		TargetID:      targetID,
		Description:   &description,
		IPAddress:     ipAddress,
		RelatedType:   relatedType,
		RelatedID:     relatedID,
	})
	return err
}

func isValidAdminProductStatus(status string) bool {
	switch status {
	case "ON_SALE", "OFF_SHELF", "LOCKED", "DELETED":
		return true
	default:
		return false
	}
}

func isValidAdminProductListStatus(status string) bool {
	switch status {
	case "ON_SALE", "OFF_SHELF", "LOCKED", "SOLD", "DELETED":
		return true
	default:
		return false
	}
}

func normalizeAdminRelated(relatedType string, relatedID uint64) (string, uint64, error) {
	relatedType = strings.ToUpper(strings.TrimSpace(relatedType))
	if (relatedType == "") != (relatedID == 0) {
		return "", 0, fmt.Errorf("relatedType and relatedId must be provided together")
	}
	if relatedType != "" && relatedType != admin.TargetTypeReport && relatedType != admin.TargetTypeAppeal {
		return "", 0, fmt.Errorf("relatedType must be REPORT or APPEAL")
	}
	return relatedType, relatedID, nil
}

func relatedTypePtr(relatedType string) *string {
	if relatedType == "" {
		return nil
	}
	return &relatedType
}

func relatedIDPtr(relatedID uint64) *uint64 {
	if relatedID == 0 {
		return nil
	}
	return &relatedID
}

func (s *Service) OffShelfOnSaleBySellerTx(ctx context.Context, tx *sql.Tx, sellerID uint64, reason string) ([]uint64, error) {
	return s.repo.OffShelfOnSaleBySellerTx(ctx, tx, sellerID, reason)
}

func (s *Service) OffShelfOnSaleBySellerWithSourceTx(ctx context.Context, tx *sql.Tx, sellerID uint64, reason string, source string) ([]uint64, error) {
	return s.repo.OffShelfOnSaleBySellerWithSourceTx(ctx, tx, sellerID, reason, source)
}

type ProductImagesInput struct {
	ProductID uint64
	SellerID  uint64
	Images    []string
}

func (s *Service) AddProductImages(ctx context.Context, input ProductImagesInput) error {
	return s.repo.AddProductImages(ctx, AddProductImagesInput{
		ProductID: input.ProductID,
		SellerID:  input.SellerID,
		Images:    input.Images,
	})
}

type ProductImageDeleteInput struct {
	ProductID uint64
	SellerID  uint64
	ImageID   uint64
}

func (s *Service) DeleteProductImage(ctx context.Context, input ProductImageDeleteInput) error {
	return s.repo.DeleteProductImage(ctx, DeleteProductImageInput{
		ProductID: input.ProductID,
		SellerID:  input.SellerID,
		ImageID:   input.ImageID,
	})
}

func (s *Service) ReplaceProductImages(ctx context.Context, input ProductImagesInput) error {
	return s.repo.ReplaceProductImages(ctx, ReplaceProductImagesInput{
		ProductID: input.ProductID,
		SellerID:  input.SellerID,
		Images:    input.Images,
	})
}

func (s *Service) checkSensitive(ctx context.Context, title string, description string) error {
	if s.sensitiveService == nil {
		return nil
	}

	checkText := strings.TrimSpace(title + " " + description)
	result, err := s.sensitiveService.CheckText(ctx, checkText)
	if err != nil {
		return err
	}

	if !result.Passed {
		return fmt.Errorf("%s：%s", result.Message, strings.Join(result.HitWords, "、"))
	}

	return nil
}

func (s *Service) LockProduct(ctx context.Context, productID uint64) error {
	return s.repo.LockProduct(ctx, productID)
}

func (s *Service) UnlockProduct(ctx context.Context, productID uint64) error {
	return s.repo.UnlockProduct(ctx, productID)
}

func (s *Service) MarkProductSold(ctx context.Context, productID uint64) error {
	return s.repo.MarkProductSold(ctx, productID)
}
