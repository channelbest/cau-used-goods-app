package message

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidReadStatus = errors.New("invalid read status")
var ErrInvalidMessageInput = errors.New("invalid message input")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateMessageInput) (uint64, error) {
	input.MessageType = strings.TrimSpace(input.MessageType)
	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)
	if input.RelatedType != nil {
		trimmed := strings.TrimSpace(*input.RelatedType)
		input.RelatedType = &trimmed
	}

	if err := validateCreateInput(input); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, input)
}

func (s *Service) ListByReceiver(ctx context.Context, receiverID uint64, readStatus string, page, pageSize int) ([]Message, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if err := validateReadStatus(readStatus); err != nil {
		return nil, 0, err
	}
	return s.repo.ListByReceiver(ctx, receiverID, readStatus, page, pageSize)
}

func (s *Service) GetByID(ctx context.Context, receiverID, messageID uint64) (*Message, error) {
	if messageID == 0 {
		return nil, ErrMessageNotFound
	}
	return s.repo.GetByID(ctx, receiverID, messageID)
}

func (s *Service) CountUnread(ctx context.Context, receiverID uint64) (int, error) {
	return s.repo.CountUnread(ctx, receiverID)
}

func (s *Service) MarkRead(ctx context.Context, receiverID, messageID uint64) error {
	if messageID == 0 {
		return ErrMessageNotFound
	}
	err := s.repo.MarkRead(ctx, receiverID, messageID)
	if errors.Is(err, ErrMessageNotFound) {
		return ErrMessageNotFound
	}
	return err
}

func (s *Service) MarkAllRead(ctx context.Context, receiverID uint64) (int64, error) {
	return s.repo.MarkAllRead(ctx, receiverID)
}

func validateReadStatus(readStatus string) error {
	if readStatus == "" || readStatus == ReadStatusUnread || readStatus == ReadStatusRead {
		return nil
	}
	return ErrInvalidReadStatus
}

func validateCreateInput(input CreateMessageInput) error {
	if input.ReceiverID == 0 {
		return ErrInvalidMessageInput
	}
	if input.MessageType == "" || len(input.MessageType) > 30 {
		return ErrInvalidMessageInput
	}
	if input.Title == "" || len(input.Title) > 100 {
		return ErrInvalidMessageInput
	}
	if input.Content == "" || len(input.Content) > 500 {
		return ErrInvalidMessageInput
	}
	if input.RelatedType != nil && len(*input.RelatedType) > 30 {
		return ErrInvalidMessageInput
	}
	return nil
}
