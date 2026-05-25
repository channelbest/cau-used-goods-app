package review

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

func (r *Repository) Create(ctx context.Context, review *Review) error {
	query := `
		INSERT INTO reviews (order_id, product_id, reviewer_id, seller_id, rating, content)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query, review.OrderID, review.ProductID, review.ReviewerID, review.SellerID, review.Rating, review.Content)
	if err != nil {
		return fmt.Errorf("insert review: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	review.ID = uint64(id)
	return nil
}

func (r *Repository) GetByOrderAndReviewer(ctx context.Context, orderID, reviewerID uint64) (*Review, error) {
	query := `
		SELECT id, order_id, product_id, reviewer_id, seller_id, rating, content, status, create_time, is_deleted
		FROM reviews WHERE order_id = ? AND reviewer_id = ? AND is_deleted = 0
	`
	row := r.db.QueryRowContext(ctx, query, orderID, reviewerID)
	return scanReview(row)
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Review, error) {
	query := `
		SELECT id, order_id, product_id, reviewer_id, seller_id, rating, content, status, create_time, is_deleted
		FROM reviews WHERE id = ? AND is_deleted = 0
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanReview(row)
}

func (r *Repository) ListByProduct(ctx context.Context, productID uint64, page, pageSize int) ([]ReviewDetail, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM reviews WHERE product_id = ? AND is_deleted = 0 AND status = 'NORMAL'`
	if err := r.db.QueryRowContext(ctx, countQuery, productID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews: %w", err)
	}

	query := `
		SELECT r.id, r.order_id, r.product_id, r.reviewer_id, r.seller_id, r.rating, r.content, r.status, r.create_time,
			u.nickname, p.title, pi.image_url
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		LEFT JOIN products p ON p.id = r.product_id
		LEFT JOIN product_images pi ON pi.product_id = p.id AND pi.sort_order = 0
		WHERE r.product_id = ? AND r.is_deleted = 0 AND r.status = 'NORMAL'
		ORDER BY r.create_time DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, productID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()

	var items []ReviewDetail
	for rows.Next() {
		var rd ReviewDetail
		var nick, title, img sql.NullString
		err := rows.Scan(
			&rd.ID, &rd.OrderID, &rd.ProductID, &rd.ReviewerID, &rd.SellerID, &rd.Rating, &rd.Content, &rd.Status, &rd.CreateTime,
			&nick, &title, &img,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}
		if nick.Valid {
			rd.ReviewerNickname = &nick.String
		}
		rd.ProductTitle = title.String
		if img.Valid {
			rd.ProductImage = &img.String
		}
		items = append(items, rd)
	}
	return items, total, nil
}

func (r *Repository) ListBySeller(ctx context.Context, sellerID uint64, page, pageSize int) ([]ReviewDetail, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM reviews WHERE seller_id = ? AND is_deleted = 0 AND status = 'NORMAL'`
	if err := r.db.QueryRowContext(ctx, countQuery, sellerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews: %w", err)
	}

	query := `
		SELECT r.id, r.order_id, r.product_id, r.reviewer_id, r.seller_id, r.rating, r.content, r.status, r.create_time,
			u.nickname, p.title, pi.image_url
		FROM reviews r
		LEFT JOIN users u ON u.id = r.reviewer_id
		LEFT JOIN products p ON p.id = r.product_id
		LEFT JOIN product_images pi ON pi.product_id = p.id AND pi.sort_order = 0
		WHERE r.seller_id = ? AND r.is_deleted = 0 AND r.status = 'NORMAL'
		ORDER BY r.create_time DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.QueryContext(ctx, query, sellerID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()

	var items []ReviewDetail
	for rows.Next() {
		var rd ReviewDetail
		var nick, title, img sql.NullString
		err := rows.Scan(
			&rd.ID, &rd.OrderID, &rd.ProductID, &rd.ReviewerID, &rd.SellerID, &rd.Rating, &rd.Content, &rd.Status, &rd.CreateTime,
			&nick, &title, &img,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan review: %w", err)
		}
		if nick.Valid {
			rd.ReviewerNickname = &nick.String
		}
		rd.ProductTitle = title.String
		if img.Valid {
			rd.ProductImage = &img.String
		}
		items = append(items, rd)
	}
	return items, total, nil
}

func (r *Repository) GetOrderInfo(ctx context.Context, orderID uint64) (productID, buyerID, sellerID uint64, status string, err error) {
	query := `SELECT product_id, buyer_id, seller_id, status FROM orders WHERE id = ?`
	if err := r.db.QueryRowContext(ctx, query, orderID).Scan(&productID, &buyerID, &sellerID, &status); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, 0, "", fmt.Errorf("order not found")
		}
		return 0, 0, 0, "", fmt.Errorf("get order info: %w", err)
	}
	return productID, buyerID, sellerID, status, nil
}

func (r *Repository) GetAverageRating(ctx context.Context, sellerID uint64) (float64, int, error) {
	query := `SELECT AVG(rating), COUNT(*) FROM reviews WHERE seller_id = ? AND is_deleted = 0 AND status = 'NORMAL'`
	var avg sql.NullFloat64
	var count int
	if err := r.db.QueryRowContext(ctx, query, sellerID).Scan(&avg, &count); err != nil {
		return 0, 0, fmt.Errorf("get average rating: %w", err)
	}
	if !avg.Valid {
		return 0, 0, nil
	}
	return avg.Float64, count, nil
}

func scanReview(row *sql.Row) (*Review, error) {
	var r Review
	var content sql.NullString
	err := row.Scan(&r.ID, &r.OrderID, &r.ProductID, &r.ReviewerID, &r.SellerID, &r.Rating, &content, &r.Status, &r.CreateTime, &r.IsDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan review: %w", err)
	}
	if content.Valid {
		r.Content = &content.String
	}
	return &r, nil
}
