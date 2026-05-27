package report

import (
	"context"
	"database/sql"
	"fmt"

	"cau-used-goods-app/backend/internal/sensitive"
)

type Service struct {
	repo      *Repository
	db        *sql.DB
	sensitive *sensitive.Service
}

func NewService(repo *Repository, db *sql.DB, sensitiveService *sensitive.Service) *Service {
	return &Service{repo: repo, db: db, sensitive: sensitiveService}
}

type CreateReportInput struct {
	ReporterID  uint64
	TargetType  string
	TargetID    uint64
	ReasonType  string
	Description *string
	Images      []string
}

func (s *Service) Create(ctx context.Context, input CreateReportInput) (*ReportDetail, error) {
	// 检查是否已举报
	reported, err := s.repo.HasReported(ctx, input.ReporterID, input.TargetType, input.TargetID)
	if err != nil {
		return nil, err
	}
	if reported {
		return nil, fmt.Errorf("you have already reported this target")
	}

	// 敏感词检测
	if input.Description != nil && *input.Description != "" {
		checkResult, err := s.sensitive.CheckText(ctx, *input.Description)
		if err != nil {
			return nil, fmt.Errorf("sensitive word check failed: %w", err)
		}
		if !checkResult.Passed {
			return nil, fmt.Errorf("report description contains sensitive words: %v", checkResult.HitWords)
		}
	}

	report := &Report{
		ReporterID:  input.ReporterID,
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
		ReasonType:  input.ReasonType,
		Description: input.Description,
		Status:      "PENDING",
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}

	if len(input.Images) > 0 {
		if err := s.repo.AddImages(ctx, report.ID, input.Images); err != nil {
			return nil, err
		}
	}

	return s.GetByID(ctx, report.ID)
}

func (s *Service) GetByID(ctx context.Context, id uint64) (*ReportDetail, error) {
	rd, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rd == nil {
		return nil, fmt.Errorf("report not found")
	}
	images, err := s.repo.GetImages(ctx, id)
	if err != nil {
		return nil, err
	}
	rd.Images = images
	return rd, nil
}

func (s *Service) ListByReporter(ctx context.Context, reporterID uint64, page, pageSize int) ([]ReportDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	items, total, err := s.repo.ListByReporter(ctx, reporterID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		images, err := s.repo.GetImages(ctx, items[i].ID)
		if err != nil {
			return nil, 0, err
		}
		items[i].Images = images
	}
	return items, total, nil
}

func (s *Service) ListAll(ctx context.Context, status string, page, pageSize int) ([]ReportDetail, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	items, total, err := s.repo.ListAll(ctx, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		images, err := s.repo.GetImages(ctx, items[i].ID)
		if err != nil {
			return nil, 0, err
		}
		items[i].Images = images
	}
	return items, total, nil
}

type HandleReportInput struct {
	ReportID     uint64
	HandlerID    uint64
	Status       string
	HandleResult *string
}

func (s *Service) Handle(ctx context.Context, input HandleReportInput) (*ReportDetail, error) {
	report, err := s.repo.GetByID(ctx, input.ReportID)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, fmt.Errorf("report not found")
	}
	if report.Status != "PENDING" && report.Status != "PROCESSING" {
		return nil, fmt.Errorf("report cannot be handled")
	}

	if err := s.repo.UpdateStatus(ctx, input.ReportID, input.Status, input.HandleResult, input.HandlerID); err != nil {
		return nil, err
	}

	// 写入管理员操作日志
	_ = s.logAdminAction(ctx, input.HandlerID, input.ReportID, input.Status, input.HandleResult)

	// TODO: 通知举报人处理结果（待 messages 模块实现后对接）
	// _ = s.notifyReporter(ctx, report.ReporterID, input.ReportID, input.Status)

	return s.GetByID(ctx, input.ReportID)
}

func (s *Service) logAdminAction(ctx context.Context, adminID, reportID uint64, status string, handleResult *string) error {
	var actionType string
	switch status {
	case "RESOLVED":
		actionType = "REPORT_RESOLVE"
	case "REJECTED":
		actionType = "REPORT_REJECT"
	case "CLOSED":
		actionType = "REPORT_CLOSE"
	default:
		actionType = "REPORT_HANDLE"
	}

	result := ""
	if handleResult != nil {
		result = *handleResult
	}

	query := `
		INSERT INTO admin_logs (admin_id, target_type, target_id, action_type, detail)
		VALUES (?, 'REPORT', ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, adminID, reportID, actionType, result)
	return err
}
