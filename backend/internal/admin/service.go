package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidAdminLogInput = errors.New("invalid admin log input")
var ErrInvalidAnnouncementInput = errors.New("invalid announcement input")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAnnouncement(ctx context.Context, input CreateAnnouncementInput) (uint64, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Status = strings.TrimSpace(input.Status)
	trimOptionalString(&input.Content)
	trimOptionalString(&input.CoverURL)
	trimOptionalString(&input.IPAddress)
	if input.Status == "" {
		input.Status = AnnouncementStatusDraft
	}

	if err := validateCreateAnnouncementInput(input); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateAnnouncement(ctx, input)
	if err != nil {
		return 0, err
	}

	description := fmt.Sprintf("create announcement: %s", input.Title)
	_, err = s.LogAction(ctx, LogActionInput{
		AdminID:       input.AdminID,
		OperationType: OperationCreateNotice,
		TargetType:    TargetTypeNotice,
		TargetID:      id,
		Description:   &description,
		IPAddress:     input.IPAddress,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) ListAnnouncements(ctx context.Context, query AnnouncementQuery) ([]Announcement, int, error) {
	query.Status = strings.TrimSpace(query.Status)
	query.Keyword = strings.TrimSpace(query.Keyword)
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}
	if query.Status != "" && !isValidAnnouncementStatus(query.Status) {
		return nil, 0, ErrInvalidAnnouncementInput
	}
	return s.repo.ListAnnouncements(ctx, query)
}

func (s *Service) UpdateAnnouncement(ctx context.Context, input UpdateAnnouncementInput) error {
	input.Title = strings.TrimSpace(input.Title)
	trimOptionalString(&input.Content)
	trimOptionalString(&input.CoverURL)
	trimOptionalString(&input.IPAddress)

	if err := validateUpdateAnnouncementInput(input); err != nil {
		return err
	}
	if err := s.repo.UpdateAnnouncement(ctx, input); err != nil {
		return err
	}

	description := fmt.Sprintf("update announcement: %s", input.Title)
	_, err := s.LogAction(ctx, LogActionInput{
		AdminID:       input.AdminID,
		OperationType: OperationUpdateNotice,
		TargetType:    TargetTypeNotice,
		TargetID:      input.ID,
		Description:   &description,
		IPAddress:     input.IPAddress,
	})
	return err
}

func (s *Service) UpdateAnnouncementStatus(ctx context.Context, input UpdateAnnouncementStatusInput) error {
	input.Status = strings.TrimSpace(input.Status)
	trimOptionalString(&input.IPAddress)

	if err := validateUpdateAnnouncementStatusInput(input); err != nil {
		return err
	}
	if err := s.repo.UpdateAnnouncementStatus(ctx, input.ID, input.Status); err != nil {
		return err
	}

	description := fmt.Sprintf("update announcement status: %s", input.Status)
	_, err := s.LogAction(ctx, LogActionInput{
		AdminID:       input.AdminID,
		OperationType: OperationStatusNotice,
		TargetType:    TargetTypeNotice,
		TargetID:      input.ID,
		Description:   &description,
		IPAddress:     input.IPAddress,
	})
	return err
}

func (s *Service) DeleteAnnouncement(ctx context.Context, adminID, id uint64, ipAddress *string) error {
	trimOptionalString(&ipAddress)
	if adminID == 0 || id == 0 {
		return ErrInvalidAnnouncementInput
	}
	if err := s.repo.UpdateAnnouncementStatus(ctx, id, AnnouncementStatusOffline); err != nil {
		return err
	}

	description := "offline announcement by delete operation"
	_, err := s.LogAction(ctx, LogActionInput{
		AdminID:       adminID,
		OperationType: OperationDeleteNotice,
		TargetType:    TargetTypeNotice,
		TargetID:      id,
		Description:   &description,
		IPAddress:     ipAddress,
	})
	return err
}

func (s *Service) LogAction(ctx context.Context, input LogActionInput) (uint64, error) {
	input.OperationType = strings.TrimSpace(input.OperationType)
	input.TargetType = strings.TrimSpace(input.TargetType)
	if input.Description != nil {
		value := strings.TrimSpace(*input.Description)
		input.Description = &value
	}
	if input.IPAddress != nil {
		value := strings.TrimSpace(*input.IPAddress)
		input.IPAddress = &value
	}

	if err := validateLogActionInput(input); err != nil {
		return 0, err
	}
	return s.repo.CreateLog(ctx, input)
}

func validateCreateAnnouncementInput(input CreateAnnouncementInput) error {
	if input.AdminID == 0 {
		return ErrInvalidAnnouncementInput
	}
	if input.Title == "" || len(input.Title) > 100 {
		return ErrInvalidAnnouncementInput
	}
	if input.Content != nil && len(*input.Content) > 65535 {
		return ErrInvalidAnnouncementInput
	}
	if input.CoverURL != nil && len(*input.CoverURL) > 255 {
		return ErrInvalidAnnouncementInput
	}
	if !isValidAnnouncementStatus(input.Status) {
		return ErrInvalidAnnouncementInput
	}
	return nil
}

func validateUpdateAnnouncementInput(input UpdateAnnouncementInput) error {
	if input.AdminID == 0 || input.ID == 0 {
		return ErrInvalidAnnouncementInput
	}
	if input.Title == "" || len(input.Title) > 100 {
		return ErrInvalidAnnouncementInput
	}
	if input.Content != nil && len(*input.Content) > 65535 {
		return ErrInvalidAnnouncementInput
	}
	if input.CoverURL != nil && len(*input.CoverURL) > 255 {
		return ErrInvalidAnnouncementInput
	}
	return nil
}

func validateUpdateAnnouncementStatusInput(input UpdateAnnouncementStatusInput) error {
	if input.AdminID == 0 || input.ID == 0 {
		return ErrInvalidAnnouncementInput
	}
	if !isValidAnnouncementStatus(input.Status) {
		return ErrInvalidAnnouncementInput
	}
	return nil
}

func isValidAnnouncementStatus(status string) bool {
	return status == AnnouncementStatusDraft ||
		status == AnnouncementStatusPublished ||
		status == AnnouncementStatusOffline
}

func trimOptionalString(value **string) {
	if value == nil || *value == nil {
		return
	}
	trimmed := strings.TrimSpace(**value)
	if trimmed == "" {
		*value = nil
		return
	}
	*value = &trimmed
}

func (s *Service) ListLogs(ctx context.Context, query LogQuery) ([]AdminLog, int, error) {
	query.OperationType = strings.TrimSpace(query.OperationType)
	query.TargetType = strings.TrimSpace(query.TargetType)
	query.StartTime = strings.TrimSpace(query.StartTime)
	query.EndTime = strings.TrimSpace(query.EndTime)

	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 100 {
		query.PageSize = 20
	}

	return s.repo.ListLogs(ctx, query)
}

func validateLogActionInput(input LogActionInput) error {
	if input.AdminID == 0 {
		return ErrInvalidAdminLogInput
	}
	if input.OperationType == "" || len(input.OperationType) > 50 {
		return ErrInvalidAdminLogInput
	}
	if input.TargetType == "" || len(input.TargetType) > 30 {
		return ErrInvalidAdminLogInput
	}
	if input.TargetID == 0 {
		return ErrInvalidAdminLogInput
	}
	if input.Description != nil && len(*input.Description) > 500 {
		return ErrInvalidAdminLogInput
	}
	if input.IPAddress != nil && len(*input.IPAddress) > 50 {
		return ErrInvalidAdminLogInput
	}
	return nil
}
