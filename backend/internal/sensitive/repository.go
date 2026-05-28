package sensitive

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type SensitiveWord struct {
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
