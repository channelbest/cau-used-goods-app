package ai

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

type CreateLogInput struct {
	UserID         uint64
	ProductID      *uint64
	GenerationType string
	InputText      string
	OutputText     string
	Status         string
}

func (r *Repository) CreateLog(ctx context.Context, input CreateLogInput) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO ai_generation_logs (
			user_id,
			product_id,
			generation_type,
			input_text,
			output_text,
			status
		) VALUES (?, ?, ?, ?, ?, ?)
	`,
		input.UserID,
		input.ProductID,
		input.GenerationType,
		input.InputText,
		input.OutputText,
		input.Status,
	)
	return err
}
