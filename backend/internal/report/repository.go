package report

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, report *Report) error {
	query := `
		INSERT INTO reports (reporter_id, target_type, target_id, reason_type, description, status)
		VALUES (?, ?, ?, ?, ?, 'PENDING')
	`
	result, err := r.db.ExecContext(ctx, query, report.ReporterID, report.TargetType, report.TargetID, report.ReasonType, report.Description)
	if err != nil {
		return fmt.Errorf("insert report: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	report.ID = uint64(id)
	return nil
}

func (r *Repository) AddImages(ctx context.Context, reportID uint64, imageURLs []string) error {
	if len(imageURLs) == 0 {
		return nil
	}
	query := "INSERT INTO report_images (report_id, image_url, sort_order) VALUES (?, ?, ?)"
	for i, url := range imageURLs {
		_, err := r.db.ExecContext(ctx, query, reportID, url, i)
		if err != nil {
			return fmt.Errorf("insert report image: %w", err)
		}
	}
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*ReportDetail, error) {
	query := `
		SELECT r.id, r.reporter_id, r.target_type, r.target_id, r.reason_type, r.description, r.status, r.handle_result, r.handler_id, r.handle_time, r.create_time, r.update_time,
			u.nickname
		FROM reports r
		LEFT JOIN users u ON u.id = r.reporter_id
		WHERE r.id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanReportDetail(row)
}

func (r *Repository) GetImages(ctx context.Context, reportID uint64) ([]string, error) {
	query := `SELECT image_url FROM report_images WHERE report_id = ? ORDER BY sort_order`
	rows, err := r.db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("get report images: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("scan image url: %w", err)
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (r *Repository) ListByReporter(ctx context.Context, reporterID uint64, page, pageSize int) ([]ReportDetail, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM reports WHERE reporter_id = ?`
	if err := r.db.QueryRowContext(ctx, countQuery, reporterID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reports: %w", err)
	}

	query := `
		SELECT r.id, r.reporter_id, r.target_type, r.target_id, r.reason_type, r.description, r.status, r.handle_result, r.handler_id, r.handle_time, r.create_time, r.update_time,
			u.nickname
		FROM reports r
		LEFT JOIN users u ON u.id = r.reporter_id
		WHERE r.reporter_id = ?
		ORDER BY r.create_time DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, reporterID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()

	var items []ReportDetail
	for rows.Next() {
		item, err := scanReportDetail(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, nil
}

func (r *Repository) ListAll(ctx context.Context, status string, page, pageSize int) ([]ReportDetail, int, error) {
	where := "1=1"
	args := []interface{}{}
	if status != "" {
		where += " AND r.status = ?"
		args = append(args, status)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM reports r WHERE ` + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reports: %w", err)
	}

	query := `
		SELECT r.id, r.reporter_id, r.target_type, r.target_id, r.reason_type, r.description, r.status, r.handle_result, r.handler_id, r.handle_time, r.create_time, r.update_time,
			u.nickname
		FROM reports r
		LEFT JOIN users u ON u.id = r.reporter_id
		WHERE ` + where + `
		ORDER BY r.create_time DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()

	var items []ReportDetail
	for rows.Next() {
		item, err := scanReportDetail(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	return items, total, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, reportID uint64, status string, handleResult *string, handlerID uint64) error {
	query := `UPDATE reports SET status = ?, handle_result = ?, handler_id = ?, handle_time = NOW() WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, handleResult, handlerID, reportID)
	if err != nil {
		return fmt.Errorf("update report status: %w", err)
	}
	return nil
}

func (r *Repository) HasReported(ctx context.Context, reporterID uint64, targetType string, targetID uint64) (bool, error) {
	query := `SELECT COUNT(*) FROM reports WHERE reporter_id = ? AND target_type = ? AND target_id = ? AND status IN ('PENDING', 'PROCESSING')`
	var count int
	if err := r.db.QueryRowContext(ctx, query, reporterID, targetType, targetID).Scan(&count); err != nil {
		return false, fmt.Errorf("check reported: %w", err)
	}
	return count > 0, nil
}

func scanReportDetail(row interface{ Scan(dest ...interface{}) error }) (*ReportDetail, error) {
	var rd ReportDetail
	var desc, handleResult sql.NullString
	var handlerID sql.NullInt64
	var handleTime sql.NullString
	var nick sql.NullString

	err := row.Scan(
		&rd.ID, &rd.ReporterID, &rd.TargetType, &rd.TargetID, &rd.ReasonType, &rd.Description, &rd.Status,
		&rd.HandleResult, &rd.HandlerID, &rd.HandleTime, &rd.CreateTime, &rd.UpdateTime,
		&nick,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan report: %w", err)
	}

	if desc.Valid {
		rd.Description = &desc.String
	}
	if handleResult.Valid {
		rd.HandleResult = &handleResult.String
	}
	if handlerID.Valid {
		v := uint64(handlerID.Int64)
		rd.HandlerID = &v
	}
	if handleTime.Valid {
		rd.HandleTime = &handleTime.String
	}
	if nick.Valid {
		rd.ReporterNickname = &nick.String
	}

	return &rd, nil
}
