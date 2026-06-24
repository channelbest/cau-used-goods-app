package review

import (
	"context"
	"fmt"

	"cau-used-goods-app/backend/internal/message"
)

type Service struct {
	repo    *Repository
	message *message.Service
}

func NewService(repo *Repository, messageService *message.Service) *Service {
	return &Service{repo: repo, message: messageService}
}

type CreateReviewInput struct {
	OrderID    uint64
	ReviewerID uint64
	Rating     int
	Content    *string
	Anonymous  bool
}

func (s *Service) Create(ctx context.Context, input CreateReviewInput) (*Review, error) {
	productID, buyerID, sellerID, status, err := s.repo.GetOrderInfo(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}

	// 只有买家可以评价
	if buyerID != input.ReviewerID {
		return nil, fmt.Errorf("only buyer can review")
	}

	// 订单必须已完成
	if status != "COMPLETED" {
		return nil, fmt.Errorf("can only review completed orders")
	}

	// 检查是否已评价
	existing, err := s.repo.GetByOrderAndReviewer(ctx, input.OrderID, input.ReviewerID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("you have already reviewed this order")
	}

	// 评分范围校验
	if input.Rating < 1 || input.Rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	review := &Review{
		OrderID:    input.OrderID,
		ProductID:  productID,
		ReviewerID: input.ReviewerID,
		SellerID:   sellerID,
		Rating:     input.Rating,
		Content:    input.Content,
		Anonymous:  input.Anonymous,
		Status:     "NORMAL",
	}

	if err := s.repo.Create(ctx, review); err != nil {
		return nil, err
	}

	// 发送评价通知给卖家
	if s.message != nil {
		relatedType := message.RelatedTypeOrder
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  sellerID,
			MessageType: message.MessageTypeSystemNotice,
			Title:       "收到新评价",
			Content:     fmt.Sprintf("您的订单收到%d星评价", input.Rating),
			RelatedType: &relatedType,
			RelatedID:   &input.OrderID,
		})
	}

	return review, nil
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*Review, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListByProduct(ctx context.Context, productID uint64, page, pageSize int) ([]ReviewDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListByProduct(ctx, productID, page, pageSize)
}

func (s *Service) ListBySeller(ctx context.Context, sellerID uint64, page, pageSize int) ([]ReviewDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListBySeller(ctx, sellerID, page, pageSize)
}

func (s *Service) GetSellerRating(ctx context.Context, sellerID uint64) (float64, int, error) {
	return s.repo.GetAverageRating(ctx, sellerID)
}
