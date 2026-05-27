package product

import (
	"context"
	"fmt"
	"strings"

	"cau-used-goods-app/backend/internal/sensitive"
)

type Service struct {
	repo             *Repository
	sensitiveService *sensitive.Service
}

func NewService(repo *Repository, sensitiveService *sensitive.Service) *Service {
	return &Service{
		repo:             repo,
		sensitiveService: sensitiveService,
	}
}

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
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
	Keyword    string
	CategoryID uint64
	Status     string
	MinPrice   *float64
	MaxPrice   *float64
	Sort       string
	Page       int
	PageSize   int
}

func (s *Service) ListProducts(ctx context.Context, input ProductListInput) (*ProductListResult, error) {
	return s.repo.ListProducts(ctx, ListProductsInput{
		Keyword:    input.Keyword,
		CategoryID: input.CategoryID,
		Status:     input.Status,
		MinPrice:   input.MinPrice,
		MaxPrice:   input.MaxPrice,
		Sort:       input.Sort,
		Page:       input.Page,
		PageSize:   input.PageSize,
	})
}

func (s *Service) GetProductByID(ctx context.Context, id uint64) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
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
