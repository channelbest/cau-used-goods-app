package order

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"cau-used-goods-app/backend/internal/admin"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/message"
)

type Service struct {
	repo    *Repository
	message *message.Service
	admin   *admin.Service
}

func NewService(repo *Repository, messageService *message.Service, adminService *admin.Service) *Service {
	return &Service{repo: repo, message: messageService, admin: adminService}
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
		if err := s.repo.LockProduct(ctx, tx, input.ProductID); err != nil {
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

	// 发送订单创建消息通知卖家
	if s.message != nil {
		relatedType := message.RelatedTypeOrder
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  sellerID,
			MessageType: message.MessageTypeOrderCreated,
			Title:       "新订单提醒",
			Content:     fmt.Sprintf("您的商品「%s」有新的订单，买家已预约，请尽快确认。", title),
			RelatedType: &relatedType,
			RelatedID:   &order.ID,
		})
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

	// 发送订单确认消息通知买家
	if s.message != nil {
		relatedType := message.RelatedTypeOrder
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  order.BuyerID,
			MessageType: message.MessageTypeOrderConfirmed,
			Title:       "订单已确认",
			Content:     fmt.Sprintf("卖家已确认您的订单「%s」，请按约定时间地点交易。", order.ProductTitleSnapshot),
			RelatedType: &relatedType,
			RelatedID:   &order.ID,
		})
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

type CancelOrderInput struct {
	OrderID uint64
	UserID  uint64
	Reason  string
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
		if err := s.repo.UnlockProduct(ctx, tx, order.ProductID); err != nil {
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

	// 发送订单取消消息通知对方
	if s.message != nil {
		relatedType := message.RelatedTypeOrder
		receiverID := order.BuyerID
		if input.UserID == order.BuyerID {
			receiverID = order.SellerID
		}
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  receiverID,
			MessageType: message.MessageTypeOrderCanceled,
			Title:       "订单已取消",
			Content:     fmt.Sprintf("订单「%s」已被取消，原因：%s", order.ProductTitleSnapshot, input.Reason),
			RelatedType: &relatedType,
			RelatedID:   &order.ID,
		})
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

type CompleteOrderInput struct {
	OrderID  uint64
	SellerID uint64
}

type AdminExceptionCloseOrderInput struct {
	OrderID          uint64
	AdminID          uint64
	Reason           string
	ResponsibleParty string
	IPAddress        string
	RelatedType      string
	RelatedID        uint64
}

type AdminUpdateOrderStatusInput struct {
	AdminID     uint64
	OrderID     uint64
	Status      string
	Reason      string
	IPAddress   *string
	RelatedType string
	RelatedID   uint64
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
		if err := s.repo.MarkProductSold(ctx, tx, order.ProductID); err != nil {
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

	// 发送订单完成消息通知买家
	if s.message != nil {
		relatedType := message.RelatedTypeOrder
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  order.BuyerID,
			MessageType: message.MessageTypeOrderConfirmed,
			Title:       "交易完成",
			Content:     fmt.Sprintf("订单「%s」已完成交易，欢迎评价。", order.ProductTitleSnapshot),
			RelatedType: &relatedType,
			RelatedID:   &order.ID,
		})
	}

	return s.repo.GetByID(ctx, input.OrderID)
}

func (s *Service) AdminExceptionClose(ctx context.Context, input AdminExceptionCloseOrderInput) (*Order, error) {
	input.Reason = strings.TrimSpace(input.Reason)
	input.ResponsibleParty = strings.ToUpper(strings.TrimSpace(input.ResponsibleParty))
	input.IPAddress = strings.TrimSpace(input.IPAddress)
	relatedType, relatedID, err := normalizeAdminRelated(input.RelatedType, input.RelatedID)
	if err != nil {
		return nil, err
	}
	if input.AdminID == 0 {
		return nil, fmt.Errorf("adminId is required")
	}
	if input.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}
	if len([]rune(input.Reason)) > 500 {
		return nil, fmt.Errorf("reason cannot exceed 500 characters")
	}
	if input.ResponsibleParty != "BUYER" && input.ResponsibleParty != "SELLER" {
		return nil, fmt.Errorf("responsibleParty must be BUYER or SELLER")
	}

	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.Status != "PENDING_CONFIRM" && order.Status != "WAIT_MEET" {
		return nil, fmt.Errorf("order cannot be exception closed")
	}

	closedOrder := AccountStatusClosedOrder{
		ID:                   order.ID,
		BuyerID:              order.BuyerID,
		SellerID:             order.SellerID,
		ProductID:            order.ProductID,
		ProductTitleSnapshot: order.ProductTitleSnapshot,
		ResponsibleParty:     input.ResponsibleParty,
	}
	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.ValidateRelatedRecordTx(ctx, tx, relatedType, relatedID, order.ID); err != nil {
			return err
		}
		description := fmt.Sprintf("responsibleParty=%s; reason=%s", input.ResponsibleParty, input.Reason)
		return s.exceptionCloseOrderTx(ctx, tx, closedOrder, input.AdminID, input.Reason, &input.IPAddress, relatedTypePtr(relatedType), relatedIDPtr(relatedID), description)
	})
	if err != nil {
		return nil, err
	}

	s.notifyOrderExceptionClosed(ctx, input.AdminID, closedOrder, fmt.Sprintf("订单「%s」已由管理员异常关闭，原因：%s", order.ProductTitleSnapshot, input.Reason))

	return s.repo.GetByID(ctx, input.OrderID)
}

func (s *Service) AdminUpdateStatus(ctx context.Context, input AdminUpdateOrderStatusInput) (*Order, error) {
	if input.AdminID == 0 {
		return nil, fmt.Errorf("adminId is required")
	}
	if input.OrderID == 0 {
		return nil, fmt.Errorf("orderId is required")
	}
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))
	input.Reason = strings.TrimSpace(input.Reason)
	relatedType, relatedID, err := normalizeAdminRelated(input.RelatedType, input.RelatedID)
	if err != nil {
		return nil, err
	}
	if input.IPAddress != nil {
		trimmed := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &trimmed
	}
	if !isValidAdminOrderStatus(input.Status) {
		return nil, fmt.Errorf("status must be PENDING_CONFIRM, WAIT_MEET, COMPLETED or CANCELED; use exception-close for EXCEPTION_CLOSED")
	}
	if len([]rune(input.Reason)) > 500 {
		return nil, fmt.Errorf("reason cannot exceed 500 characters")
	}

	order, err := s.repo.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	updates := map[string]interface{}{}
	switch input.Status {
	case "WAIT_MEET":
		updates["confirm_time"] = now
	case "COMPLETED":
		updates["finish_time"] = now
	case "CANCELED":
		updates["cancel_reason"] = input.Reason
		updates["cancel_by"] = input.AdminID
		updates["close_time"] = now
	}

	err = db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.ValidateRelatedRecordTx(ctx, tx, relatedType, relatedID, order.ID); err != nil {
			return err
		}
		switch input.Status {
		case "PENDING_CONFIRM", "WAIT_MEET":
			if err := s.repo.UpdateProductStatusForAdmin(ctx, tx, order.ProductID, "LOCKED"); err != nil {
				return err
			}
		case "COMPLETED":
			if err := s.repo.UpdateProductStatusForAdmin(ctx, tx, order.ProductID, "SOLD"); err != nil {
				return err
			}
		case "CANCELED":
			if err := s.repo.UpdateProductStatusForAdmin(ctx, tx, order.ProductID, "ON_SALE"); err != nil {
				return err
			}
		}
		if err := s.repo.UpdateStatus(ctx, tx, input.OrderID, input.Status, updates); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, input.AdminID, admin.OperationUpdateOrderStatus, admin.TargetTypeOrder, input.OrderID, buildStatusDescription(input.Status, input.Reason), input.IPAddress, relatedTypePtr(relatedType), relatedIDPtr(relatedID))
	})
	if err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, input.OrderID)
}

func (s *Service) logAdminActionTx(ctx context.Context, tx *sql.Tx, adminID uint64, operationType, targetType string, targetID uint64, description string, ipAddress *string, relatedType *string, relatedID *uint64) error {
	if s.admin == nil {
		return fmt.Errorf("admin logger is not configured")
	}
	_, err := s.admin.LogActionTx(ctx, tx, admin.LogActionInput{
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

func (s *Service) exceptionCloseOrderTx(ctx context.Context, tx *sql.Tx, item AccountStatusClosedOrder, adminID uint64, reason string, ipAddress *string, relatedType *string, relatedID *uint64, description string) error {
	if item.ResponsibleParty == "BUYER" {
		if err := s.repo.UnlockProduct(ctx, tx, item.ProductID); err != nil {
			return err
		}
	} else {
		if err := s.repo.OffShelfLockedProduct(ctx, tx, item.ProductID, reason); err != nil {
			return err
		}
	}
	if err := s.repo.UpdateStatus(ctx, tx, item.ID, "EXCEPTION_CLOSED", map[string]interface{}{
		"cancel_reason": reason,
		"cancel_by":     adminID,
		"close_time":    time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		return err
	}
	return s.logAdminActionTx(ctx, tx, adminID, admin.OperationOrderExceptionClose, admin.TargetTypeOrder, item.ID, description, ipAddress, relatedType, relatedID)
}

func (s *Service) notifyOrderExceptionClosed(ctx context.Context, adminID uint64, item AccountStatusClosedOrder, content string) {
	if s.message == nil {
		return
	}
	relatedType := message.RelatedTypeOrder
	for _, receiverID := range []uint64{item.BuyerID, item.SellerID} {
		_, _ = s.message.Create(ctx, message.CreateMessageInput{
			ReceiverID:  receiverID,
			SenderID:    &adminID,
			MessageType: message.MessageTypeSystemNotice,
			Title:       "订单异常关闭",
			Content:     content,
			RelatedType: &relatedType,
			RelatedID:   &item.ID,
		})
	}
}

func isValidAdminOrderStatus(status string) bool {
	switch status {
	case "PENDING_CONFIRM", "WAIT_MEET", "COMPLETED", "CANCELED":
		return true
	default:
		return false
	}
}

func buildStatusDescription(status string, reason string) string {
	description := fmt.Sprintf("update order status to %s", status)
	if reason != "" {
		description = fmt.Sprintf("%s: %s", description, reason)
	}
	return description
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

func (s *Service) ListAll(ctx context.Context, status string, page, pageSize int) ([]OrderDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListAll(ctx, status, page, pageSize)
}

func (s *Service) CountBlockingOrdersTx(ctx context.Context, tx *sql.Tx, userID uint64) (int, int, int, int, error) {
	return s.repo.CountBlockingOrdersTx(ctx, tx, userID)
}

func (s *Service) AutoExceptionClosePendingConfirmByUserTx(ctx context.Context, tx *sql.Tx, userID, adminID uint64, reason string, ipAddress *string, relatedType *string, relatedID *uint64) ([]AccountStatusClosedOrder, error) {
	if s.admin == nil {
		return nil, fmt.Errorf("admin logger is not configured")
	}
	orders, err := s.repo.ListPendingConfirmByUserTx(ctx, tx, userID)
	if err != nil {
		return nil, err
	}
	for i := range orders {
		responsibleParty := "BUYER"
		if orders[i].SellerID == userID {
			responsibleParty = "SELLER"
		}
		orders[i].ResponsibleParty = responsibleParty

		description := fmt.Sprintf("因用户账号状态变更，系统自动异常关闭订单；责任方=%s，原因=%s", responsibleParty, reason)
		if err := s.exceptionCloseOrderTx(ctx, tx, orders[i], adminID, reason, ipAddress, relatedType, relatedID, description); err != nil {
			return nil, err
		}
	}
	return orders, nil
}

func (s *Service) CancelExpiredOrders(ctx context.Context) (int, error) {
	orders, err := s.repo.ListExpiredOrders(ctx)
	if err != nil {
		return 0, err
	}

	cancelled := 0
	var firstErr error
	for _, order := range orders {
		cancelledThisOrder := false
		err = db.WithTx(ctx, func(tx *sql.Tx) error {
			affected, err := s.repo.CancelExpiredOrderTx(ctx, tx, order.ID, order.BuyerID)
			if err != nil {
				return err
			}
			if !affected {
				return nil
			}
			cancelledThisOrder = true
			return s.repo.UnlockProductIfLocked(ctx, tx, order.ProductID)
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		if !cancelledThisOrder {
			continue
		}

		// 发送超时取消消息通知买卖双方
		if s.message != nil {
			relatedType := message.RelatedTypeOrder
			_, _ = s.message.Create(ctx, message.CreateMessageInput{
				ReceiverID:  order.BuyerID,
				MessageType: message.MessageTypeOrderTimeout,
				Title:       "订单超时取消",
				Content:     fmt.Sprintf("订单「%s」因超时未确认，已自动取消。", order.ProductTitleSnapshot),
				RelatedType: &relatedType,
				RelatedID:   &order.ID,
			})
			_, _ = s.message.Create(ctx, message.CreateMessageInput{
				ReceiverID:  order.SellerID,
				MessageType: message.MessageTypeOrderTimeout,
				Title:       "订单超时取消",
				Content:     fmt.Sprintf("订单「%s」因超时未确认，已自动取消。", order.ProductTitleSnapshot),
				RelatedType: &relatedType,
				RelatedID:   &order.ID,
			})
		}

		cancelled++
	}
	return cancelled, firstErr
}
