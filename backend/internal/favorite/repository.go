package favorite

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

func (r *Repository) Add(ctx context.Context, userID, productID uint64) error {
	query := `
		INSERT INTO favorites (user_id, product_id)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE is_deleted = 0
	`
	_, err := r.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return fmt.Errorf("add favorite: %w", err)
	}
	return nil
}

func (r *Repository) Remove(ctx context.Context, userID, productID uint64) error {
	query := `UPDATE favorites SET is_deleted = 1 WHERE user_id = ? AND product_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID, productID)
	if err != nil {
		return fmt.Errorf("remove favorite: %w", err)
	}
	return nil
}

func (r *Repository) IsFavorited(ctx context.Context, userID, productID uint64) (bool, error) {
	query := `SELECT COUNT(*) FROM favorites WHERE user_id = ? AND product_id = ? AND is_deleted = 0`
	var count int
	if err := r.db.QueryRowContext(ctx, query, userID, productID).Scan(&count); err != nil {
		return false, fmt.Errorf("check favorite: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) ListByUser(ctx context.Context, userID uint64, page, pageSize int) ([]FavoriteDetail, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM favorites WHERE user_id = ? AND is_deleted = 0`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count favorites: %w", err)
	}

	query := `
		SELECT f.id, f.user_id, f.product_id, f.create_time,
			p.title, p.price, p.status, pi.image_url, u.nickname
		FROM favorites f
		JOIN products p ON p.id = f.product_id
		LEFT JOIN product_images pi ON pi.product_id = p.id AND pi.sort_order = 0
		LEFT JOIN users u ON u.id = p.seller_id
		WHERE f.user_id = ? AND f.is_deleted = 0
		ORDER BY f.create_time DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list favorites: %w", err)
	}
	defer rows.Close()

	var items []FavoriteDetail
	for rows.Next() {
		var fd FavoriteDetail
		var img, nick sql.NullString
		err := rows.Scan(&fd.ID, &fd.UserID, &fd.ProductID, &fd.CreateTime,
			&fd.ProductTitle, &fd.ProductPrice, &fd.ProductStatus, &img, &nick)
		if err != nil {
			return nil, 0, fmt.Errorf("scan favorite: %w", err)
		}
		if img.Valid {
			fd.ProductImage = &img.String
		}
		if nick.Valid {
			fd.SellerNickname = &nick.String
		}
		items = append(items, fd)
	}
	return items, total, nil
}
