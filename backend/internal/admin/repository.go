package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrAnnouncementNotFound = errors.New("announcement not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAnnouncement(ctx context.Context, input CreateAnnouncementInput) (uint64, error) {
	query := `
		INSERT INTO announcements (title, content, cover_url, status, publish_time, create_by)
		VALUES (?, ?, ?, ?, CASE WHEN ? = 'PUBLISHED' THEN NOW() ELSE NULL END, ?)
	`
	result, err := r.db.ExecContext(ctx, query, input.Title, input.Content, input.CoverURL, input.Status, input.Status, input.AdminID)
	if err != nil {
		return 0, fmt.Errorf("create announcement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created announcement id: %w", err)
	}
	return uint64(id), nil
}

func (r *Repository) ListAnnouncements(ctx context.Context, query AnnouncementQuery) ([]Announcement, int, error) {
	whereSQL, args := buildAnnouncementWhere(query)

	var total int
	countSQL := `SELECT COUNT(*) FROM announcements` + whereSQL
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count announcements: %w", err)
	}

	listSQL := `
		SELECT id, title, content, cover_url, status, publish_time, create_by, create_time, update_time
		FROM announcements
	` + whereSQL + ` ORDER BY update_time DESC LIMIT ? OFFSET ?`
	args = append(args, query.PageSize, (query.Page-1)*query.PageSize)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list announcements: %w", err)
	}
	defer rows.Close()

	items := make([]Announcement, 0)
	for rows.Next() {
		item, err := scanAnnouncement(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate announcements: %w", err)
	}

	return items, total, nil
}

func (r *Repository) UpdateAnnouncement(ctx context.Context, input UpdateAnnouncementInput) error {
	query := `
		UPDATE announcements
		SET title = ?, content = ?, cover_url = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, input.Title, input.Content, input.CoverURL, input.ID)
	if err != nil {
		return fmt.Errorf("update announcement: %w", err)
	}
	return checkAnnouncementRowsAffected(result)
}

func (r *Repository) UpdateAnnouncementStatus(ctx context.Context, id uint64, status string) error {
	query := `
		UPDATE announcements
		SET status = ?, publish_time = CASE WHEN ? = 'PUBLISHED' THEN NOW() ELSE NULL END
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, status, status, id)
	if err != nil {
		return fmt.Errorf("update announcement status: %w", err)
	}
	return checkAnnouncementRowsAffected(result)
}

func (r *Repository) CreateLog(ctx context.Context, input LogActionInput) (uint64, error) {
	query := `
		INSERT INTO admin_logs (
			admin_id, operation_type, target_type, target_id, description, ip_address
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		input.AdminID,
		input.OperationType,
		input.TargetType,
		input.TargetID,
		input.Description,
		input.IPAddress,
	)
	if err != nil {
		return 0, fmt.Errorf("create admin log: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created admin log id: %w", err)
	}
	return uint64(id), nil
}

func buildAnnouncementWhere(query AnnouncementQuery) (string, []interface{}) {
	whereSQL := " WHERE 1 = 1"
	args := make([]interface{}, 0)

	if query.Status != "" {
		whereSQL += " AND status = ?"
		args = append(args, query.Status)
	}
	if query.Keyword != "" {
		whereSQL += " AND (title LIKE ? OR content LIKE ?)"
		keyword := "%" + query.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	return whereSQL, args
}

func scanAnnouncement(scanner messageScanner) (Announcement, error) {
	var item Announcement
	var content sql.NullString
	var coverURL sql.NullString
	var publishTime sql.NullTime
	if err := scanner.Scan(
		&item.ID,
		&item.Title,
		&content,
		&coverURL,
		&item.Status,
		&publishTime,
		&item.CreateBy,
		&item.CreateTime,
		&item.UpdateTime,
	); err != nil {
		return Announcement{}, fmt.Errorf("scan announcement: %w", err)
	}
	if content.Valid {
		item.Content = &content.String
	}
	if coverURL.Valid {
		item.CoverURL = &coverURL.String
	}
	if publishTime.Valid {
		item.PublishTime = &publishTime.Time
	}
	return item, nil
}

type messageScanner interface {
	Scan(dest ...interface{}) error
}

func checkAnnouncementRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check announcement result: %w", err)
	}
	if rowsAffected == 0 {
		return ErrAnnouncementNotFound
	}
	return nil
}

func (r *Repository) ListLogs(ctx context.Context, query LogQuery) ([]AdminLog, int, error) {
	whereSQL, args := buildLogWhere(query)

	var total int
	countSQL := `SELECT COUNT(*) FROM admin_logs` + whereSQL
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin logs: %w", err)
	}

	listSQL := `
		SELECT id, admin_id, operation_type, target_type, target_id, description, ip_address, create_time
		FROM admin_logs
	` + whereSQL + ` ORDER BY create_time DESC LIMIT ? OFFSET ?`
	args = append(args, query.PageSize, (query.Page-1)*query.PageSize)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin logs: %w", err)
	}
	defer rows.Close()

	items := make([]AdminLog, 0)
	for rows.Next() {
		var item AdminLog
		var description sql.NullString
		var ipAddress sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.AdminID,
			&item.OperationType,
			&item.TargetType,
			&item.TargetID,
			&description,
			&ipAddress,
			&item.CreateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan admin log: %w", err)
		}
		if description.Valid {
			item.Description = &description.String
		}
		if ipAddress.Valid {
			item.IPAddress = &ipAddress.String
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin logs: %w", err)
	}

	return items, total, nil
}

func buildLogWhere(query LogQuery) (string, []interface{}) {
	whereSQL := " WHERE 1 = 1"
	args := make([]interface{}, 0)

	if query.AdminID > 0 {
		whereSQL += " AND admin_id = ?"
		args = append(args, query.AdminID)
	}
	if query.OperationType != "" {
		whereSQL += " AND operation_type = ?"
		args = append(args, query.OperationType)
	}
	if query.TargetType != "" {
		whereSQL += " AND target_type = ?"
		args = append(args, query.TargetType)
	}
	if query.TargetID > 0 {
		whereSQL += " AND target_id = ?"
		args = append(args, query.TargetID)
	}
	if query.StartTime != "" {
		whereSQL += " AND create_time >= ?"
		args = append(args, query.StartTime)
	}
	if query.EndTime != "" {
		whereSQL += " AND create_time <= ?"
		args = append(args, query.EndTime)
	}

	return whereSQL, args
}
