package appeal

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

func (r *Repository) Create(ctx context.Context, input CreateAppealInput) (*Appeal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create appeal tx: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		INSERT INTO appeals (appellant_id, target_type, target_id, reason, status)
		VALUES (?, ?, ?, ?, ?)
	`, input.AppellantID, input.TargetType, input.TargetID, input.Reason, StatusPending)
	if err != nil {
		return nil, fmt.Errorf("create appeal: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get created appeal id: %w", err)
	}
	appealID := uint64(id)

	for i, imageURL := range input.EvidenceURLs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO appeal_images (appeal_id, image_url, sort_order)
			VALUES (?, ?, ?)
		`, appealID, imageURL, i); err != nil {
			return nil, fmt.Errorf("create appeal image: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create appeal tx: %w", err)
	}
	return r.GetByID(ctx, appealID)
}

func (r *Repository) HasActiveAppeal(ctx context.Context, appellantID uint64, targetType string, targetID uint64) (bool, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM appeals
		WHERE appellant_id = ? AND target_type = ? AND target_id = ?
		  AND status IN (?, ?)
	`, appellantID, targetType, targetID, StatusPending, StatusProcessing).Scan(&count); err != nil {
		return false, fmt.Errorf("check active appeal: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) TargetExists(ctx context.Context, targetType string, targetID uint64) (bool, error) {
	var query string
	switch targetType {
	case TargetTypeProduct:
		query = "SELECT COUNT(*) FROM products WHERE id = ? AND is_deleted = 0"
	case TargetTypeUser:
		query = "SELECT COUNT(*) FROM users WHERE id = ? AND is_deleted = 0"
	case TargetTypeOrder:
		query = "SELECT COUNT(*) FROM orders WHERE id = ?"
	case TargetTypeReport:
		query = "SELECT COUNT(*) FROM reports WHERE id = ?"
	default:
		return false, nil
	}

	var count int
	if err := r.db.QueryRowContext(ctx, query, targetID).Scan(&count); err != nil {
		return false, fmt.Errorf("check appeal target: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) TargetAppealableBy(ctx context.Context, appellantID uint64, targetType string, targetID uint64) (bool, error) {
	var query string
	args := []interface{}{targetID, appellantID}
	switch targetType {
	case TargetTypeProduct:
		query = "SELECT COUNT(*) FROM products WHERE id = ? AND seller_id = ? AND is_deleted = 0"
	case TargetTypeUser:
		query = "SELECT COUNT(*) FROM users WHERE id = ? AND id = ? AND is_deleted = 0"
	case TargetTypeOrder:
		query = "SELECT COUNT(*) FROM orders WHERE id = ? AND (buyer_id = ? OR seller_id = ?)"
		args = append(args, appellantID)
	case TargetTypeReport:
		query = "SELECT COUNT(*) FROM reports WHERE id = ? AND reporter_id = ?"
	default:
		return false, nil
	}

	var count int
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return false, fmt.Errorf("check appeal target permission: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Appeal, error) {
	query := `
		SELECT id, appellant_id, target_type, target_id, reason, status,
			handle_result, handler_id, DATE_FORMAT(handle_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s')
		FROM appeals
		WHERE id = ?
	`
	item, err := scanAppeal(r.db.QueryRowContext(ctx, query, id))
	if err != nil || item == nil {
		return item, err
	}
	images, err := r.ListImages(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	item.Images = images
	return item, nil
}

func (r *Repository) GetDetailByID(ctx context.Context, id uint64) (*AppealDetail, error) {
	query := `
		SELECT a.id, a.appellant_id, a.target_type, a.target_id, a.reason, a.status,
			a.handle_result, a.handler_id, DATE_FORMAT(a.handle_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(a.create_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(a.update_time, '%Y-%m-%d %H:%i:%s'),
			appellant.nickname, handler.nickname
		FROM appeals a
		LEFT JOIN users appellant ON appellant.id = a.appellant_id
		LEFT JOIN users handler ON handler.id = a.handler_id
		WHERE a.id = ?
	`

	var item AppealDetail
	var handleResult, handleTime sql.NullString
	var handlerID sql.NullInt64
	var appellantNickname, handlerNickname sql.NullString
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.AppellantID, &item.TargetType, &item.TargetID, &item.Reason,
		&item.Status, &handleResult, &handlerID, &handleTime, &item.CreateTime,
		&item.UpdateTime, &appellantNickname, &handlerNickname,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan appeal detail: %w", err)
	}
	fillAppealNullableFields(&item.Appeal, handleResult, handlerID, handleTime)
	if appellantNickname.Valid {
		item.AppellantNickname = &appellantNickname.String
	}
	if handlerNickname.Valid {
		item.HandlerNickname = &handlerNickname.String
	}

	images, err := r.ListImages(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	item.Images = images
	return &item, nil
}

func (r *Repository) List(ctx context.Context, query AppealQuery) ([]AppealDetail, int, error) {
	where := " WHERE 1 = 1 "
	args := []interface{}{}
	if query.AppellantID > 0 {
		where += " AND a.appellant_id = ? "
		args = append(args, query.AppellantID)
	}
	if query.TargetType != "" {
		where += " AND a.target_type = ? "
		args = append(args, query.TargetType)
	}
	if query.TargetID > 0 {
		where += " AND a.target_id = ? "
		args = append(args, query.TargetID)
	}
	if query.Status != "" {
		where += " AND a.status = ? "
		args = append(args, query.Status)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM appeals a"+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count appeals: %w", err)
	}

	sqlText := `
		SELECT a.id, a.appellant_id, a.target_type, a.target_id, a.reason, a.status,
			a.handle_result, a.handler_id, DATE_FORMAT(a.handle_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(a.create_time, '%Y-%m-%d %H:%i:%s'),
			DATE_FORMAT(a.update_time, '%Y-%m-%d %H:%i:%s'),
			appellant.nickname, handler.nickname
		FROM appeals a
		LEFT JOIN users appellant ON appellant.id = a.appellant_id
		LEFT JOIN users handler ON handler.id = a.handler_id
	` + where + `
		ORDER BY a.create_time DESC, a.id DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, query.PageSize, (query.Page-1)*query.PageSize)
	rows, err := r.db.QueryContext(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list appeals: %w", err)
	}
	defer rows.Close()

	items := make([]AppealDetail, 0)
	for rows.Next() {
		var item AppealDetail
		var handleResult, handleTime sql.NullString
		var handlerID sql.NullInt64
		var appellantNickname, handlerNickname sql.NullString
		if err := rows.Scan(
			&item.ID, &item.AppellantID, &item.TargetType, &item.TargetID, &item.Reason,
			&item.Status, &handleResult, &handlerID, &handleTime, &item.CreateTime,
			&item.UpdateTime, &appellantNickname, &handlerNickname,
		); err != nil {
			return nil, 0, fmt.Errorf("scan appeal list item: %w", err)
		}
		fillAppealNullableFields(&item.Appeal, handleResult, handlerID, handleTime)
		if appellantNickname.Valid {
			item.AppellantNickname = &appellantNickname.String
		}
		if handlerNickname.Valid {
			item.HandlerNickname = &handlerNickname.String
		}
		images, err := r.ListImages(ctx, item.ID)
		if err != nil {
			return nil, 0, err
		}
		item.Images = images
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate appeals: %w", err)
	}
	return items, total, nil
}

func (r *Repository) ListImages(ctx context.Context, appealID uint64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT image_url
		FROM appeal_images
		WHERE appeal_id = ?
		ORDER BY sort_order ASC, id ASC
	`, appealID)
	if err != nil {
		return nil, fmt.Errorf("list appeal images: %w", err)
	}
	defer rows.Close()

	images := make([]string, 0)
	for rows.Next() {
		var imageURL string
		if err := rows.Scan(&imageURL); err != nil {
			return nil, fmt.Errorf("scan appeal image: %w", err)
		}
		images = append(images, imageURL)
	}
	return images, rows.Err()
}

func (r *Repository) Handle(ctx context.Context, input HandleAppealInput) (*Appeal, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin handle appeal tx: %w", err)
	}
	defer tx.Rollback()
	if err := r.HandleTx(ctx, tx, input); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit handle appeal tx: %w", err)
	}
	return r.GetByID(ctx, input.AppealID)
}

func (r *Repository) HandleTx(ctx context.Context, tx *sql.Tx, input HandleAppealInput) error {
	var currentStatus, targetType string
	var targetID uint64
	if err := tx.QueryRowContext(ctx, `
		SELECT status, target_type, target_id
		FROM appeals
		WHERE id = ?
		FOR UPDATE
	`, input.AppealID).Scan(&currentStatus, &targetType, &targetID); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("appeal not found")
		}
		return fmt.Errorf("lock appeal: %w", err)
	}
	if currentStatus != StatusPending && currentStatus != StatusProcessing {
		return fmt.Errorf("appeal already handled")
	}
	if currentStatus != StatusProcessing {
		return fmt.Errorf("appeal must be marked processing before handle")
	}
	if input.Status == StatusApproved {
		if err := r.validateApprovedActionTx(ctx, tx, targetType, targetID); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE appeals
		SET status = ?, handle_result = ?, handler_id = ?, handle_time = NOW(), update_time = CURRENT_TIMESTAMP
		WHERE id = ?
	`, input.Status, input.HandleResult, input.AdminID, input.AppealID); err != nil {
		return fmt.Errorf("handle appeal: %w", err)
	}
	return nil
}

func (r *Repository) validateApprovedActionTx(ctx context.Context, tx *sql.Tx, targetType string, targetID uint64) error {
	if targetType != TargetTypeUser {
		return nil
	}
	var accountStatus string
	if err := tx.QueryRowContext(ctx, `
		SELECT account_status
		FROM users
		WHERE id = ? AND is_deleted = 0
		FOR UPDATE
	`, targetID).Scan(&accountStatus); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("appeal target not found")
		}
		return fmt.Errorf("check appealed user status: %w", err)
	}
	if accountStatus != "NORMAL" {
		return fmt.Errorf("user appeal target must be restored before approval")
	}
	return nil
}

func (r *Repository) Close(ctx context.Context, input CloseAppealInput) (*Appeal, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE appeals
		SET status = ?, handle_result = ?, handler_id = NULL, handle_time = NOW(), update_time = CURRENT_TIMESTAMP
		WHERE id = ? AND appellant_id = ? AND status = ?
	`, StatusClosed, input.CloseReason, input.AppealID, input.AppellantID, StatusPending)
	if err != nil {
		return nil, fmt.Errorf("close appeal: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("check close appeal result: %w", err)
	}
	if affected == 0 {
		return nil, fmt.Errorf("appeal cannot be closed")
	}
	return r.GetByID(ctx, input.AppealID)
}

func (r *Repository) MarkProcessing(ctx context.Context, appealID uint64, adminID uint64) (*Appeal, error) {
	if err := r.MarkProcessingTx(ctx, nil, appealID, adminID); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, appealID)
}

func (r *Repository) MarkProcessingTx(ctx context.Context, tx *sql.Tx, appealID uint64, adminID uint64) error {
	execer := appealExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, `
		UPDATE appeals
		SET status = ?, handler_id = ?, update_time = CURRENT_TIMESTAMP
		WHERE id = ? AND status = ?
	`, StatusProcessing, adminID, appealID, StatusPending)
	if err != nil {
		return fmt.Errorf("mark appeal processing: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check mark appeal processing result: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("appeal cannot be marked processing")
	}
	return nil
}

type appealExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (r *Repository) applyApprovedAction(ctx context.Context, tx *sql.Tx, item *Appeal) error {
	switch item.TargetType {
	case TargetTypeProduct:
		_, err := tx.ExecContext(ctx, `
			UPDATE products
			SET status = 'ON_SALE', off_shelf_reason = NULL, update_time = CURRENT_TIMESTAMP
			WHERE id = ? AND is_deleted = 0 AND status = 'OFF_SHELF'
		`, item.TargetID)
		if err != nil {
			return fmt.Errorf("restore appealed product: %w", err)
		}
	case TargetTypeUser:
		// User account restoration is executed through the user-management
		// status API after the appeal is approved.
		return nil
	case TargetTypeOrder, TargetTypeReport:
		return nil
	}
	return nil
}

func scanAppeal(scanner interface {
	Scan(dest ...interface{}) error
}) (*Appeal, error) {
	var item Appeal
	var handleResult, handleTime sql.NullString
	var handlerID sql.NullInt64
	if err := scanner.Scan(
		&item.ID, &item.AppellantID, &item.TargetType, &item.TargetID, &item.Reason,
		&item.Status, &handleResult, &handlerID, &handleTime, &item.CreateTime, &item.UpdateTime,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan appeal: %w", err)
	}
	fillAppealNullableFields(&item, handleResult, handlerID, handleTime)
	return &item, nil
}

func fillAppealNullableFields(item *Appeal, handleResult sql.NullString, handlerID sql.NullInt64, handleTime sql.NullString) {
	if handleResult.Valid {
		item.HandleResult = &handleResult.String
	}
	if handlerID.Valid {
		value := uint64(handlerID.Int64)
		item.HandlerID = &value
	}
	if handleTime.Valid {
		item.HandleTime = &handleTime.String
	}
}
