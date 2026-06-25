package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Category struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	ParentID  uint64 `json:"parentId"`
	SortOrder int    `json:"sortOrder"`
	Status    string `json:"status"`
}

type Product struct {
	ID             uint64   `json:"id"`
	SellerID       uint64   `json:"sellerId"`
	BuyerID        *uint64  `json:"buyerId,omitempty"`
	CategoryID     uint64   `json:"categoryId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	OriginalPrice  *float64 `json:"originalPrice"`
	Price          float64  `json:"price"`
	ConditionLevel string   `json:"conditionLevel"`
	MeetLocation   string   `json:"meetLocation"`
	Status         string   `json:"status"`
	OffShelfBy     *string  `json:"offShelfBy,omitempty"`
	OffShelfReason *string  `json:"offShelfReason,omitempty"`
	ViewCount      int      `json:"viewCount"`
	FavoriteCount  int      `json:"favoriteCount"`
	CreateTime     string   `json:"createTime"`
	Images         []string `json:"images"`
}

const (
	OffShelfByUser          = "USER"
	OffShelfByAdmin         = "ADMIN"
	OffShelfBySystem        = "SYSTEM"
	OffShelfByAccountStatus = "ACCOUNT_STATUS"
)

type ProductNoticeInfo struct {
	ID       uint64
	SellerID uint64
	Title    string
	Status   string
}

func (r *Repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, parent_id, sort_order, status
		FROM categories
		WHERE status = 'ENABLED'
		ORDER BY sort_order ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.SortOrder, &c.Status); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *Repository) ListAllCategories(ctx context.Context, status string) ([]Category, error) {
	query := `
		SELECT id, name, parent_id, sort_order, status
		FROM categories
		WHERE 1 = 1
	`
	args := []interface{}{}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY sort_order ASC, id ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.ParentID, &c.SortOrder, &c.Status); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

type CreateCategoryInput struct {
	Name      string
	ParentID  uint64
	SortOrder int
	Status    string
}

func (r *Repository) CreateCategory(ctx context.Context, input CreateCategoryInput) (uint64, error) {
	return r.CreateCategoryTx(ctx, nil, input)
}

func (r *Repository) CreateCategoryTx(ctx context.Context, tx *sql.Tx, input CreateCategoryInput) (uint64, error) {
	execer := productExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, `
		INSERT INTO categories (name, parent_id, sort_order, status)
		VALUES (?, ?, ?, ?)
	`, input.Name, input.ParentID, input.SortOrder, input.Status)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

type UpdateCategoryInput struct {
	ID        uint64
	Name      string
	ParentID  uint64
	SortOrder int
	Status    string
}

func (r *Repository) UpdateCategory(ctx context.Context, input UpdateCategoryInput) error {
	return r.UpdateCategoryTx(ctx, nil, input)
}

func (r *Repository) UpdateCategoryTx(ctx context.Context, tx *sql.Tx, input UpdateCategoryInput) error {
	execer := productExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, `
		UPDATE categories
		SET name = ?, parent_id = ?, sort_order = ?, status = ?, update_time = CURRENT_TIMESTAMP
		WHERE id = ?
	`, input.Name, input.ParentID, input.SortOrder, input.Status, input.ID)
	if err != nil {
		return err
	}
	return checkAffected(result)
}

func (r *Repository) UpdateCategoryStatus(ctx context.Context, id uint64, status string) error {
	return r.UpdateCategoryStatusTx(ctx, nil, id, status)
}

func (r *Repository) UpdateCategoryStatusTx(ctx context.Context, tx *sql.Tx, id uint64, status string) error {
	execer := productExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, `
		UPDATE categories
		SET status = ?, update_time = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, id)
	if err != nil {
		return err
	}
	return checkAffected(result)
}

type CreateProductInput struct {
	SellerID       uint64
	CategoryID     uint64
	Title          string
	Description    string
	OriginalPrice  *float64
	Price          float64
	ConditionLevel string
	MeetLocation   string
}

func (r *Repository) CreateProduct(ctx context.Context, input CreateProductInput) (uint64, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO products (
			seller_id, category_id, title, description, original_price,
			price, condition_level, meet_location, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'ON_SALE')
	`,
		input.SellerID,
		input.CategoryID,
		input.Title,
		input.Description,
		input.OriginalPrice,
		input.Price,
		input.ConditionLevel,
		input.MeetLocation,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

type ListProductsInput struct {
	Keyword        string
	CategoryID     uint64
	ConditionLevel string
	Status         string
	MinPrice       *float64
	MaxPrice       *float64
	Sort           string
	Page           int
	PageSize       int
	IncludeDeleted bool
}

type ProductListResult struct {
	List     []Product `json:"list"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total"`
}

func (r *Repository) ListProducts(ctx context.Context, input ListProductsInput) (*ProductListResult, error) {
	input.Status = "ON_SALE"
	input.IncludeDeleted = false
	return r.listProducts(ctx, input)
}

func (r *Repository) ListAdminProducts(ctx context.Context, input ListProductsInput) (*ProductListResult, error) {
	input.IncludeDeleted = true
	return r.listProducts(ctx, input)
}

func (r *Repository) listProducts(ctx context.Context, input ListProductsInput) (*ProductListResult, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}

	where := ` WHERE 1 = 1 `
	args := []any{}
	if !input.IncludeDeleted {
		where += " AND p.is_deleted = 0 "
	}
	if input.Status != "" {
		where += " AND p.status = ? "
		args = append(args, input.Status)
	}

	for _, term := range splitProductSearchKeyword(input.Keyword) {
		where += " AND (p.title LIKE ? OR p.description LIKE ? OR c.name LIKE ?) "
		keyword := "%" + term + "%"
		args = append(args, keyword, keyword, keyword)
	}

	if input.CategoryID > 0 {
		where += " AND p.category_id = ? "
		args = append(args, input.CategoryID)
	}

	if input.ConditionLevel != "" {
		where += " AND p.condition_level = ? "
		args = append(args, input.ConditionLevel)
	}

	if input.MinPrice != nil {
		where += " AND p.price >= ? "
		args = append(args, *input.MinPrice)
	}

	if input.MaxPrice != nil {
		where += " AND p.price <= ? "
		args = append(args, *input.MaxPrice)
	}

	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
	` + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	orderBy := " ORDER BY p.create_time DESC "
	switch input.Sort {
	case "price_asc":
		orderBy = " ORDER BY p.price ASC, p.create_time DESC "
	case "price_desc":
		orderBy = " ORDER BY p.price DESC, p.create_time DESC "
	case "popular":
		orderBy = " ORDER BY p.view_count DESC, p.favorite_count DESC, p.create_time DESC "
	}

	offset := (input.Page - 1) * input.PageSize

	query := `
		SELECT p.id, p.seller_id, p.category_id, p.title, p.description, p.original_price,
		       p.price, p.condition_level, p.meet_location, p.status, p.off_shelf_by, p.view_count,
		       p.favorite_count, DATE_FORMAT(p.create_time, '%Y-%m-%d %H:%i:%s')
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
	` + where + orderBy + ` LIMIT ? OFFSET ?`

	queryArgs := append(args, input.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Product, 0)
	for rows.Next() {
		var p Product
		var desc sql.NullString
		var originalPrice sql.NullFloat64
		var condition sql.NullString
		var location sql.NullString
		var offShelfBy sql.NullString

		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
			&p.Price, &condition, &location, &p.Status, &offShelfBy, &p.ViewCount,
			&p.FavoriteCount, &p.CreateTime,
		); err != nil {
			return nil, err
		}

		fillProductNullableFields(&p, desc, originalPrice, condition, location)
		fillProductOffShelfBy(&p, offShelfBy)

		images, _ := r.ListProductImages(ctx, p.ID)
		p.Images = images

		list = append(list, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &ProductListResult{
		List:     list,
		Page:     input.Page,
		PageSize: input.PageSize,
		Total:    total,
	}, nil
}

func splitProductSearchKeyword(keyword string) []string {
	terms := strings.FieldsFunc(strings.TrimSpace(keyword), func(r rune) bool {
		return unicode.IsSpace(r) || r == ',' || r == '\uFF0C' || r == ';' || r == '\uFF1B'
	})
	if len(terms) == 0 {
		return nil
	}

	result := make([]string, 0, len(terms))
	seen := make(map[string]bool, len(terms))
	for _, term := range terms {
		if seen[term] {
			continue
		}
		seen[term] = true
		result = append(result, term)
	}
	return result
}
func (r *Repository) IncrementViewCount(ctx context.Context, productID uint64, viewer ProductViewer) error {
	if viewer.Role == "ADMIN" || viewer.Role == "SUPER_ADMIN" {
		return nil
	}

	if viewer.UserID == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin product view count tx: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE products
		SET view_count = view_count + 1,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND is_deleted = 0
		  AND status = 'ON_SALE'
		  AND seller_id <> ?
		  AND NOT EXISTS (
		      SELECT 1
		      FROM browse_history
		      WHERE user_id = ?
		        AND product_id = ?
		        AND create_time >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
		  )
	`, productID, viewer.UserID, viewer.UserID, productID)
	if err != nil {
		return fmt.Errorf("increment product view count: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check increment product view count result: %w", err)
	}
	if affected > 0 {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO browse_history (user_id, product_id)
			VALUES (?, ?)
		`, viewer.UserID, productID); err != nil {
			return fmt.Errorf("create browse history: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit product view count tx: %w", err)
	}
	return nil
}
func (r *Repository) GetProductByID(ctx context.Context, id uint64, viewer ProductViewer) (*Product, error) {
	var p Product
	var desc sql.NullString
	var originalPrice sql.NullFloat64
	var condition sql.NullString
	var location sql.NullString
	var offShelfBy sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, seller_id, category_id, title, description, original_price,
		       price, condition_level, meet_location, status, off_shelf_by, view_count,
		       favorite_count, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s')
		FROM products
		WHERE id = ?
		  AND is_deleted = 0
		  AND (
		      status = 'ON_SALE'
		      OR seller_id = ?
		      OR ? IN ('ADMIN', 'SUPER_ADMIN')
		      OR EXISTS (
		          SELECT 1
		          FROM orders o
		          WHERE o.product_id = products.id
		            AND (o.buyer_id = ? OR o.seller_id = ?)
		            AND o.status IN ('PENDING_CONFIRM', 'WAIT_MEET', 'COMPLETED')
		      )
		  )
	`, id, viewer.UserID, viewer.Role, viewer.UserID, viewer.UserID).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
		&p.Price, &condition, &location, &p.Status, &offShelfBy, &p.ViewCount,
		&p.FavoriteCount, &p.CreateTime,
	)
	if err != nil {
		return nil, err
	}

	fillProductNullableFields(&p, desc, originalPrice, condition, location)
	fillProductOffShelfBy(&p, offShelfBy)

	images, err := r.ListProductImages(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Images = images

	return &p, nil
}

func (r *Repository) GetProductNoticeInfo(ctx context.Context, productID uint64) (*ProductNoticeInfo, error) {
	var item ProductNoticeInfo
	err := r.db.QueryRowContext(ctx, `
		SELECT id, seller_id, title, status
		FROM products
		WHERE id = ? AND is_deleted = 0
	`, productID).Scan(&item.ID, &item.SellerID, &item.Title, &item.Status)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) AdminGetProductByID(ctx context.Context, id uint64) (*Product, error) {
	var p Product
	var desc sql.NullString
	var originalPrice sql.NullFloat64
	var condition sql.NullString
	var location sql.NullString
	var offShelfBy sql.NullString
	var buyerID sql.NullInt64

	err := r.db.QueryRowContext(ctx, `
		SELECT p.id, p.seller_id,
		       (
		           SELECT o.buyer_id
		           FROM orders o
		           WHERE o.product_id = p.id
		             AND o.status IN ('PENDING_CONFIRM', 'WAIT_MEET', 'COMPLETED')
		           ORDER BY CASE o.status
		               WHEN 'WAIT_MEET' THEN 1
		               WHEN 'PENDING_CONFIRM' THEN 2
		               WHEN 'COMPLETED' THEN 3
		               ELSE 4
		           END, o.update_time DESC
		           LIMIT 1
		       ) AS buyer_id,
		       p.category_id, p.title, p.description, p.original_price,
		       p.price, p.condition_level, p.meet_location, p.status, p.off_shelf_by, p.view_count,
		       p.favorite_count, DATE_FORMAT(p.create_time, '%Y-%m-%d %H:%i:%s')
		FROM products p
		WHERE p.id = ?
	`, id).Scan(
		&p.ID, &p.SellerID, &buyerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
		&p.Price, &condition, &location, &p.Status, &offShelfBy, &p.ViewCount,
		&p.FavoriteCount, &p.CreateTime,
	)
	if err != nil {
		return nil, err
	}

	fillProductNullableFields(&p, desc, originalPrice, condition, location)
	fillProductOffShelfBy(&p, offShelfBy)
	fillProductBuyerID(&p, buyerID)

	images, err := r.ListProductImages(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Images = images

	return &p, nil
}

func (r *Repository) ListProductImages(ctx context.Context, productID uint64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT image_url
		FROM product_images
		WHERE product_id = ?
		ORDER BY sort_order ASC, id ASC
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]string, 0)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		images = append(images, url)
	}

	return images, rows.Err()
}

func (r *Repository) ListMyProducts(ctx context.Context, sellerID uint64) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, seller_id, category_id, title, description, original_price,
		       price, condition_level, meet_location, status, off_shelf_by, off_shelf_reason, view_count,
		       favorite_count, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s')
		FROM products
		WHERE seller_id = ? AND is_deleted = 0
		ORDER BY create_time DESC
	`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]Product, 0)
	for rows.Next() {
		var p Product
		var desc sql.NullString
		var originalPrice sql.NullFloat64
		var condition sql.NullString
		var location sql.NullString
		var offShelfBy sql.NullString
		var offShelfReason sql.NullString

		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
			&p.Price, &condition, &location, &p.Status, &offShelfBy, &offShelfReason, &p.ViewCount,
			&p.FavoriteCount, &p.CreateTime,
		); err != nil {
			return nil, err
		}

		fillProductNullableFields(&p, desc, originalPrice, condition, location)
		fillProductOffShelfBy(&p, offShelfBy)
		fillProductOffShelfReason(&p, offShelfReason)

		images, _ := r.ListProductImages(ctx, p.ID)
		p.Images = images

		list = append(list, p)
	}

	return list, rows.Err()
}

func (r *Repository) DeleteProduct(ctx context.Context, productID uint64, sellerID uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'DELETED',
		    is_deleted = 1,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND seller_id = ?
		  AND is_deleted = 0
		  AND status IN ('ON_SALE', 'OFF_SHELF')
	`, productID, sellerID)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

type UpdateProductInput struct {
	ProductID      uint64
	SellerID       uint64
	CategoryID     uint64
	Title          string
	Description    string
	OriginalPrice  *float64
	Price          float64
	ConditionLevel string
	MeetLocation   string
}

func (r *Repository) UpdateProduct(ctx context.Context, input UpdateProductInput) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET category_id = ?,
		    title = ?,
		    description = ?,
		    original_price = ?,
		    price = ?,
		    condition_level = ?,
		    meet_location = ?,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND seller_id = ?
		  AND is_deleted = 0
		  AND status IN ('ON_SALE', 'OFF_SHELF')
	`,
		input.CategoryID,
		input.Title,
		input.Description,
		input.OriginalPrice,
		input.Price,
		input.ConditionLevel,
		input.MeetLocation,
		input.ProductID,
		input.SellerID,
	)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

func (r *Repository) UpdateProductStatus(ctx context.Context, productID uint64, sellerID uint64, status string, reason string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = ?,
		    off_shelf_reason = CASE WHEN ? = 'OFF_SHELF' THEN ? ELSE NULL END,
		    off_shelf_by = CASE WHEN ? = 'OFF_SHELF' THEN ? ELSE NULL END,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND seller_id = ?
		  AND is_deleted = 0
		  AND (
		      (? = 'OFF_SHELF' AND status = 'ON_SALE')
		      OR
		      (? = 'ON_SALE' AND status = 'OFF_SHELF' AND UPPER(TRIM(off_shelf_by)) IN (?, ?))
		  )
	`,
		status, status, reason, status, OffShelfByUser, productID, sellerID, status, status,
		OffShelfByUser, OffShelfByAccountStatus,
	)
	if err != nil {
		return err
	}

	if err := checkAffected(result); err != nil {
		if errors.Is(err, sql.ErrNoRows) && status == "ON_SALE" {
			var offShelfBy sql.NullString
			checkErr := r.db.QueryRowContext(ctx, `
				SELECT off_shelf_by
				FROM products
				WHERE id = ?
				  AND seller_id = ?
				  AND is_deleted = 0
				  AND status = 'OFF_SHELF'
				LIMIT 1
			`, productID, sellerID).Scan(&offShelfBy)
			if checkErr == nil && offShelfBy.Valid && offShelfBy.String == OffShelfByAdmin {
				return fmt.Errorf("管理员下架的商品不能自行上架")
			}
			if checkErr == nil && offShelfBy.Valid && offShelfBy.String == OffShelfBySystem {
				return fmt.Errorf("系统下架的商品不能自行上架")
			}
		}
		return err
	}
	return nil
}

func (r *Repository) AdminUpdateProductStatus(ctx context.Context, input AdminUpdateProductStatusInput) error {
	return r.AdminUpdateProductStatusTx(ctx, nil, input)
}

func (r *Repository) AdminUpdateProductStatusTx(ctx context.Context, tx *sql.Tx, input AdminUpdateProductStatusInput) error {
	isDeleted := 0
	if input.Status == "DELETED" {
		isDeleted = 1
	}
	execer := productExecutor(r.db)
	if tx != nil {
		execer = tx
	}
	result, err := execer.ExecContext(ctx, `
		UPDATE products
		SET status = ?,
		    off_shelf_reason = CASE WHEN ? = 'OFF_SHELF' THEN ? ELSE NULL END,
		    off_shelf_by = CASE WHEN ? = 'OFF_SHELF' THEN ? ELSE NULL END,
		    is_deleted = ?,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
	`, input.Status, input.Status, input.Reason, input.Status, OffShelfByAdmin, isDeleted, input.ProductID)
	if err != nil {
		return fmt.Errorf("update product status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check update product status result: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

type productExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (r *Repository) OffShelfOnSaleBySellerTx(ctx context.Context, tx *sql.Tx, sellerID uint64, reason string) ([]uint64, error) {
	return r.OffShelfOnSaleBySellerWithSourceTx(ctx, tx, sellerID, reason, OffShelfBySystem)
}

func (r *Repository) OffShelfOnSaleBySellerWithSourceTx(ctx context.Context, tx *sql.Tx, sellerID uint64, reason string, source string) ([]uint64, error) {
	if source == "" {
		source = OffShelfBySystem
	}
	rows, err := tx.QueryContext(ctx, `
		SELECT id
		FROM products
		WHERE seller_id = ?
		  AND is_deleted = 0
		  AND status = 'ON_SALE'
		ORDER BY id ASC
		FOR UPDATE
	`, sellerID)
	if err != nil {
		return nil, fmt.Errorf("list seller on-sale products: %w", err)
	}
	defer rows.Close()

	var productIDs []uint64
	for rows.Next() {
		var productID uint64
		if err := rows.Scan(&productID); err != nil {
			return nil, fmt.Errorf("scan seller on-sale product: %w", err)
		}
		productIDs = append(productIDs, productID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate seller on-sale products: %w", err)
	}
	if len(productIDs) == 0 {
		return productIDs, nil
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE products
		SET status = 'OFF_SHELF',
		    off_shelf_reason = ?,
		    off_shelf_by = ?,
		    update_time = CURRENT_TIMESTAMP
		WHERE seller_id = ?
		  AND is_deleted = 0
		  AND status = 'ON_SALE'
	`, reason, source, sellerID)
	if err != nil {
		return nil, fmt.Errorf("off shelf seller products: %w", err)
	}
	return productIDs, nil
}

func (r *Repository) BatchPutOnSaleRestorable(ctx context.Context, sellerID uint64) (int64, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'ON_SALE',
		    off_shelf_reason = NULL,
		    off_shelf_by = NULL,
		    update_time = CURRENT_TIMESTAMP
		WHERE seller_id = ?
		  AND is_deleted = 0
		  AND status = 'OFF_SHELF'
		  AND UPPER(TRIM(off_shelf_by)) IN (?, ?)
	`, sellerID, OffShelfByUser, OffShelfByAccountStatus)
	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (r *Repository) ValidateRelatedRecordTx(ctx context.Context, tx *sql.Tx, relatedType string, relatedID, productID uint64) error {
	if relatedType == "" {
		return nil
	}
	if tx == nil {
		return fmt.Errorf("transaction is required")
	}

	var query string
	var args []interface{}
	switch relatedType {
	case "REPORT":
		query = "SELECT COUNT(*) FROM reports WHERE id = ? AND target_type = 'PRODUCT' AND target_id = ?"
		args = []interface{}{relatedID, productID}
	case "APPEAL":
		query = "SELECT COUNT(*) FROM appeals WHERE id = ? AND target_type = 'PRODUCT' AND target_id = ?"
		args = []interface{}{relatedID, productID}
	default:
		return fmt.Errorf("relatedType must be REPORT or APPEAL")
	}

	var count int
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return fmt.Errorf("check related record: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("related record not found")
	}
	return nil
}

type AddProductImagesInput struct {
	ProductID uint64
	SellerID  uint64
	Images    []string
}

func (r *Repository) AddProductImages(ctx context.Context, input AddProductImagesInput) error {
	if len(input.Images) == 0 {
		return sql.ErrNoRows
	}
	if len(input.Images) > 9 {
		return fmt.Errorf("最多只能上传9张图片")
	}

	var sellerID uint64
	var status string
	err := r.db.QueryRowContext(ctx, `
		SELECT seller_id, status
		FROM products
		WHERE id = ? AND is_deleted = 0
	`, input.ProductID).Scan(&sellerID, &status)
	if err != nil {
		return err
	}

	if sellerID != input.SellerID {
		return fmt.Errorf("无权限操作该商品")
	}
	if status == "SOLD" || status == "DELETED" {
		return fmt.Errorf("当前商品状态不允许添加图片")
	}

	var currentCount int
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM product_images
		WHERE product_id = ?
	`, input.ProductID).Scan(&currentCount)
	if err != nil {
		return err
	}

	if currentCount+len(input.Images) > 9 {
		return fmt.Errorf("商品图片最多9张")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, imageURL := range input.Images {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO product_images (product_id, image_url, sort_order)
			VALUES (?, ?, ?)
		`, input.ProductID, imageURL, currentCount+i)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

type DeleteProductImageInput struct {
	ProductID uint64
	SellerID  uint64
	ImageID   uint64
}

func (r *Repository) DeleteProductImage(ctx context.Context, input DeleteProductImageInput) error {
	var sellerID uint64
	var status string
	err := r.db.QueryRowContext(ctx, `
		SELECT seller_id, status
		FROM products
		WHERE id = ? AND is_deleted = 0
	`, input.ProductID).Scan(&sellerID, &status)
	if err != nil {
		return err
	}

	if sellerID != input.SellerID {
		return fmt.Errorf("无权限操作该商品")
	}
	if status == "SOLD" || status == "DELETED" {
		return fmt.Errorf("当前商品状态不允许删除图片")
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM product_images
		WHERE id = ? AND product_id = ?
	`, input.ImageID, input.ProductID)
	if err != nil {
		return err
	}

	return checkAffected(result)
}

type ReplaceProductImagesInput struct {
	ProductID uint64
	SellerID  uint64
	Images    []string
}

func (r *Repository) ReplaceProductImages(ctx context.Context, input ReplaceProductImagesInput) error {
	if len(input.Images) > 9 {
		return fmt.Errorf("商品图片最多9张")
	}

	var sellerID uint64
	var status string
	err := r.db.QueryRowContext(ctx, `
		SELECT seller_id, status
		FROM products
		WHERE id = ? AND is_deleted = 0
	`, input.ProductID).Scan(&sellerID, &status)
	if err != nil {
		return err
	}

	if sellerID != input.SellerID {
		return fmt.Errorf("无权限操作该商品")
	}
	if status == "SOLD" || status == "DELETED" {
		return fmt.Errorf("当前商品状态不允许替换图片")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM product_images
		WHERE product_id = ?
	`, input.ProductID); err != nil {
		return err
	}

	for i, imageURL := range input.Images {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO product_images (product_id, image_url, sort_order)
			VALUES (?, ?, ?)
		`, input.ProductID, imageURL, i); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) LockProduct(ctx context.Context, productID uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'LOCKED',
		    off_shelf_reason = NULL,
		    off_shelf_by = NULL,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND is_deleted = 0
		  AND status = 'ON_SALE'
	`, productID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product not available")
	}

	return nil
}

func (r *Repository) UnlockProduct(ctx context.Context, productID uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'ON_SALE',
		    off_shelf_reason = NULL,
		    off_shelf_by = NULL,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND is_deleted = 0
		  AND status = 'LOCKED'
	`, productID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product not locked")
	}

	return nil
}

func (r *Repository) MarkProductSold(ctx context.Context, productID uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'SOLD',
		    off_shelf_reason = NULL,
		    off_shelf_by = NULL,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND is_deleted = 0
		  AND status = 'LOCKED'
	`, productID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("product not locked")
	}

	return nil
}
func fillProductNullableFields(p *Product, desc sql.NullString, originalPrice sql.NullFloat64, condition sql.NullString, location sql.NullString) {
	if desc.Valid {
		p.Description = desc.String
	}
	if originalPrice.Valid {
		p.OriginalPrice = &originalPrice.Float64
	}
	if condition.Valid {
		p.ConditionLevel = condition.String
	}
	if location.Valid {
		p.MeetLocation = location.String
	}
}

func fillProductOffShelfBy(p *Product, offShelfBy sql.NullString) {
	if offShelfBy.Valid {
		p.OffShelfBy = &offShelfBy.String
	}
}

func fillProductOffShelfReason(p *Product, offShelfReason sql.NullString) {
	if offShelfReason.Valid {
		p.OffShelfReason = &offShelfReason.String
	}
}

func fillProductBuyerID(p *Product, buyerID sql.NullInt64) {
	if buyerID.Valid {
		value := uint64(buyerID.Int64)
		p.BuyerID = &value
	}
}

func checkAffected(result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
