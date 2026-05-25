package order

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tx *sql.Tx, o *Order) error {
	query := `
		INSERT INTO orders (order_no, product_id, buyer_id, seller_id, product_title_snapshot, product_price_snapshot, status, remark, meet_time, meet_location, expire_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	args := []interface{}{o.OrderNo, o.ProductID, o.BuyerID, o.SellerID, o.ProductTitleSnapshot, o.ProductPriceSnapshot, o.Status, o.Remark, o.MeetTime, o.MeetLocation, o.ExpireTime}

	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, args...)
	} else {
		result, err = r.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	o.ID = uint64(id)
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64) (*Order, error) {
	query := `
		SELECT id, order_no, product_id, buyer_id, seller_id, product_title_snapshot, product_price_snapshot, status, remark, meet_time, meet_location, cancel_reason, cancel_by, expire_time, confirm_time, finish_time, close_time, create_time, update_time
		FROM orders WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanOrder(row)
}

func (r *Repository) GetByOrderNo(ctx context.Context, orderNo string) (*Order, error) {
	query := `
		SELECT id, order_no, product_id, buyer_id, seller_id, product_title_snapshot, product_price_snapshot, status, remark, meet_time, meet_location, cancel_reason, cancel_by, expire_time, confirm_time, finish_time, close_time, create_time, update_time
		FROM orders WHERE order_no = ?
	`
	row := r.db.QueryRowContext(ctx, query, orderNo)
	return scanOrder(row)
}

func (r *Repository) UpdateStatus(ctx context.Context, tx *sql.Tx, id uint64, status string, updates map[string]interface{}) error {
	query := "UPDATE orders SET status = ?"
	args := []interface{}{status}

	for col, val := range updates {
		query += ", " + col + " = ?"
		args = append(args, val)
	}
	query += " WHERE id = ?"
	args = append(args, id)

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}

func (r *Repository) LockProduct(ctx context.Context, tx *sql.Tx, productID uint64) error {
	query := "UPDATE products SET status = 'LOCKED' WHERE id = ? AND status = 'ON_SALE'"
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, productID)
	} else {
		result, err = r.db.ExecContext(ctx, query, productID)
	}
	if err != nil {
		return fmt.Errorf("lock product: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("product not available")
	}
	return nil
}

func (r *Repository) UnlockProduct(ctx context.Context, tx *sql.Tx, productID uint64) error {
	query := "UPDATE products SET status = 'ON_SALE' WHERE id = ? AND status = 'LOCKED'"
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, productID)
	} else {
		result, err = r.db.ExecContext(ctx, query, productID)
	}
	if err != nil {
		return fmt.Errorf("unlock product: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("product not locked")
	}
	return nil
}

func (r *Repository) MarkProductSold(ctx context.Context, tx *sql.Tx, productID uint64) error {
	query := "UPDATE products SET status = 'SOLD' WHERE id = ? AND status = 'LOCKED'"
	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, productID)
	} else {
		result, err = r.db.ExecContext(ctx, query, productID)
	}
	if err != nil {
		return fmt.Errorf("mark product sold: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("product not locked")
	}
	return nil
}

func (r *Repository) ListByBuyer(ctx context.Context, buyerID uint64, status string, page, pageSize int) ([]OrderDetail, int, error) {
	where := "o.buyer_id = ?"
	args := []interface{}{buyerID}
	if status != "" {
		where += " AND o.status = ?"
		args = append(args, status)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM orders o WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	query := `
		SELECT o.id, o.order_no, o.product_id, o.buyer_id, o.seller_id, o.product_title_snapshot, o.product_price_snapshot, o.status, o.remark, o.meet_time, o.meet_location, o.cancel_reason, o.cancel_by, o.expire_time, o.confirm_time, o.finish_time, o.close_time, o.create_time, o.update_time,
			ub.nickname, us.nickname, pi.image_url
		FROM orders o
		LEFT JOIN users ub ON ub.id = o.buyer_id
		LEFT JOIN users us ON us.id = o.seller_id
		LEFT JOIN product_images pi ON pi.product_id = o.product_id AND pi.sort_order = 0
		WHERE ` + where + `
		ORDER BY o.create_time DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []OrderDetail
	for rows.Next() {
		var od OrderDetail
		var buyerNick, sellerNick, productImg sql.NullString
		err := rows.Scan(
			&od.ID, &od.OrderNo, &od.ProductID, &od.BuyerID, &od.SellerID, &od.ProductTitleSnapshot, &od.ProductPriceSnapshot, &od.Status, &od.Remark, &od.MeetTime, &od.MeetLocation, &od.CancelReason, &od.CancelBy, &od.ExpireTime, &od.ConfirmTime, &od.FinishTime, &od.CloseTime, &od.CreateTime, &od.UpdateTime,
			&buyerNick, &sellerNick, &productImg,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan order: %w", err)
		}
		if buyerNick.Valid {
			od.BuyerNickname = &buyerNick.String
		}
		if sellerNick.Valid {
			od.SellerNickname = &sellerNick.String
		}
		if productImg.Valid {
			od.ProductImage = &productImg.String
		}
		orders = append(orders, od)
	}
	return orders, total, nil
}

func (r *Repository) ListBySeller(ctx context.Context, sellerID uint64, status string, page, pageSize int) ([]OrderDetail, int, error) {
	where := "o.seller_id = ?"
	args := []interface{}{sellerID}
	if status != "" {
		where += " AND o.status = ?"
		args = append(args, status)
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM orders o WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count orders: %w", err)
	}

	query := `
		SELECT o.id, o.order_no, o.product_id, o.buyer_id, o.seller_id, o.product_title_snapshot, o.product_price_snapshot, o.status, o.remark, o.meet_time, o.meet_location, o.cancel_reason, o.cancel_by, o.expire_time, o.confirm_time, o.finish_time, o.close_time, o.create_time, o.update_time,
			ub.nickname, us.nickname, pi.image_url
		FROM orders o
		LEFT JOIN users ub ON ub.id = o.buyer_id
		LEFT JOIN users us ON us.id = o.seller_id
		LEFT JOIN product_images pi ON pi.product_id = o.product_id AND pi.sort_order = 0
		WHERE ` + where + `
		ORDER BY o.create_time DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	defer rows.Close()

	var orders []OrderDetail
	for rows.Next() {
		var od OrderDetail
		var buyerNick, sellerNick, productImg sql.NullString
		err := rows.Scan(
			&od.ID, &od.OrderNo, &od.ProductID, &od.BuyerID, &od.SellerID, &od.ProductTitleSnapshot, &od.ProductPriceSnapshot, &od.Status, &od.Remark, &od.MeetTime, &od.MeetLocation, &od.CancelReason, &od.CancelBy, &od.ExpireTime, &od.ConfirmTime, &od.FinishTime, &od.CloseTime, &od.CreateTime, &od.UpdateTime,
			&buyerNick, &sellerNick, &productImg,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan order: %w", err)
		}
		if buyerNick.Valid {
			od.BuyerNickname = &buyerNick.String
		}
		if sellerNick.Valid {
			od.SellerNickname = &sellerNick.String
		}
		if productImg.Valid {
			od.ProductImage = &productImg.String
		}
		orders = append(orders, od)
	}
	return orders, total, nil
}

func (r *Repository) HasActiveOrderByBuyer(ctx context.Context, buyerID, productID uint64) (bool, error) {
	query := `
		SELECT COUNT(*) FROM orders
		WHERE buyer_id = ? AND product_id = ? AND status IN ('PENDING_CONFIRM', 'WAIT_MEET')
	`
	var count int
	if err := r.db.QueryRowContext(ctx, query, buyerID, productID).Scan(&count); err != nil {
		return false, fmt.Errorf("check active order: %w", err)
	}
	return count > 0, nil
}

func (r *Repository) GetProductSeller(ctx context.Context, productID uint64) (uint64, error) {
	query := "SELECT seller_id FROM products WHERE id = ?"
	var sellerID uint64
	if err := r.db.QueryRowContext(ctx, query, productID).Scan(&sellerID); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("product not found")
		}
		return 0, fmt.Errorf("get product seller: %w", err)
	}
	return sellerID, nil
}

func (r *Repository) GetProductInfo(ctx context.Context, productID uint64) (title string, price float64, err error) {
	query := "SELECT title, price FROM products WHERE id = ?"
	if err := r.db.QueryRowContext(ctx, query, productID).Scan(&title, &price); err != nil {
		if err == sql.ErrNoRows {
			return "", 0, fmt.Errorf("product not found")
		}
		return "", 0, fmt.Errorf("get product info: %w", err)
	}
	return title, price, nil
}

func (r *Repository) ListExpiredOrders(ctx context.Context) ([]Order, error) {
	query := `
		SELECT id, order_no, product_id, buyer_id, seller_id, product_title_snapshot, product_price_snapshot, status, remark, meet_time, meet_location, cancel_reason, cancel_by, expire_time, confirm_time, finish_time, close_time, create_time, update_time
		FROM orders
		WHERE status = 'PENDING_CONFIRM' AND expire_time < NOW()
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list expired orders: %w", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		var remark, meetTime, meetLocation, cancelReason sql.NullString
		var cancelBy sql.NullInt64
		var confirmTime, finishTime, closeTime sql.NullString

		err := rows.Scan(
			&o.ID, &o.OrderNo, &o.ProductID, &o.BuyerID, &o.SellerID, &o.ProductTitleSnapshot, &o.ProductPriceSnapshot, &o.Status,
			&remark, &meetTime, &meetLocation, &cancelReason, &cancelBy,
			&o.ExpireTime, &confirmTime, &finishTime, &closeTime, &o.CreateTime, &o.UpdateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scan expired order: %w", err)
		}
		if remark.Valid {
			o.Remark = &remark.String
		}
		if meetTime.Valid {
			o.MeetTime = &meetTime.String
		}
		if meetLocation.Valid {
			o.MeetLocation = &meetLocation.String
		}
		if cancelReason.Valid {
			o.CancelReason = &cancelReason.String
		}
		if cancelBy.Valid {
			v := uint64(cancelBy.Int64)
			o.CancelBy = &v
		}
		if confirmTime.Valid {
			o.ConfirmTime = &confirmTime.String
		}
		if finishTime.Valid {
			o.FinishTime = &finishTime.String
		}
		if closeTime.Valid {
			o.CloseTime = &closeTime.String
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func scanOrder(row *sql.Row) (*Order, error) {
	var o Order
	var remark, meetTime, meetLocation, cancelReason sql.NullString
	var cancelBy sql.NullInt64
	var confirmTime, finishTime, closeTime sql.NullString

	err := row.Scan(
		&o.ID, &o.OrderNo, &o.ProductID, &o.BuyerID, &o.SellerID, &o.ProductTitleSnapshot, &o.ProductPriceSnapshot, &o.Status,
		&remark, &meetTime, &meetLocation, &cancelReason, &cancelBy,
		&o.ExpireTime, &confirmTime, &finishTime, &closeTime, &o.CreateTime, &o.UpdateTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan order: %w", err)
	}

	if remark.Valid {
		o.Remark = &remark.String
	}
	if meetTime.Valid {
		o.MeetTime = &meetTime.String
	}
	if meetLocation.Valid {
		o.MeetLocation = &meetLocation.String
	}
	if cancelReason.Valid {
		o.CancelReason = &cancelReason.String
	}
	if cancelBy.Valid {
		v := uint64(cancelBy.Int64)
		o.CancelBy = &v
	}
	if confirmTime.Valid {
		o.ConfirmTime = &confirmTime.String
	}
	if finishTime.Valid {
		o.FinishTime = &finishTime.String
	}
	if closeTime.Valid {
		o.CloseTime = &closeTime.String
	}

	return &o, nil
}

func generateOrderNo() string {
	return fmt.Sprintf("O%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000000)
}
