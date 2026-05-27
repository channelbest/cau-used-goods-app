package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/product"
)

type Service struct {
	repo    *Repository
	product *product.Service
}

func NewService(repo *Repository, productService *product.Service) *Service {
	return &Service{repo: repo, product: productService}
}

type CreateOrderInput struct {
	ProductID    uint64
	BuyerID      uint64
	Remark       *string
	MeetTime     *string
	MeetLocation *string
}

func (s *Service) Create(ctx context.Context, input CreateOrderInput) (*Order, error) {
	// 不能购买自己的商品
	sellerID, err := s.repo.GetProductSeller(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}
	if sellerID == input.BuyerID {
		return nil, fmt.Errorf("cannot buy your own product")
	}

	// 检查是否已有有效订单
	exists, err := s.repo.HasActiveOrderByBuyer(ctx, input.BuyerID, input.ProductID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("you already have an active order for this product")
	}

	title, price, err := s.repo.GetProductInfo(ctx, input.ProductID)
	if err != nil {
		return nil, err
	}

	expireTime := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")

	order := &Order{
		OrderNo:              generateOrderNo(),
		ProductID:            input.ProductID,
		BuyerID:              input.BuyerID,
		SellerID:             sellerID,
		ProductTitleSnapshot: title,
		ProductPriceSnapshot: price,
		Status:               "PENDING_CONFIRM",
		Remark:               input.Remark,
		MeetTime:             input.MeetTime,
		MeetLocation:         input.MeetLocation,
		ExpireTime:           expireTime,
	}

	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.product.LockProduct(ctx, input.ProductID); err != nil {
			return err
		}
		if err := s.repo.Create(ctx, tx, order); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

type ConfirmOrderInput struct {
	OrderID  uint64
	SellerID uint64
}

func (s *Service) Confirm(ctx context.Context, input ConfirmOrderInput) (*Order, error) {
	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.SellerID != input.SellerID {
		return nil, fmt.Errorf("permission denied")
	}
	if order.Status != "PENDING_CONFIRM" {
		return nil, fmt.Errorf("order cannot be confirmed")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	err = s.repo.UpdateStatus(ctx, nil, input.OrderID, "WAIT_MEET", map[string]interface{}{
		"confirm_time": now,
	})
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

type CancelOrderInput struct {
	OrderID  uint64
	UserID   uint64
	Reason   string
}

func (s *Service) Cancel(ctx context.Context, input CancelOrderInput) (*Order, error) {
	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.BuyerID != input.UserID && order.SellerID != input.UserID {
		return nil, fmt.Errorf("permission denied")
	}
	if order.Status != "PENDING_CONFIRM" && order.Status != "WAIT_MEET" {
		return nil, fmt.Errorf("order cannot be cancelled")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.product.UnlockProduct(ctx, order.ProductID); err != nil {
			return err
		}
		if err := s.repo.UpdateStatus(ctx, tx, input.OrderID, "CANCELED", map[string]interface{}{
			"cancel_reason": input.Reason,
			"cancel_by":     input.UserID,
			"close_time":    now,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

type CompleteOrderInput struct {
	OrderID  uint64
	SellerID uint64
}

type ExceptionCloseOrderInput struct {
	OrderID uint64
	UserID  uint64
	Reason  string
}

func (s *Service) Complete(ctx context.Context, input CompleteOrderInput) (*Order, error) {
	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.SellerID != input.SellerID {
		return nil, fmt.Errorf("permission denied")
	}
	if order.Status != "WAIT_MEET" {
		return nil, fmt.Errorf("order cannot be completed")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.product.MarkProductSold(ctx, order.ProductID); err != nil {
			return err
		}
		if err := s.repo.UpdateStatus(ctx, tx, input.OrderID, "COMPLETED", map[string]interface{}{
			"finish_time": now,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

func (s *Service) ExceptionClose(ctx context.Context, input ExceptionCloseOrderInput) (*Order, error) {
	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.BuyerID != input.UserID && order.SellerID != input.UserID {
		return nil, fmt.Errorf("permission denied")
	}
	if order.Status != "PENDING_CONFIRM" && order.Status != "WAIT_MEET" {
		return nil, fmt.Errorf("order cannot be exception closed")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.product.UnlockProduct(ctx, order.ProductID); err != nil {
			return err
		}
		if err := s.repo.UpdateStatus(ctx, tx, input.OrderID, "EXCEPTION_CLOSED", map[string]interface{}{
			"cancel_reason": input.Reason,
			"cancel_by":     input.UserID,
			"close_time":    now,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

func (s *Service) GetByID(ctx context.Context, orderID uint64) (*Order, error) {
	return s.repo.GetByID(ctx, orderID)
}

func (s *Service) ListByBuyer(ctx context.Context, buyerID uint64, status string, page, pageSize int) ([]OrderDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListByBuyer(ctx, buyerID, status, page, pageSize)
}

func (s *Service) ListBySeller(ctx context.Context, sellerID uint64, status string, page, pageSize int) ([]OrderDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListBySeller(ctx, sellerID, status, page, pageSize)
}

func (s *Service) CancelExpiredOrders(ctx context.Context) (int, error) {
	orders, err := s.repo.ListExpiredOrders(ctx)
	if err != nil {
		return 0, err
	}

	cancelled := 0
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, order := range orders {
		err = db.WithTx(ctx, func(tx *sql.Tx) error {
			if err := s.product.UnlockProduct(ctx, order.ProductID); err != nil {
				return err
			}
			if err := s.repo.UpdateStatus(ctx, tx, order.ID, "CANCELED", map[string]interface{}{
				"cancel_reason": "订单超时未确认，系统自动取消",
				"cancel_by":     order.BuyerID,
				"close_time":    now,
			}); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			continue
		}
		cancelled++
	}
	return cancelled, nil
}
