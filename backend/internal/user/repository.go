package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cau-used-goods-app/backend/internal/admin"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	ID            uint64  `json:"id"`
	Nickname      *string `json:"nickname"`
	AvatarURL     *string `json:"avatarUrl"`
	StudentID     *string `json:"studentId"`
	RealName      *string `json:"realName"`
	College       *string `json:"college"`
	Phone         *string `json:"phone"`
	Role          string  `json:"role"`
	AuthStatus    string  `json:"authStatus"`
	AccountStatus string  `json:"accountStatus"`
}

type AdminUser struct {
	User
	OpenID        *string `json:"openid,omitempty"`
	LastLoginTime *string `json:"lastLoginTime,omitempty"`
	CreateTime    string  `json:"createTime"`
	UpdateTime    string  `json:"updateTime"`
	IsDeleted     bool    `json:"isDeleted"`
}

type AdminUserQuery struct {
	Keyword       string
	AuthStatus    string
	AccountStatus string
	Role          string
	IncludeOpenID bool
	Page          int
	PageSize      int
}

type PublicProfile struct {
	ID             uint64  `json:"id"`
	Nickname       *string `json:"nickname"`
	AvatarURL      *string `json:"avatarUrl"`
	AuthStatus     string  `json:"authStatus"`
	TradeAvailable bool    `json:"tradeAvailable"`
}

type PublicHomepageStats struct {
	OnSaleProductCount  int      `json:"onSaleProductCount"`
	CompletedOrderCount int      `json:"completedOrderCount"`
	ReviewReceivedCount int      `json:"reviewReceivedCount"`
	AverageRating       *float64 `json:"averageRating,omitempty"`
}

type PublicHomepage struct {
	Profile  *PublicProfile               `json:"profile"`
	Stats    PublicHomepageStats          `json:"stats"`
	Products PagedResult[UserProductItem] `json:"products"`
}

type Restriction struct {
	AccountStatus string  `json:"accountStatus"`
	Reason        *string `json:"reason,omitempty"`
	OperationType *string `json:"operationType,omitempty"`
	RelatedType   *string `json:"relatedType,omitempty"`
	RelatedID     *uint64 `json:"relatedId,omitempty"`
	CreateTime    *string `json:"createTime,omitempty"`
}

type UserStats struct {
	ProductCount         int `json:"productCount"`
	OnSaleProductCount   int `json:"onSaleProductCount"`
	OrderCount           int `json:"orderCount"`
	CompletedOrderCount  int `json:"completedOrderCount"`
	ReportSubmittedCount int `json:"reportSubmittedCount"`
	ReportedCount        int `json:"reportedCount"`
	AppealCount          int `json:"appealCount"`
	ReviewGivenCount     int `json:"reviewGivenCount"`
	ReviewReceivedCount  int `json:"reviewReceivedCount"`
	AdminOperationCount  int `json:"adminOperationCount"`
}

type AdminUserDetail struct {
	User            *AdminUser       `json:"user"`
	Stats           UserStats        `json:"stats"`
	RecentReports   []UserReportItem `json:"recentReports"`
	RecentAppeals   []UserAppealItem `json:"recentAppeals"`
	RecentAdminLogs []UserLogItem    `json:"recentAdminLogs"`
}

type UserProductItem struct {
	ID            uint64  `json:"id"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	Status        string  `json:"status"`
	ViewCount     int     `json:"viewCount"`
	FavoriteCount int     `json:"favoriteCount"`
	CreateTime    string  `json:"createTime"`
	ImageURL      *string `json:"imageUrl,omitempty"`
}

type UserOrderItem struct {
	ID                   uint64  `json:"id"`
	OrderNo              string  `json:"orderNo"`
	ProductID            uint64  `json:"productId"`
	ProductTitleSnapshot string  `json:"productTitleSnapshot"`
	ProductPriceSnapshot float64 `json:"productPriceSnapshot"`
	BuyerID              uint64  `json:"buyerId"`
	BuyerNickname        *string `json:"buyerNickname,omitempty"`
	SellerID             uint64  `json:"sellerId"`
	SellerNickname       *string `json:"sellerNickname,omitempty"`
	Status               string  `json:"status"`
	CreateTime           string  `json:"createTime"`
	UpdateTime           string  `json:"updateTime"`
}

type UserReportItem struct {
	ID               uint64  `json:"id"`
	Relation         string  `json:"relation"`
	ReporterID       uint64  `json:"reporterId"`
	ReporterNickname *string `json:"reporterNickname,omitempty"`
	TargetType       string  `json:"targetType"`
	TargetID         uint64  `json:"targetId"`
	ReasonType       string  `json:"reasonType"`
	Status           string  `json:"status"`
	HandleResult     *string `json:"handleResult,omitempty"`
	CreateTime       string  `json:"createTime"`
	UpdateTime       string  `json:"updateTime"`
}

type UserAppealItem struct {
	ID                uint64  `json:"id"`
	Relation          string  `json:"relation"`
	AppellantID       uint64  `json:"appellantId"`
	AppellantNickname *string `json:"appellantNickname,omitempty"`
	TargetType        string  `json:"targetType"`
	TargetID          uint64  `json:"targetId"`
	Reason            string  `json:"reason"`
	Status            string  `json:"status"`
	HandleResult      *string `json:"handleResult,omitempty"`
	CreateTime        string  `json:"createTime"`
	UpdateTime        string  `json:"updateTime"`
}

type UserReviewItem struct {
	ID               uint64  `json:"id"`
	Relation         string  `json:"relation"`
	OrderID          uint64  `json:"orderId"`
	ProductID        uint64  `json:"productId"`
	ProductTitle     string  `json:"productTitle"`
	ReviewerID       uint64  `json:"reviewerId"`
	ReviewerNickname *string `json:"reviewerNickname,omitempty"`
	SellerID         uint64  `json:"sellerId"`
	SellerNickname   *string `json:"sellerNickname,omitempty"`
	Rating           int     `json:"rating"`
	Content          *string `json:"content,omitempty"`
	Anonymous        bool    `json:"anonymous"`
	Status           string  `json:"status"`
	CreateTime       string  `json:"createTime"`
}

type UserLogItem struct {
	ID            uint64  `json:"id"`
	AdminID       uint64  `json:"adminId"`
	AdminNickname *string `json:"adminNickname,omitempty"`
	OperationType string  `json:"operationType"`
	TargetType    string  `json:"targetType"`
	TargetID      uint64  `json:"targetId"`
	Description   *string `json:"description,omitempty"`
	IPAddress     *string `json:"ipAddress,omitempty"`
	RelatedType   *string `json:"relatedType,omitempty"`
	RelatedID     *uint64 `json:"relatedId,omitempty"`
	CreateTime    string  `json:"createTime"`
}

type PagedResult[T any] struct {
	Items    []T `json:"items"`
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type StudentVerification struct {
	UserID        uint64  `json:"userId,omitempty"`
	StudentID     *string `json:"studentId"`
	RealName      *string `json:"realName"`
	College       *string `json:"college"`
	AuthStatus    string  `json:"authStatus"`
	AccountStatus string  `json:"accountStatus"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, userID uint64) (*User, error) {
	const query = `
SELECT id, nickname, avatar_url, student_id, real_name, college, phone, role, auth_status, account_status
FROM users
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	var user User
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Nickname,
		&user.AvatarURL,
		&user.StudentID,
		&user.RealName,
		&user.College,
		&user.Phone,
		&user.Role,
		&user.AuthStatus,
		&user.AccountStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *Repository) FindAdminUserByID(ctx context.Context, userID uint64) (*AdminUser, error) {
	const query = `
SELECT id, openid, nickname, avatar_url, student_id, real_name, college, phone, role,
       auth_status, account_status,
       DATE_FORMAT(last_login_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'),
       is_deleted
FROM users
WHERE id = ?
LIMIT 1`

	item, err := scanAdminUser(r.db.QueryRowContext(ctx, query, userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find admin user by id: %w", err)
	}
	return item, nil
}

func (r *Repository) ListAdminUsers(ctx context.Context, query AdminUserQuery) ([]AdminUser, int, error) {
	whereSQL, args := buildAdminUserWhere(query)

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users"+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count admin users: %w", err)
	}

	listSQL := `
SELECT id, openid, nickname, avatar_url, student_id, real_name, college, phone, role,
       auth_status, account_status,
       DATE_FORMAT(last_login_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'),
       is_deleted
FROM users` + whereSQL + `
ORDER BY create_time DESC, id DESC
LIMIT ? OFFSET ?`
	args = append(args, query.PageSize, (query.Page-1)*query.PageSize)

	rows, err := r.db.QueryContext(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list admin users: %w", err)
	}
	defer rows.Close()

	items := make([]AdminUser, 0)
	for rows.Next() {
		item, err := scanAdminUser(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin users: %w", err)
	}
	return items, total, nil
}

func (r *Repository) UserStats(ctx context.Context, userID uint64) (UserStats, error) {
	var stats UserStats
	queries := []struct {
		dest *int
		sql  string
		args []interface{}
	}{
		{&stats.ProductCount, "SELECT COUNT(*) FROM products WHERE seller_id = ? AND is_deleted = 0", []interface{}{userID}},
		{&stats.OnSaleProductCount, "SELECT COUNT(*) FROM products WHERE seller_id = ? AND status = 'ON_SALE' AND is_deleted = 0", []interface{}{userID}},
		{&stats.OrderCount, "SELECT COUNT(*) FROM orders WHERE buyer_id = ? OR seller_id = ?", []interface{}{userID, userID}},
		{&stats.CompletedOrderCount, "SELECT COUNT(*) FROM orders WHERE (buyer_id = ? OR seller_id = ?) AND status = 'COMPLETED'", []interface{}{userID, userID}},
		{&stats.ReportSubmittedCount, "SELECT COUNT(*) FROM reports WHERE reporter_id = ?", []interface{}{userID}},
		{&stats.ReportedCount, "SELECT COUNT(*) FROM reports WHERE target_type = 'USER' AND target_id = ?", []interface{}{userID}},
		{&stats.AppealCount, "SELECT COUNT(*) FROM appeals WHERE appellant_id = ? OR (target_type = 'USER' AND target_id = ?)", []interface{}{userID, userID}},
		{&stats.ReviewGivenCount, "SELECT COUNT(*) FROM reviews WHERE reviewer_id = ? AND is_deleted = 0", []interface{}{userID}},
		{&stats.ReviewReceivedCount, "SELECT COUNT(*) FROM reviews WHERE seller_id = ? AND is_deleted = 0", []interface{}{userID}},
		{&stats.AdminOperationCount, "SELECT COUNT(*) FROM admin_logs WHERE target_type = 'USER' AND target_id = ?", []interface{}{userID}},
	}
	for _, item := range queries {
		if err := r.db.QueryRowContext(ctx, item.sql, item.args...).Scan(item.dest); err != nil {
			if isMissingTable(err) {
				continue
			}
			return UserStats{}, fmt.Errorf("query user stats: %w", err)
		}
	}
	return stats, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uint64, nickname, avatarURL, phone *string) error {
	const execSQL = `
UPDATE users
SET nickname = COALESCE(?, nickname),
    avatar_url = COALESCE(?, avatar_url),
    phone = COALESCE(?, phone),
    update_time = NOW()
WHERE id = ? AND account_status = 'NORMAL' AND is_deleted = 0`

	result, err := r.db.ExecContext(ctx, execSQL, nickname, avatarURL, phone, userID)
	if err != nil {
		return fmt.Errorf("update user profile: %w", err)
	}
	return checkAffected(result, "当前账号状态不可操作")
}

func (r *Repository) SubmitStudentVerification(ctx context.Context, userID uint64, studentID, realName, college string) error {
	const execSQL = `
UPDATE users
SET student_id = ?,
    real_name = ?,
    college = ?,
    auth_status = 'PENDING',
    update_time = NOW()
WHERE id = ? AND account_status = 'NORMAL' AND is_deleted = 0`

	result, err := r.db.ExecContext(ctx, execSQL, studentID, realName, college, userID)
	if err != nil {
		if isDuplicateEntry(err) {
			return fmt.Errorf("学号已被使用")
		}
		return fmt.Errorf("submit student verification: %w", err)
	}
	return checkAffected(result, "当前账号状态不可操作")
}

func (r *Repository) StudentIDUsedByOther(ctx context.Context, studentID string, userID uint64) (bool, error) {
	const query = `
SELECT 1
FROM users
WHERE student_id = ? AND id <> ? AND is_deleted = 0
LIMIT 1`

	var exists int
	if err := r.db.QueryRowContext(ctx, query, studentID, userID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check student id uniqueness: %w", err)
	}
	return true, nil
}

func (r *Repository) FindStudentVerification(ctx context.Context, userID uint64) (*StudentVerification, error) {
	const query = `
SELECT id, student_id, real_name, college, auth_status, account_status
FROM users
WHERE id = ? AND is_deleted = 0
LIMIT 1`

	var verification StudentVerification
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&verification.UserID,
		&verification.StudentID,
		&verification.RealName,
		&verification.College,
		&verification.AuthStatus,
		&verification.AccountStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("find student verification: %w", err)
	}
	return &verification, nil
}

func (r *Repository) ListStudentVerifications(ctx context.Context, status string) ([]StudentVerification, error) {
	const query = `
SELECT id, student_id, real_name, college, auth_status, account_status
FROM users
WHERE auth_status = ? AND is_deleted = 0
ORDER BY update_time DESC`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("list student verifications: %w", err)
	}
	defer rows.Close()

	items := make([]StudentVerification, 0)
	for rows.Next() {
		var item StudentVerification
		if err := rows.Scan(
			&item.UserID,
			&item.StudentID,
			&item.RealName,
			&item.College,
			&item.AuthStatus,
			&item.AccountStatus,
		); err != nil {
			return nil, fmt.Errorf("scan student verification: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate student verifications: %w", err)
	}
	return items, nil
}

func (r *Repository) ReviewStudentVerification(ctx context.Context, adminID uint64, userID uint64, authStatus string, description string, adminLogger *admin.Service) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin review student verification tx: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	result, err := tx.ExecContext(ctx, `
UPDATE users
SET auth_status = ?, update_time = NOW()
WHERE id = ? AND auth_status = 'PENDING' AND account_status = 'NORMAL' AND is_deleted = 0`, authStatus, userID)
	if err != nil {
		return fmt.Errorf("review student verification: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get review affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("认证状态已变化，请刷新后重试")
	}

	operationType := admin.OperationStudentVerifyReject
	if authStatus == "VERIFIED" {
		operationType = admin.OperationStudentVerifyApprove
	}

	if adminLogger == nil {
		return fmt.Errorf("admin logger is not configured")
	}
	if _, err := adminLogger.LogActionTx(ctx, tx, admin.LogActionInput{
		AdminID:       adminID,
		OperationType: operationType,
		TargetType:    admin.TargetTypeUser,
		TargetID:      userID,
		Description:   &description,
	}); err != nil {
		return fmt.Errorf("create student verification admin log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit review student verification tx: %w", err)
	}
	committed = true
	return nil
}

func (r *Repository) ListUserProducts(ctx context.Context, userID uint64, page, pageSize int, dataScope string) ([]UserProductItem, int, error) {
	whereSQL := " WHERE p.seller_id = ?" + dataScopeSQL("p", dataScope)
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products p"+whereSQL, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user products: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT p.id, p.title, p.price, p.status, p.view_count, p.favorite_count,
       DATE_FORMAT(p.create_time, '%Y-%m-%d %H:%i:%s'),
       (SELECT image_url FROM product_images pi WHERE pi.product_id = p.id ORDER BY pi.sort_order, pi.id LIMIT 1)
FROM products p`+whereSQL+`
ORDER BY p.create_time DESC, p.id DESC
LIMIT ? OFFSET ?`, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list user products: %w", err)
	}
	defer rows.Close()

	items := make([]UserProductItem, 0)
	for rows.Next() {
		var item UserProductItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Price, &item.Status, &item.ViewCount, &item.FavoriteCount, &item.CreateTime, &item.ImageURL); err != nil {
			return nil, 0, fmt.Errorf("scan user product: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListUserOrders(ctx context.Context, userID uint64, page, pageSize int) ([]UserOrderItem, int, error) {
	const whereSQL = " WHERE o.buyer_id = ? OR o.seller_id = ?"
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM orders o"+whereSQL, userID, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user orders: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT o.id, o.order_no, o.product_id, o.product_title_snapshot, o.product_price_snapshot,
       o.buyer_id, buyer.nickname, o.seller_id, seller.nickname, o.status,
       DATE_FORMAT(o.create_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(o.update_time, '%Y-%m-%d %H:%i:%s')
FROM orders o
LEFT JOIN users buyer ON buyer.id = o.buyer_id
LEFT JOIN users seller ON seller.id = o.seller_id`+whereSQL+`
ORDER BY o.create_time DESC, o.id DESC
LIMIT ? OFFSET ?`, userID, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list user orders: %w", err)
	}
	defer rows.Close()

	items := make([]UserOrderItem, 0)
	for rows.Next() {
		var item UserOrderItem
		if err := rows.Scan(
			&item.ID, &item.OrderNo, &item.ProductID, &item.ProductTitleSnapshot, &item.ProductPriceSnapshot,
			&item.BuyerID, &item.BuyerNickname, &item.SellerID, &item.SellerNickname, &item.Status,
			&item.CreateTime, &item.UpdateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user order: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListUserReports(ctx context.Context, userID uint64, page, pageSize int) ([]UserReportItem, int, error) {
	const whereSQL = " WHERE r.reporter_id = ? OR (r.target_type = 'USER' AND r.target_id = ?)"
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reports r"+whereSQL, userID, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user reports: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT r.id,
       CASE WHEN r.reporter_id = ? THEN 'SUBMITTED' ELSE 'REPORTED' END,
       r.reporter_id, u.nickname, r.target_type, r.target_id, r.reason_type, r.status,
       r.handle_result,
       DATE_FORMAT(r.create_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(r.update_time, '%Y-%m-%d %H:%i:%s')
FROM reports r
LEFT JOIN users u ON u.id = r.reporter_id`+whereSQL+`
ORDER BY r.create_time DESC, r.id DESC
LIMIT ? OFFSET ?`, userID, userID, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list user reports: %w", err)
	}
	defer rows.Close()

	items := make([]UserReportItem, 0)
	for rows.Next() {
		var item UserReportItem
		if err := rows.Scan(
			&item.ID, &item.Relation, &item.ReporterID, &item.ReporterNickname, &item.TargetType,
			&item.TargetID, &item.ReasonType, &item.Status, &item.HandleResult, &item.CreateTime, &item.UpdateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user report: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListUserAppeals(ctx context.Context, userID uint64, page, pageSize int) ([]UserAppealItem, int, error) {
	const whereSQL = " WHERE a.appellant_id = ? OR (a.target_type = 'USER' AND a.target_id = ?)"
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM appeals a"+whereSQL, userID, userID).Scan(&total); err != nil {
		if isMissingTable(err) {
			return []UserAppealItem{}, 0, nil
		}
		return nil, 0, fmt.Errorf("count user appeals: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT a.id,
       CASE WHEN a.appellant_id = ? THEN 'SUBMITTED' ELSE 'TARGET' END,
       a.appellant_id, u.nickname, a.target_type, a.target_id, a.reason, a.status,
       a.handle_result,
       DATE_FORMAT(a.create_time, '%Y-%m-%d %H:%i:%s'),
       DATE_FORMAT(a.update_time, '%Y-%m-%d %H:%i:%s')
FROM appeals a
LEFT JOIN users u ON u.id = a.appellant_id`+whereSQL+`
ORDER BY a.create_time DESC, a.id DESC
LIMIT ? OFFSET ?`, userID, userID, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		if isMissingTable(err) {
			return []UserAppealItem{}, 0, nil
		}
		return nil, 0, fmt.Errorf("list user appeals: %w", err)
	}
	defer rows.Close()

	items := make([]UserAppealItem, 0)
	for rows.Next() {
		var item UserAppealItem
		if err := rows.Scan(
			&item.ID, &item.Relation, &item.AppellantID, &item.AppellantNickname, &item.TargetType,
			&item.TargetID, &item.Reason, &item.Status, &item.HandleResult, &item.CreateTime, &item.UpdateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user appeal: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListUserReviews(ctx context.Context, userID uint64, page, pageSize int, dataScope string) ([]UserReviewItem, int, error) {
	whereSQL := " WHERE (rv.reviewer_id = ? OR rv.seller_id = ?)" + dataScopeSQL("rv", dataScope)
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews rv"+whereSQL, userID, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user reviews: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT rv.id,
       CASE WHEN rv.reviewer_id = ? THEN 'GIVEN' ELSE 'RECEIVED' END,
       rv.order_id, rv.product_id, p.title,
       rv.reviewer_id, reviewer.nickname, rv.seller_id, seller.nickname,
       rv.rating, rv.content, rv.is_anonymous, rv.status,
       DATE_FORMAT(rv.create_time, '%Y-%m-%d %H:%i:%s')
FROM reviews rv
LEFT JOIN products p ON p.id = rv.product_id
LEFT JOIN users reviewer ON reviewer.id = rv.reviewer_id
LEFT JOIN users seller ON seller.id = rv.seller_id`+whereSQL+`
ORDER BY rv.create_time DESC, rv.id DESC
LIMIT ? OFFSET ?`, userID, userID, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list user reviews: %w", err)
	}
	defer rows.Close()

	items := make([]UserReviewItem, 0)
	for rows.Next() {
		var item UserReviewItem
		if err := rows.Scan(
			&item.ID, &item.Relation, &item.OrderID, &item.ProductID, &item.ProductTitle,
			&item.ReviewerID, &item.ReviewerNickname, &item.SellerID, &item.SellerNickname,
			&item.Rating, &item.Content, &item.Anonymous, &item.Status, &item.CreateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user review: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (r *Repository) ListUserLogs(ctx context.Context, userID uint64, page, pageSize int) ([]UserLogItem, int, error) {
	const whereSQL = " WHERE l.target_type = 'USER' AND l.target_id = ?"
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM admin_logs l"+whereSQL, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count user logs: %w", err)
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT l.id, l.admin_id, u.nickname, l.operation_type, l.target_type, l.target_id,
       l.description, l.ip_address, l.related_type, l.related_id,
       DATE_FORMAT(l.create_time, '%Y-%m-%d %H:%i:%s')
FROM admin_logs l
LEFT JOIN users u ON u.id = l.admin_id`+whereSQL+`
ORDER BY l.create_time DESC, l.id DESC
LIMIT ? OFFSET ?`, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list user logs: %w", err)
	}
	defer rows.Close()

	items := make([]UserLogItem, 0)
	for rows.Next() {
		var item UserLogItem
		if err := rows.Scan(
			&item.ID, &item.AdminID, &item.AdminNickname, &item.OperationType, &item.TargetType,
			&item.TargetID, &item.Description, &item.IPAddress, &item.RelatedType, &item.RelatedID, &item.CreateTime,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user log: %w", err)
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func checkAffected(result sql.Result, notFoundMsg string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf(notFoundMsg)
	}
	return nil
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func isMissingTable(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1146
}

func buildAdminUserWhere(query AdminUserQuery) (string, []interface{}) {
	whereSQL := " WHERE 1 = 1"
	args := make([]interface{}, 0)

	if query.AccountStatus == "" {
		whereSQL += " AND account_status <> 'CANCELED'"
	}
	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		if query.IncludeOpenID {
			whereSQL += " AND (openid LIKE ? OR nickname LIKE ? OR student_id LIKE ? OR real_name LIKE ? OR phone LIKE ?)"
			args = append(args, keyword, keyword, keyword, keyword, keyword)
		} else {
			whereSQL += " AND (nickname LIKE ? OR student_id LIKE ? OR real_name LIKE ? OR phone LIKE ?)"
			args = append(args, keyword, keyword, keyword, keyword)
		}
	}
	if query.AuthStatus != "" {
		whereSQL += " AND auth_status = ?"
		args = append(args, query.AuthStatus)
	}
	if query.AccountStatus != "" {
		whereSQL += " AND account_status = ?"
		args = append(args, query.AccountStatus)
	}
	if query.Role != "" {
		whereSQL += " AND role = ?"
		args = append(args, query.Role)
	}
	return whereSQL, args
}

func scanAdminUser(scanner interface {
	Scan(dest ...interface{}) error
}) (*AdminUser, error) {
	var item AdminUser
	var openID, nickname, avatarURL, studentID, realName, college, phone, lastLoginTime sql.NullString
	if err := scanner.Scan(
		&item.ID,
		&openID,
		&nickname,
		&avatarURL,
		&studentID,
		&realName,
		&college,
		&phone,
		&item.Role,
		&item.AuthStatus,
		&item.AccountStatus,
		&lastLoginTime,
		&item.CreateTime,
		&item.UpdateTime,
		&item.IsDeleted,
	); err != nil {
		return nil, err
	}
	item.OpenID = nullStringPtr(openID)
	item.Nickname = nullStringPtr(nickname)
	item.AvatarURL = nullStringPtr(avatarURL)
	item.StudentID = nullStringPtr(studentID)
	item.RealName = nullStringPtr(realName)
	item.College = nullStringPtr(college)
	item.Phone = nullStringPtr(phone)
	item.LastLoginTime = nullStringPtr(lastLoginTime)
	return &item, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func nullableString(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}

func dataScopeSQL(alias, dataScope string) string {
	switch dataScope {
	case "DELETED":
		return " AND " + alias + ".is_deleted = 1"
	case "ALL":
		return ""
	default:
		return " AND " + alias + ".is_deleted = 0"
	}
}
