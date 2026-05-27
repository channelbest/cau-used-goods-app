package product

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
	CategoryID     uint64   `json:"categoryId"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	OriginalPrice  *float64 `json:"originalPrice"`
	Price          float64  `json:"price"`
	ConditionLevel string   `json:"conditionLevel"`
	MeetLocation   string   `json:"meetLocation"`
	Status         string   `json:"status"`
	ViewCount      int      `json:"viewCount"`
	FavoriteCount  int      `json:"favoriteCount"`
	CreateTime     string   `json:"createTime"`
	Images         []string `json:"images"`
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
	Keyword    string
	CategoryID uint64
	Status     string
	MinPrice   *float64
	MaxPrice   *float64
	Sort       string
	Page       int
	PageSize   int
}

type ProductListResult struct {
	List     []Product `json:"list"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	Total    int       `json:"total"`
}

func (r *Repository) ListProducts(ctx context.Context, input ListProductsInput) (*ProductListResult, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageSize > 50 {
		input.PageSize = 50
	}
	if input.Status == "" {
		input.Status = "ON_SALE"
	}

	where := ` WHERE is_deleted = 0 `
	args := []any{}

	if input.Status != "ALL" {
		where += " AND status = ? "
		args = append(args, input.Status)
	}

	if input.Keyword != "" {
		where += " AND (title LIKE ? OR description LIKE ?) "
		keyword := "%" + input.Keyword + "%"
		args = append(args, keyword, keyword)
	}

	if input.CategoryID > 0 {
		where += " AND category_id = ? "
		args = append(args, input.CategoryID)
	}

	if input.MinPrice != nil {
		where += " AND price >= ? "
		args = append(args, *input.MinPrice)
	}

	if input.MaxPrice != nil {
		where += " AND price <= ? "
		args = append(args, *input.MaxPrice)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products "+where, args...).Scan(&total); err != nil {
		return nil, err
	}

	orderBy := " ORDER BY create_time DESC "
	switch input.Sort {
	case "price_asc":
		orderBy = " ORDER BY price ASC, create_time DESC "
	case "price_desc":
		orderBy = " ORDER BY price DESC, create_time DESC "
	case "popular":
		orderBy = " ORDER BY view_count DESC, favorite_count DESC, create_time DESC "
	}

	offset := (input.Page - 1) * input.PageSize

	query := `
		SELECT id, seller_id, category_id, title, description, original_price,
		       price, condition_level, meet_location, status, view_count,
		       favorite_count, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s')
		FROM products
	` + where + orderBy + ` LIMIT ? OFFSET ?`

	queryArgs := append(args, input.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		var desc sql.NullString
		var originalPrice sql.NullFloat64
		var condition sql.NullString
		var location sql.NullString

		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
			&p.Price, &condition, &location, &p.Status, &p.ViewCount,
			&p.FavoriteCount, &p.CreateTime,
		); err != nil {
			return nil, err
		}

		fillProductNullableFields(&p, desc, originalPrice, condition, location)

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

func (r *Repository) GetProductByID(ctx context.Context, id uint64) (*Product, error) {
	var p Product
	var desc sql.NullString
	var originalPrice sql.NullFloat64
	var condition sql.NullString
	var location sql.NullString

	err := r.db.QueryRowContext(ctx, `
		SELECT id, seller_id, category_id, title, description, original_price,
		       price, condition_level, meet_location, status, view_count,
		       favorite_count, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s')
		FROM products
		WHERE id = ? AND is_deleted = 0 AND status <> 'DELETED'
	`, id).Scan(
		&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
		&p.Price, &condition, &location, &p.Status, &p.ViewCount,
		&p.FavoriteCount, &p.CreateTime,
	)
	if err != nil {
		return nil, err
	}

	fillProductNullableFields(&p, desc, originalPrice, condition, location)

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

	var images []string
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
		       price, condition_level, meet_location, status, view_count,
		       favorite_count, DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s')
		FROM products
		WHERE seller_id = ? AND is_deleted = 0
		ORDER BY create_time DESC
	`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Product
	for rows.Next() {
		var p Product
		var desc sql.NullString
		var originalPrice sql.NullFloat64
		var condition sql.NullString
		var location sql.NullString

		if err := rows.Scan(
			&p.ID, &p.SellerID, &p.CategoryID, &p.Title, &desc, &originalPrice,
			&p.Price, &condition, &location, &p.Status, &p.ViewCount,
			&p.FavoriteCount, &p.CreateTime,
		); err != nil {
			return nil, err
		}

		fillProductNullableFields(&p, desc, originalPrice, condition, location)

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
		    off_shelf_reason = ?,
		    update_time = CURRENT_TIMESTAMP
		WHERE id = ?
		  AND seller_id = ?
		  AND is_deleted = 0
		  AND (
		      (? = 'OFF_SHELF' AND status = 'ON_SALE')
		      OR
		      (? = 'ON_SALE' AND status = 'OFF_SHELF')
		  )
	`,
		status, reason, productID, sellerID, status, status,
	)
	if err != nil {
		return err
	}

	return checkAffected(result)
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

func (r *Repository) LockProduct(ctx context.Context, productID uint64) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET status = 'LOCKED',
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
