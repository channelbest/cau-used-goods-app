package stats

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

type ProductOverview struct {
	TotalProducts    int     `json:"totalProducts"`
	OnSaleProducts   int     `json:"onSaleProducts"`
	OffShelfProducts int     `json:"offShelfProducts"`
	LockedProducts   int     `json:"lockedProducts"`
	SoldProducts     int     `json:"soldProducts"`
	DeletedProducts  int     `json:"deletedProducts"`
	TotalViews       int     `json:"totalViews"`
	TotalFavorites   int     `json:"totalFavorites"`
	AveragePrice     float64 `json:"averagePrice"`
}

type CategoryDistributionItem struct {
	CategoryID     uint64  `json:"categoryId"`
	CategoryName   string  `json:"categoryName"`
	ProductCount   int     `json:"productCount"`
	OnSaleCount    int     `json:"onSaleCount"`
	AveragePrice   float64 `json:"averagePrice"`
	TotalViews     int     `json:"totalViews"`
	TotalFavorites int     `json:"totalFavorites"`
}

type StatusDistributionItem struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type ProductTrendItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

func (r *Repository) ProductOverview(ctx context.Context) (*ProductOverview, error) {
	var overview ProductOverview

	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total_products,
			SUM(CASE WHEN status = 'ON_SALE' AND is_deleted = 0 THEN 1 ELSE 0 END) AS on_sale_products,
			SUM(CASE WHEN status = 'OFF_SHELF' AND is_deleted = 0 THEN 1 ELSE 0 END) AS off_shelf_products,
			SUM(CASE WHEN status = 'LOCKED' AND is_deleted = 0 THEN 1 ELSE 0 END) AS locked_products,
			SUM(CASE WHEN status = 'SOLD' AND is_deleted = 0 THEN 1 ELSE 0 END) AS sold_products,
			SUM(CASE WHEN status = 'DELETED' OR is_deleted = 1 THEN 1 ELSE 0 END) AS deleted_products,
			COALESCE(SUM(view_count), 0) AS total_views,
			COALESCE(SUM(favorite_count), 0) AS total_favorites,
			COALESCE(AVG(CASE WHEN is_deleted = 0 THEN price ELSE NULL END), 0) AS average_price
		FROM products
	`).Scan(
		&overview.TotalProducts,
		&overview.OnSaleProducts,
		&overview.OffShelfProducts,
		&overview.LockedProducts,
		&overview.SoldProducts,
		&overview.DeletedProducts,
		&overview.TotalViews,
		&overview.TotalFavorites,
		&overview.AveragePrice,
	)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

func (r *Repository) CategoryDistribution(ctx context.Context) ([]CategoryDistributionItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			c.id,
			c.name,
			COUNT(p.id) AS product_count,
			COALESCE(SUM(CASE WHEN p.status = 'ON_SALE' AND p.is_deleted = 0 THEN 1 ELSE 0 END), 0) AS on_sale_count,
			COALESCE(AVG(CASE WHEN p.is_deleted = 0 THEN p.price ELSE NULL END), 0) AS average_price,
			COALESCE(SUM(CASE WHEN p.is_deleted = 0 THEN p.view_count ELSE 0 END), 0) AS total_views,
			COALESCE(SUM(CASE WHEN p.is_deleted = 0 THEN p.favorite_count ELSE 0 END), 0) AS total_favorites
		FROM categories c
		LEFT JOIN products p ON p.category_id = c.id
		WHERE c.status = 'ENABLED'
		GROUP BY c.id, c.name
		ORDER BY product_count DESC, c.sort_order ASC, c.id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []CategoryDistributionItem
	for rows.Next() {
		var item CategoryDistributionItem
		if err := rows.Scan(
			&item.CategoryID,
			&item.CategoryName,
			&item.ProductCount,
			&item.OnSaleCount,
			&item.AveragePrice,
			&item.TotalViews,
			&item.TotalFavorites,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, rows.Err()
}

func (r *Repository) StatusDistribution(ctx context.Context) ([]StatusDistributionItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT status, COUNT(*) AS count
		FROM products
		GROUP BY status
		ORDER BY count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []StatusDistributionItem
	for rows.Next() {
		var item StatusDistributionItem
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, rows.Err()
}

func (r *Repository) ProductTrend(ctx context.Context, days int) ([]ProductTrendItem, error) {
	if days <= 0 {
		days = 7
	}
	if days > 90 {
		days = 90
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT DATE(create_time) AS date, COUNT(*) AS count
		FROM products
		WHERE create_time >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(create_time)
		ORDER BY date ASC
	`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ProductTrendItem
	for rows.Next() {
		var item ProductTrendItem
		if err := rows.Scan(&item.Date, &item.Count); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, rows.Err()
}
