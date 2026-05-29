package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrMessageNotFound = errors.New("message not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, input CreateMessageInput) (uint64, error) {
	query := `
		INSERT INTO messages (
			receiver_id, sender_id, message_type, title, content,
			related_type, related_id, read_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		input.ReceiverID,
		input.SenderID,
		input.MessageType,
		input.Title,
		input.Content,
		input.RelatedType,
		input.RelatedID,
		ReadStatusUnread,
	)
	if err != nil {
		return 0, fmt.Errorf("create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created message id: %w", err)
	}
	return uint64(id), nil
}

func (r *Repository) ListByReceiver(ctx context.Context, receiverID uint64, readStatus string, page, pageSize int) ([]Message, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM messages WHERE receiver_id = ?`
	countArgs := []interface{}{receiverID}
	if readStatus != "" {
		countQuery += ` AND read_status = ?`
		countArgs = append(countArgs, readStatus)
	}
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count messages: %w", err)
	}

	query := `
		SELECT id, receiver_id, sender_id, message_type, title, content,
			related_type, related_id, read_status, create_time, read_time
		FROM messages
		WHERE receiver_id = ?
	`
	args := []interface{}{receiverID}
	if readStatus != "" {
		query += ` AND read_status = ?`
		args = append(args, readStatus)
	}
	query += ` ORDER BY create_time DESC LIMIT ? OFFSET ?`
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()

	items := make([]Message, 0)
	for rows.Next() {
		item, err := scanMessage(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate messages: %w", err)
	}

	return items, total, nil
}

func (r *Repository) GetByID(ctx context.Context, receiverID, messageID uint64) (*Message, error) {
	query := `
		SELECT id, receiver_id, sender_id, message_type, title, content,
			related_type, related_id, read_status, create_time, read_time
		FROM messages
		WHERE id = ? AND receiver_id = ?
	`
	item, err := scanMessage(r.db.QueryRowContext(ctx, query, messageID, receiverID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *Repository) CountUnread(ctx context.Context, receiverID uint64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND read_status = ?`
	var count int
	if err := r.db.QueryRowContext(ctx, query, receiverID, ReadStatusUnread).Scan(&count); err != nil {
		return 0, fmt.Errorf("count unread messages: %w", err)
	}
	return count, nil
}

func (r *Repository) MarkRead(ctx context.Context, receiverID, messageID uint64) error {
	query := `
		UPDATE messages
		SET read_status = ?, read_time = NOW()
		WHERE id = ? AND receiver_id = ?
	`
	result, err := r.db.ExecContext(ctx, query, ReadStatusRead, messageID, receiverID)
	if err != nil {
		return fmt.Errorf("mark message read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check message read result: %w", err)
	}
	if rowsAffected == 0 {
		return ErrMessageNotFound
	}
	return nil
}

func (r *Repository) MarkAllRead(ctx context.Context, receiverID uint64) (int64, error) {
	query := `
		UPDATE messages
		SET read_status = ?, read_time = NOW()
		WHERE receiver_id = ? AND read_status = ?
	`
	result, err := r.db.ExecContext(ctx, query, ReadStatusRead, receiverID, ReadStatusUnread)
	if err != nil {
		return 0, fmt.Errorf("mark all messages read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("check mark all messages read result: %w", err)
	}
	return rowsAffected, nil
}

type messageScanner interface {
	Scan(dest ...interface{}) error
}

func scanMessage(scanner messageScanner) (Message, error) {
	var item Message
	var senderID sql.NullInt64
	var relatedType sql.NullString
	var relatedID sql.NullInt64
	var readTime sql.NullTime

	if err := scanner.Scan(
		&item.ID,
		&item.ReceiverID,
		&senderID,
		&item.MessageType,
		&item.Title,
		&item.Content,
		&relatedType,
		&relatedID,
		&item.ReadStatus,
		&item.CreateTime,
		&readTime,
	); err != nil {
		return Message{}, fmt.Errorf("scan message: %w", err)
	}

	if senderID.Valid {
		value := uint64(senderID.Int64)
		item.SenderID = &value
	}
	if relatedType.Valid {
		value := relatedType.String
		item.RelatedType = &value
	}
	if relatedID.Valid {
		value := uint64(relatedID.Int64)
		item.RelatedID = &value
	}
	if readTime.Valid {
		value := readTime.Time
		item.ReadTime = &value
	}

	return item, nil
}
