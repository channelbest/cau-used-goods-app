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

type OrderOverview struct {
	TotalOrders           int     `json:"totalOrders"`
	PendingConfirmOrders  int     `json:"pendingConfirmOrders"`
	WaitMeetOrders        int     `json:"waitMeetOrders"`
	CompletedOrders       int     `json:"completedOrders"`
	CanceledOrders        int     `json:"canceledOrders"`
	ExceptionClosedOrders int     `json:"exceptionClosedOrders"`
	TotalCompletedAmount  float64 `json:"totalCompletedAmount"`
	AverageCompletedPrice float64 `json:"averageCompletedPrice"`
}

type UserOverview struct {
	TotalUsers      int `json:"totalUsers"`
	NormalUsers     int `json:"normalUsers"`
	DisabledUsers   int `json:"disabledUsers"`
	VerifiedUsers   int `json:"verifiedUsers"`
	PendingUsers    int `json:"pendingUsers"`
	UnverifiedUsers int `json:"unverifiedUsers"`
	AdminUsers      int `json:"adminUsers"`
}

type ReportOverview struct {
	TotalReports      int `json:"totalReports"`
	PendingReports    int `json:"pendingReports"`
	ProcessingReports int `json:"processingReports"`
	ResolvedReports   int `json:"resolvedReports"`
	RejectedReports   int `json:"rejectedReports"`
	ClosedReports     int `json:"closedReports"`
	ProductReports    int `json:"productReports"`
	UserReports       int `json:"userReports"`
	OrderReports      int `json:"orderReports"`
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

func (r *Repository) OrderOverview(ctx context.Context) (*OrderOverview, error) {
	var overview OrderOverview
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total_orders,
			COALESCE(SUM(CASE WHEN status = 'PENDING_CONFIRM' THEN 1 ELSE 0 END), 0) AS pending_confirm_orders,
			COALESCE(SUM(CASE WHEN status = 'WAIT_MEET' THEN 1 ELSE 0 END), 0) AS wait_meet_orders,
			COALESCE(SUM(CASE WHEN status = 'COMPLETED' THEN 1 ELSE 0 END), 0) AS completed_orders,
			COALESCE(SUM(CASE WHEN status = 'CANCELED' THEN 1 ELSE 0 END), 0) AS canceled_orders,
			COALESCE(SUM(CASE WHEN status = 'EXCEPTION_CLOSED' THEN 1 ELSE 0 END), 0) AS exception_closed_orders,
			COALESCE(SUM(CASE WHEN status = 'COMPLETED' THEN product_price_snapshot ELSE 0 END), 0) AS total_completed_amount,
			COALESCE(AVG(CASE WHEN status = 'COMPLETED' THEN product_price_snapshot ELSE NULL END), 0) AS average_completed_price
		FROM orders
	`).Scan(
		&overview.TotalOrders,
		&overview.PendingConfirmOrders,
		&overview.WaitMeetOrders,
		&overview.CompletedOrders,
		&overview.CanceledOrders,
		&overview.ExceptionClosedOrders,
		&overview.TotalCompletedAmount,
		&overview.AverageCompletedPrice,
	)
	if err != nil {
		return nil, err
	}
	return &overview, nil
}

func (r *Repository) UserOverview(ctx context.Context) (*UserOverview, error) {
	var overview UserOverview
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total_users,
			COALESCE(SUM(CASE WHEN account_status = 'NORMAL' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS normal_users,
			COALESCE(SUM(CASE WHEN account_status = 'DISABLED' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS disabled_users,
			COALESCE(SUM(CASE WHEN auth_status = 'VERIFIED' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS verified_users,
			COALESCE(SUM(CASE WHEN auth_status = 'PENDING' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS pending_users,
			COALESCE(SUM(CASE WHEN auth_status = 'UNVERIFIED' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS unverified_users,
			COALESCE(SUM(CASE WHEN role = 'ADMIN' AND is_deleted = 0 THEN 1 ELSE 0 END), 0) AS admin_users
		FROM users
	`).Scan(
		&overview.TotalUsers,
		&overview.NormalUsers,
		&overview.DisabledUsers,
		&overview.VerifiedUsers,
		&overview.PendingUsers,
		&overview.UnverifiedUsers,
		&overview.AdminUsers,
	)
	if err != nil {
		return nil, err
	}
	return &overview, nil
}

func (r *Repository) ReportOverview(ctx context.Context) (*ReportOverview, error) {
	var overview ReportOverview
	err := r.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) AS total_reports,
			COALESCE(SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_reports,
			COALESCE(SUM(CASE WHEN status = 'PROCESSING' THEN 1 ELSE 0 END), 0) AS processing_reports,
			COALESCE(SUM(CASE WHEN status = 'RESOLVED' THEN 1 ELSE 0 END), 0) AS resolved_reports,
			COALESCE(SUM(CASE WHEN status = 'REJECTED' THEN 1 ELSE 0 END), 0) AS rejected_reports,
			COALESCE(SUM(CASE WHEN status = 'CLOSED' THEN 1 ELSE 0 END), 0) AS closed_reports,
			COALESCE(SUM(CASE WHEN target_type = 'PRODUCT' THEN 1 ELSE 0 END), 0) AS product_reports,
			COALESCE(SUM(CASE WHEN target_type = 'USER' THEN 1 ELSE 0 END), 0) AS user_reports,
			COALESCE(SUM(CASE WHEN target_type = 'ORDER' THEN 1 ELSE 0 END), 0) AS order_reports
		FROM reports
	`).Scan(
		&overview.TotalReports,
		&overview.PendingReports,
		&overview.ProcessingReports,
		&overview.ResolvedReports,
		&overview.RejectedReports,
		&overview.ClosedReports,
		&overview.ProductReports,
		&overview.UserReports,
		&overview.OrderReports,
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
