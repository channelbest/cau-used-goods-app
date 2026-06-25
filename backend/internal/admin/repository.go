package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"cau-used-goods-app/backend/internal/message"
)

var ErrAnnouncementNotFound = errors.New("announcement not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAnnouncement(ctx context.Context, input CreateAnnouncementInput) (uint64, error) {
	return r.CreateAnnouncementTx(ctx, nil, input)
}

func (r *Repository) CreateAnnouncementTx(ctx context.Context, tx *sql.Tx, input CreateAnnouncementInput) (uint64, error) {
	query := `
		INSERT INTO announcements (title, content, cover_url, status, publish_time, create_by)
		VALUES (?, ?, ?, ?, CASE WHEN ? = 'PUBLISHED' THEN NOW() ELSE NULL END, ?)
	`
	execer := announcementExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, query, input.Title, input.Content, input.CoverURL, input.Status, input.Status, input.AdminID)
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
	return r.UpdateAnnouncementTx(ctx, nil, input)
}

func (r *Repository) UpdateAnnouncementTx(ctx context.Context, tx *sql.Tx, input UpdateAnnouncementInput) error {
	query := `
		UPDATE announcements
		SET title = ?, content = ?, cover_url = ?
		WHERE id = ?
	`
	execer := announcementExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, query, input.Title, input.Content, input.CoverURL, input.ID)
	if err != nil {
		return fmt.Errorf("update announcement: %w", err)
	}
	return r.checkAnnouncementRowsAffected(ctx, input.ID, result)
}

func (r *Repository) UpdateAnnouncementStatus(ctx context.Context, id uint64, status string) error {
	return r.UpdateAnnouncementStatusTx(ctx, nil, id, status)
}

func (r *Repository) UpdateAnnouncementStatusTx(ctx context.Context, tx *sql.Tx, id uint64, status string) error {
	query := `
		UPDATE announcements
		SET status = ?, publish_time = CASE WHEN ? = 'PUBLISHED' THEN NOW() ELSE NULL END
		WHERE id = ?
	`
	execer := announcementExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, query, status, status, id)
	if err != nil {
		return fmt.Errorf("update announcement status: %w", err)
	}
	return r.checkAnnouncementRowsAffected(ctx, id, result)
}

func (r *Repository) BroadcastPublishedAnnouncementTx(ctx context.Context, tx *sql.Tx, announcementID, senderID uint64) (int64, error) {
	if tx == nil {
		return 0, fmt.Errorf("broadcast announcement requires transaction")
	}
	result, err := tx.ExecContext(ctx, `
		INSERT INTO messages (
			receiver_id, sender_id, message_type, title, content,
			related_type, related_id, read_status
		)
		SELECT
			u.id,
			?,
			?,
			CONCAT('平台公告：', a.title),
			COALESCE(NULLIF(a.content, ''), a.title),
			?,
			a.id,
			?
		FROM users u
		JOIN announcements a ON a.id = ?
		WHERE a.status = 'PUBLISHED'
		  AND u.is_deleted = 0
		  AND u.account_status <> 'CANCELED'
		  AND u.role NOT IN ('ADMIN', 'SUPER_ADMIN')
		  AND NOT EXISTS (
			SELECT 1
			FROM messages m
			WHERE m.receiver_id = u.id
			  AND m.message_type = ?
			  AND m.related_type = ?
			  AND m.related_id = a.id
		  )
	`,
		senderID,
		message.MessageTypeSystemNotice,
		message.RelatedTypeNotice,
		message.ReadStatusUnread,
		announcementID,
		message.MessageTypeSystemNotice,
		message.RelatedTypeNotice,
	)
	if err != nil {
		return 0, fmt.Errorf("broadcast announcement message: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("check broadcast announcement result: %w", err)
	}
	return affected, nil
}

type announcementExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (r *Repository) CreateLog(ctx context.Context, input LogActionInput) (uint64, error) {
	return createLog(ctx, r.db, input)
}

func (r *Repository) CreateLogTx(ctx context.Context, tx *sql.Tx, input LogActionInput) (uint64, error) {
	return createLog(ctx, tx, input)
}

type logExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func createLog(ctx context.Context, execer logExecutor, input LogActionInput) (uint64, error) {
	query := `
		INSERT INTO admin_logs (
			admin_id, operation_type, target_type, target_id, description, ip_address, related_type, related_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := execer.ExecContext(
		ctx,
		query,
		input.AdminID,
		input.OperationType,
		input.TargetType,
		input.TargetID,
		input.Description,
		input.IPAddress,
		input.RelatedType,
		input.RelatedID,
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

func (r *Repository) checkAnnouncementRowsAffected(ctx context.Context, id uint64, result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check announcement result: %w", err)
	}
	if rowsAffected == 0 {
		exists, err := r.announcementExists(ctx, id)
		if err != nil {
			return err
		}
		if !exists {
			return ErrAnnouncementNotFound
		}
	}
	return nil
}

func (r *Repository) announcementExists(ctx context.Context, id uint64) (bool, error) {
	query := `SELECT COUNT(*) FROM announcements WHERE id = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&count); err != nil {
		return false, fmt.Errorf("check announcement exists: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) GetLogByID(ctx context.Context, id uint64) (*AdminLog, error) {
	query := `
		SELECT l.id, l.admin_id, l.operation_type, l.target_type, l.target_id, l.description, l.ip_address,
		       l.related_type, l.related_id, l.create_time,
		       u.nickname, u.avatar_url
		FROM admin_logs l
		LEFT JOIN users u ON u.id = l.admin_id
		WHERE l.id = ?
	`
	var item AdminLog
	var description sql.NullString
	var ipAddress sql.NullString
	var relatedType sql.NullString
	var relatedID sql.NullInt64
	var adminName sql.NullString
	var adminAvatar sql.NullString
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.AdminID,
		&item.OperationType,
		&item.TargetType,
		&item.TargetID,
		&description,
		&ipAddress,
		&relatedType,
		&relatedID,
		&item.CreateTime,
		&adminName,
		&adminAvatar,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get admin log by id: %w", err)
	}
	if description.Valid {
		item.Description = &description.String
	}
	if ipAddress.Valid {
		item.IPAddress = &ipAddress.String
	}
	if relatedType.Valid {
		item.RelatedType = &relatedType.String
	}
	if relatedID.Valid {
		value := uint64(relatedID.Int64)
		item.RelatedID = &value
	}
	if adminName.Valid {
		item.AdminName = adminName.String
	}
	if adminAvatar.Valid {
		item.AdminAvatar = adminAvatar.String
	}
	return &item, nil
}

func (r *Repository) ListLogs(ctx context.Context, query LogQuery) ([]AdminLog, int, error) {
	whereSQL, args := buildLogWhere(query)

	var total int
	countSQL := `SELECT COUNT(*) FROM admin_logs l` + whereSQL
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin logs: %w", err)
	}

	listSQL := `
		SELECT l.id, l.admin_id, l.operation_type, l.target_type, l.target_id, l.description, l.ip_address,
		       l.related_type, l.related_id, l.create_time,
		       a.nickname, a.avatar_url,
		       CASE
		         WHEN l.target_type = 'USER' THEN u.nickname
		         WHEN l.target_type = 'PRODUCT' THEN p.title
		         ELSE NULL
		       END as target_name,
		       CASE
		         WHEN l.target_type = 'USER' THEN u.college
		         ELSE NULL
		       END as target_college,
		       CASE
		         WHEN l.target_type = 'REPORT' THEN r.reason_type
		         ELSE NULL
		       END as report_reason,
		       CASE
		         WHEN l.target_type = 'REPORT' THEN ru.nickname
		         ELSE NULL
		       END as reporter_name
		FROM admin_logs l
		LEFT JOIN users a ON a.id = l.admin_id
		LEFT JOIN users u ON u.id = l.target_id AND l.target_type = 'USER'
		LEFT JOIN products p ON p.id = l.target_id AND l.target_type = 'PRODUCT'
		LEFT JOIN reports r ON r.id = l.target_id AND l.target_type = 'REPORT'
		LEFT JOIN users ru ON ru.id = r.reporter_id AND l.target_type = 'REPORT'
	` + whereSQL + ` ORDER BY l.create_time DESC LIMIT ? OFFSET ?`
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
		var relatedType sql.NullString
		var relatedID sql.NullInt64
		var adminName sql.NullString
		var adminAvatar sql.NullString
		var targetName sql.NullString
		var targetCollege sql.NullString
		var reportReason sql.NullString
		var reporterName sql.NullString
		if err := rows.Scan(
			&item.ID,
			&item.AdminID,
			&item.OperationType,
			&item.TargetType,
			&item.TargetID,
			&description,
			&ipAddress,
			&relatedType,
			&relatedID,
			&item.CreateTime,
			&adminName,
			&adminAvatar,
			&targetName,
			&targetCollege,
			&reportReason,
			&reporterName,
		); err != nil {
			return nil, 0, fmt.Errorf("scan admin log: %w", err)
		}
		if description.Valid {
			item.Description = &description.String
		}
		if ipAddress.Valid {
			item.IPAddress = &ipAddress.String
		}
		if relatedType.Valid {
			item.RelatedType = &relatedType.String
		}
		if relatedID.Valid {
			value := uint64(relatedID.Int64)
			item.RelatedID = &value
		}
		if adminName.Valid {
			item.AdminName = adminName.String
		}
		if adminAvatar.Valid {
			item.AdminAvatar = adminAvatar.String
		}
		if targetName.Valid {
			item.TargetName = targetName.String
		}
		if targetCollege.Valid {
			item.TargetCollege = targetCollege.String
		}
		if reportReason.Valid {
			item.ReportReason = reportReason.String
		}
		if reporterName.Valid {
			item.ReporterName = reporterName.String
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
		whereSQL += " AND l.admin_id = ?"
		args = append(args, query.AdminID)
	}
	if query.OperationType != "" {
		whereSQL += " AND l.operation_type = ?"
		args = append(args, query.OperationType)
	}
	if query.TargetType != "" {
		whereSQL += " AND l.target_type = ?"
		args = append(args, query.TargetType)
	}
	if query.TargetID > 0 {
		whereSQL += " AND l.target_id = ?"
		args = append(args, query.TargetID)
	}
	if query.StartTime != "" {
		whereSQL += " AND l.create_time >= ?"
		args = append(args, query.StartTime)
	}
	if query.EndTime != "" {
		whereSQL += " AND l.create_time <= ?"
		args = append(args, query.EndTime)
	}

	return whereSQL, args
}
