package sensitive

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrSensitiveWordNotFound = errors.New("sensitive word not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type SensitiveWord struct {
	ID         uint64    `json:"id"`
	Word       string    `json:"word"`
	WordType   string    `json:"wordType"`
	Status     string    `json:"status"`
	CreateBy   *uint64   `json:"createBy,omitempty"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type WordQuery struct {
	Status   string
	WordType string
	Keyword  string
	Page     int
	PageSize int
}

type CreateWordInput struct {
	Word     string
	WordType string
	Status   string
	CreateBy uint64
}

type UpdateWordInput struct {
	ID       uint64
	Word     string
	WordType string
	Status   string
}

func (r *Repository) ListEnabledWords(ctx context.Context) ([]SensitiveWord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, word, word_type, status
		FROM sensitive_words
		WHERE status = 'ENABLED'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []SensitiveWord
	for rows.Next() {
		var w SensitiveWord
		if err := rows.Scan(&w.ID, &w.Word, &w.WordType, &w.Status); err != nil {
			return nil, err
		}
		list = append(list, w)
	}

	return list, rows.Err()
}

func (r *Repository) ListWords(ctx context.Context, query WordQuery) ([]SensitiveWord, int, error) {
	whereSQL, args := buildWordWhere(query)

	var total int
	countSQL := `SELECT COUNT(*) FROM sensitive_words` + whereSQL
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count sensitive words: %w", err)
	}

	listSQL := `
		SELECT id, word, word_type, status, create_by, create_time, update_time
		FROM sensitive_words
	` + whereSQL + ` ORDER BY update_time DESC LIMIT ? OFFSET ?`
	args = append(args, query.PageSize, (query.Page-1)*query.PageSize)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list sensitive words: %w", err)
	}
	defer rows.Close()

	items := make([]SensitiveWord, 0)
	for rows.Next() {
		item, err := scanSensitiveWord(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate sensitive words: %w", err)
	}

	return items, total, nil
}

func (r *Repository) CreateWord(ctx context.Context, input CreateWordInput) (uint64, error) {
	query := `
		INSERT INTO sensitive_words (word, word_type, status, create_by)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, input.Word, input.WordType, input.Status, input.CreateBy)
	if err != nil {
		return 0, fmt.Errorf("create sensitive word: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created sensitive word id: %w", err)
	}
	return uint64(id), nil
}

func (r *Repository) UpdateWord(ctx context.Context, input UpdateWordInput) error {
	query := `
		UPDATE sensitive_words
		SET word = ?, word_type = ?, status = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, query, input.Word, input.WordType, input.Status, input.ID)
	if err != nil {
		return fmt.Errorf("update sensitive word: %w", err)
	}
	return checkRowsAffected(result)
}

func (r *Repository) DisableWord(ctx context.Context, id uint64) error {
	query := `UPDATE sensitive_words SET status = 'DISABLED' WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("disable sensitive word: %w", err)
	}
	return checkRowsAffected(result)
}

func buildWordWhere(query WordQuery) (string, []interface{}) {
	whereSQL := " WHERE 1 = 1"
	args := make([]interface{}, 0)

	if query.Status != "" {
		whereSQL += " AND status = ?"
		args = append(args, query.Status)
	}
	if query.WordType != "" {
		whereSQL += " AND word_type = ?"
		args = append(args, query.WordType)
	}
	if query.Keyword != "" {
		whereSQL += " AND word LIKE ?"
		args = append(args, "%"+query.Keyword+"%")
	}

	return whereSQL, args
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanSensitiveWord(scanner scanner) (SensitiveWord, error) {
	var item SensitiveWord
	var createBy sql.NullInt64
	if err := scanner.Scan(
		&item.ID,
		&item.Word,
		&item.WordType,
		&item.Status,
		&createBy,
		&item.CreateTime,
		&item.UpdateTime,
	); err != nil {
		return SensitiveWord{}, fmt.Errorf("scan sensitive word: %w", err)
	}
	if createBy.Valid {
		value := uint64(createBy.Int64)
		item.CreateBy = &value
	}
	return item, nil
}

func checkRowsAffected(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check sensitive word result: %w", err)
	}
	if rowsAffected == 0 {
		return ErrSensitiveWordNotFound
	}
	return nil
}
