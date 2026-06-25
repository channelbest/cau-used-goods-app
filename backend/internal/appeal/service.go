package appeal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"cau-used-goods-app/backend/internal/admin"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/message"
)

type Service struct {
	repo     *Repository
	admin    *admin.Service
	messages *message.Service
}

func NewService(repo *Repository, adminService *admin.Service, messageService *message.Service) *Service {
	return &Service{repo: repo, admin: adminService, messages: messageService}
}

func (s *Service) Create(ctx context.Context, input CreateAppealInput) (*Appeal, error) {
	input.TargetType = strings.ToUpper(strings.TrimSpace(input.TargetType))
	input.Reason = strings.TrimSpace(input.Reason)
	input.EvidenceURLs = normalizeEvidenceURLs(input.EvidenceURLs)

	if err := validateCreateInput(input); err != nil {
		return nil, err
	}
	exists, err := s.repo.TargetExists(ctx, input.TargetType, input.TargetID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("appeal target not found")
	}
	allowed, err := s.repo.TargetAppealableBy(ctx, input.AppellantID, input.TargetType, input.TargetID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, fmt.Errorf("permission denied")
	}
	if input.TargetType == TargetTypeReport {
		appealable, err := s.repo.ReportResultAppealable(ctx, input.TargetID)
		if err != nil {
			return nil, err
		}
		if !appealable {
			return nil, fmt.Errorf("report result can only be appealed after handled")
		}
	}
	active, err := s.repo.HasActiveAppeal(ctx, input.AppellantID, input.TargetType, input.TargetID)
	if err != nil {
		return nil, err
	}
	if active {
		return nil, fmt.Errorf("you already have an active appeal for this target")
	}
	return s.repo.Create(ctx, input)
}

func (s *Service) ListMy(ctx context.Context, appellantID uint64, query AppealQuery) ([]AppealDetail, int, error) {
	query.AppellantID = appellantID
	return s.List(ctx, query)
}

func (s *Service) List(ctx context.Context, query AppealQuery) ([]AppealDetail, int, error) {
	query.TargetType = strings.ToUpper(strings.TrimSpace(query.TargetType))
	query.Status = strings.ToUpper(strings.TrimSpace(query.Status))
	if query.TargetType != "" && !isValidTargetType(query.TargetType) {
		return nil, 0, fmt.Errorf("invalid targetType")
	}
	if query.Status != "" && !isValidStatus(query.Status) {
		return nil, 0, fmt.Errorf("invalid status")
	}
	query.Page, query.PageSize = normalizePage(query.Page, query.PageSize)
	return s.repo.List(ctx, query)
}

func (s *Service) GetByID(ctx context.Context, appealID, userID uint64, isAdmin bool) (*AppealDetail, error) {
	item, err := s.repo.GetDetailByID(ctx, appealID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("appeal not found")
	}
	if !isAdmin && item.AppellantID != userID {
		return nil, fmt.Errorf("permission denied")
	}
	return item, nil
}

func (s *Service) Handle(ctx context.Context, input HandleAppealInput) (*Appeal, error) {
	input.Status = strings.ToUpper(strings.TrimSpace(input.Status))
	input.HandleResult = strings.TrimSpace(input.HandleResult)
	if input.IPAddress != nil {
		trimmed := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &trimmed
	}

	if err := validateHandleInput(input); err != nil {
		return nil, err
	}

	description := fmt.Sprintf("handle appeal #%d: %s", input.AppealID, input.Status)
	if err := db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.HandleTx(ctx, tx, input); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, input.AdminID, input.AppealID, appealOperationType(input.Status), description, input.IPAddress)
	}); err != nil {
		return nil, err
	}
	item, err := s.repo.GetByID(ctx, input.AppealID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("appeal not found")
	}

	title := "申诉处理结果"
	content := limitRunes(fmt.Sprintf("你的申诉已处理，结果：%s。处理说明：%s", input.Status, input.HandleResult), 500)
	relatedType := message.RelatedTypeAppeal
	relatedID := item.ID
	_, _ = s.messages.Create(ctx, message.CreateMessageInput{
		ReceiverID:  item.AppellantID,
		SenderID:    &input.AdminID,
		MessageType: message.MessageTypeSystemNotice,
		Title:       title,
		Content:     content,
		RelatedType: &relatedType,
		RelatedID:   &relatedID,
	})

	return item, nil
}

func (s *Service) Close(ctx context.Context, input CloseAppealInput) (*Appeal, error) {
	input.CloseReason = strings.TrimSpace(input.CloseReason)
	if input.AppealID == 0 {
		return nil, fmt.Errorf("appealId is required")
	}
	if input.AppellantID == 0 {
		return nil, fmt.Errorf("appellantId is required")
	}
	if len([]rune(input.CloseReason)) > 500 {
		return nil, fmt.Errorf("closeReason cannot exceed 500 characters")
	}

	item, err := s.repo.GetDetailByID(ctx, input.AppealID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("appeal not found")
	}
	if item.AppellantID != input.AppellantID {
		return nil, fmt.Errorf("permission denied")
	}
	if item.Status != StatusPending {
		return nil, fmt.Errorf("appeal cannot be closed")
	}

	return s.repo.Close(ctx, input)
}

func (s *Service) MarkProcessing(ctx context.Context, appealID, adminID uint64, ipAddress *string) (*Appeal, error) {
	if appealID == 0 {
		return nil, fmt.Errorf("appealId is required")
	}
	if adminID == 0 {
		return nil, fmt.Errorf("adminId is required")
	}
	if ipAddress != nil {
		trimmed := strings.TrimSpace(*ipAddress)
		ipAddress = &trimmed
	}

	description := fmt.Sprintf("mark appeal #%d as PROCESSING", appealID)
	if err := db.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.MarkProcessingTx(ctx, tx, appealID, adminID); err != nil {
			return err
		}
		return s.logAdminActionTx(ctx, tx, adminID, appealID, admin.OperationMarkAppealProcessing, description, ipAddress)
	}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, appealID)
}

func appealOperationType(status string) string {
	switch status {
	case StatusApproved:
		return admin.OperationApproveAppeal
	case StatusRejected:
		return admin.OperationRejectAppeal
	default:
		return admin.OperationMarkAppealProcessing
	}
}

func (s *Service) logAdminActionTx(ctx context.Context, tx *sql.Tx, adminID, appealID uint64, operationType string, description string, ipAddress *string) error {
	if s.admin == nil {
		return fmt.Errorf("admin logger is not configured")
	}
	_, err := s.admin.LogActionTx(ctx, tx, admin.LogActionInput{
		AdminID:       adminID,
		OperationType: operationType,
		TargetType:    admin.TargetTypeAppeal,
		TargetID:      appealID,
		Description:   &description,
		IPAddress:     ipAddress,
	})
	return err
}

func validateCreateInput(input CreateAppealInput) error {
	if input.AppellantID == 0 {
		return fmt.Errorf("appellantId is required")
	}
	if !isValidTargetType(input.TargetType) {
		return fmt.Errorf("targetType must be PRODUCT, USER, ORDER or REPORT")
	}
	if input.TargetID == 0 {
		return fmt.Errorf("targetId is required")
	}
	if input.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if len([]rune(input.Reason)) > 500 {
		return fmt.Errorf("reason cannot exceed 500 characters")
	}
	if len(input.EvidenceURLs) > 9 {
		return fmt.Errorf("evidenceUrls cannot exceed 9")
	}
	return nil
}

func validateHandleInput(input HandleAppealInput) error {
	if input.AppealID == 0 {
		return fmt.Errorf("appealId is required")
	}
	if input.AdminID == 0 {
		return fmt.Errorf("adminId is required")
	}
	if input.Status != StatusApproved && input.Status != StatusRejected {
		return fmt.Errorf("status must be APPROVED or REJECTED")
	}
	if input.HandleResult == "" {
		return fmt.Errorf("handleResult is required")
	}
	if len([]rune(input.HandleResult)) > 500 {
		return fmt.Errorf("handleResult cannot exceed 500 characters")
	}
	return nil
}

func normalizeEvidenceURLs(urls []string) []string {
	result := make([]string, 0, len(urls))
	for _, url := range urls {
		trimmed := strings.TrimSpace(url)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func isValidTargetType(targetType string) bool {
	return targetType == TargetTypeProduct ||
		targetType == TargetTypeUser ||
		targetType == TargetTypeOrder ||
		targetType == TargetTypeReport
}

func isValidStatus(status string) bool {
	return status == StatusPending ||
		status == StatusProcessing ||
		status == StatusApproved ||
		status == StatusRejected ||
		status == StatusClosed
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func limitRunes(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max])
}
