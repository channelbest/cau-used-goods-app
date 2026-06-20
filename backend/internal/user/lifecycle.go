package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cau-used-goods-app/backend/internal/admin"
)

func (r *Repository) ChangeAccountStatus(
	ctx context.Context,
	adminID uint64,
	adminRole string,
	input UpdateAccountStatusInput,
	products ProductLifecycle,
	orders OrderLifecycle,
	adminLogger *admin.Service,
) (*AccountStatusChangeEffects, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin change account status tx: %w", err)
	}
	defer tx.Rollback()
	effects := &AccountStatusChangeEffects{}

	var currentRole, currentStatus string
	var isDeleted bool
	if err := tx.QueryRowContext(ctx, `
		SELECT role, account_status, is_deleted
		FROM users
		WHERE id = ?
		FOR UPDATE
	`, input.UserID).Scan(&currentRole, &currentStatus, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("lock target user: %w", err)
	}

	if isDeleted || currentStatus == accountStatusCanceled {
		return nil, fmt.Errorf("目标用户账号状态不可操作")
	}
	if currentRole == roleSuperAdmin {
		return nil, fmt.Errorf("不能修改超级管理员账号状态")
	}
	if adminRole == roleAdmin && currentRole != roleUser {
		return nil, fmt.Errorf("普通管理员只能修改普通用户账号状态")
	}
	if adminRole != roleAdmin && adminRole != roleSuperAdmin {
		return nil, fmt.Errorf("需要管理员权限")
	}
	if !allowedAccountTransition(currentStatus, input.AccountStatus, adminRole) {
		return nil, fmt.Errorf("不允许的账号状态流转")
	}
	if currentStatus == accountStatusBanned && input.AccountStatus == accountStatusNormal &&
		(input.RelatedType != "APPEAL" || input.RelatedID == 0) {
		return nil, fmt.Errorf("撤销永久封禁必须关联已通过的申诉")
	}
	if input.RelatedType != "" {
		if err := validateRelatedRecordTx(ctx, tx, input.RelatedType, input.RelatedID, input.UserID); err != nil {
			return nil, err
		}
	}
	var relatedType *string
	var relatedID *uint64
	if input.RelatedType != "" {
		relatedType = &input.RelatedType
		relatedID = &input.RelatedID
	}
	if adminLogger == nil {
		return nil, fmt.Errorf("admin logger is not configured")
	}

	if input.AccountStatus == accountStatusDisabled || input.AccountStatus == accountStatusBanned {
		if orders == nil {
			return nil, fmt.Errorf("订单状态协同服务不可用")
		}
		sellerPending, sellerWait, buyerPending, buyerWait, err := orders.CountBlockingOrdersTx(ctx, tx, input.UserID)
		if err != nil {
			return nil, err
		}
		if msg := waitMeetOrderMessage(sellerWait, buyerWait); msg != "" {
			return nil, fmt.Errorf("%s，请先处理待面交订单后再修改账号状态", msg)
		}
		if sellerPending > 0 || buyerPending > 0 {
			closedOrders, err := orders.AutoExceptionClosePendingConfirmByUserTx(ctx, tx, input.UserID, adminID, input.Reason, optionalStringPtr(input.IPAddress), relatedType, relatedID)
			if err != nil {
				return nil, err
			}
			effects.ClosedOrders = closedOrders
		}
		if products == nil {
			return nil, fmt.Errorf("商品状态协同服务不可用")
		}
		productIDs, err := products.OffShelfOnSaleBySellerTx(ctx, tx, input.UserID, input.Reason)
		if err != nil {
			return nil, err
		}
		effects.OffShelfProducts = productIDs
		for _, productID := range productIDs {
			description := fmt.Sprintf("off shelf product due account status %s: %s", input.AccountStatus, input.Reason)
			if _, err := adminLogger.LogActionTx(ctx, tx, admin.LogActionInput{
				AdminID:       adminID,
				OperationType: admin.OperationUpdateProductStatus,
				TargetType:    admin.TargetTypeProduct,
				TargetID:      productID,
				Description:   &description,
				IPAddress:     optionalStringPtr(input.IPAddress),
				RelatedType:   relatedType,
				RelatedID:     relatedID,
			}); err != nil {
				return nil, fmt.Errorf("create product status admin log: %w", err)
			}
		}
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET account_status = ?, is_deleted = 0, update_time = NOW()
		WHERE id = ?
	`, input.AccountStatus, input.UserID); err != nil {
		return nil, fmt.Errorf("update account status: %w", err)
	}

	if _, err := adminLogger.LogActionTx(ctx, tx, admin.LogActionInput{
		AdminID:       adminID,
		OperationType: accountOperationType(currentStatus, input.AccountStatus),
		TargetType:    admin.TargetTypeUser,
		TargetID:      input.UserID,
		Description:   &input.Reason,
		IPAddress:     optionalStringPtr(input.IPAddress),
		RelatedType:   relatedType,
		RelatedID:     relatedID,
	}); err != nil {
		return nil, fmt.Errorf("create account status admin log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit change account status tx: %w", err)
	}
	return effects, nil
}

func (r *Repository) ChangeRole(ctx context.Context, adminID, userID uint64, role, reason, ipAddress string, adminLogger *admin.Service) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin change role tx: %w", err)
	}
	defer tx.Rollback()

	var currentRole, accountStatus string
	var isDeleted bool
	if err := tx.QueryRowContext(ctx, `
		SELECT role, account_status, is_deleted
		FROM users
		WHERE id = ?
		FOR UPDATE
	`, userID).Scan(&currentRole, &accountStatus, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("lock target user role: %w", err)
	}
	if isDeleted || accountStatus != accountStatusNormal {
		return fmt.Errorf("只有正常账号可以修改角色")
	}
	if currentRole == roleSuperAdmin {
		return fmt.Errorf("不能通过接口修改超级管理员角色")
	}
	if role != roleUser && role != roleAdmin {
		return fmt.Errorf("role 只能是 USER 或 ADMIN")
	}
	if currentRole == role {
		return fmt.Errorf("目标用户已经是该角色")
	}

	if _, err := tx.ExecContext(ctx, `UPDATE users SET role = ?, update_time = NOW() WHERE id = ?`, role, userID); err != nil {
		return fmt.Errorf("update user role: %w", err)
	}
	description := fmt.Sprintf("%s -> %s；原因：%s", currentRole, role, reason)
	if adminLogger == nil {
		return fmt.Errorf("admin logger is not configured")
	}
	if _, err := adminLogger.LogActionTx(ctx, tx, admin.LogActionInput{
		AdminID:       adminID,
		OperationType: admin.OperationUserRoleChange,
		TargetType:    admin.TargetTypeUser,
		TargetID:      userID,
		Description:   &description,
		IPAddress:     optionalStringPtr(ipAddress),
	}); err != nil {
		return fmt.Errorf("create role admin log: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit change role tx: %w", err)
	}
	return nil
}

func (r *Repository) CancelAccount(ctx context.Context, userID uint64, products ProductLifecycle, orders OrderLifecycle) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin cancel account tx: %w", err)
	}
	defer tx.Rollback()

	var role, status string
	var isDeleted bool
	if err := tx.QueryRowContext(ctx, `
		SELECT role, account_status, is_deleted
		FROM users
		WHERE id = ?
		FOR UPDATE
	`, userID).Scan(&role, &status, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("lock canceling user: %w", err)
	}
	if role != roleUser {
		return fmt.Errorf("只有普通用户可以主动注销")
	}
	if isDeleted || status != accountStatusNormal {
		return fmt.Errorf("当前账号状态不可注销")
	}
	if orders == nil || products == nil {
		return fmt.Errorf("账号注销协同服务不可用")
	}
	sellerPending, sellerWait, buyerPending, buyerWait, err := orders.CountBlockingOrdersTx(ctx, tx, userID)
	if err != nil {
		return err
	}
	if blockingOrderMessage(sellerPending, sellerWait, buyerPending, buyerWait) != "" {
		return fmt.Errorf("存在进行中订单，暂不能注销")
	}
	if _, err := products.OffShelfOnSaleBySellerTx(ctx, tx, userID, "用户主动注销"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET account_status = 'CANCELED', is_deleted = 1, token_version = token_version + 1, update_time = NOW()
		WHERE id = ?
	`, userID); err != nil {
		return fmt.Errorf("cancel account: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit cancel account tx: %w", err)
	}
	return nil
}

func (r *Repository) ReactivateAccount(ctx context.Context, userID uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reactivate account tx: %w", err)
	}
	defer tx.Rollback()

	var status string
	var isDeleted bool
	if err := tx.QueryRowContext(ctx, `
		SELECT account_status, is_deleted
		FROM users
		WHERE id = ?
		FOR UPDATE
	`, userID).Scan(&status, &isDeleted); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("lock reactivating user: %w", err)
	}
	if status != accountStatusCanceled || !isDeleted {
		return fmt.Errorf("当前账号状态不可重新激活")
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET account_status = 'NORMAL', is_deleted = 0, role = 'USER',
		    token_version = token_version + 1, update_time = NOW()
		WHERE id = ?
	`, userID); err != nil {
		return fmt.Errorf("reactivate account: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reactivate account tx: %w", err)
	}
	return nil
}

func allowedAccountTransition(current, target, adminRole string) bool {
	switch current {
	case accountStatusNormal:
		return target == accountStatusDisabled || target == accountStatusBanned
	case accountStatusDisabled:
		return target == accountStatusNormal || target == accountStatusBanned
	case accountStatusBanned:
		return target == accountStatusNormal && adminRole == roleSuperAdmin
	default:
		return false
	}
}

func accountOperationType(current, target string) string {
	switch target {
	case accountStatusDisabled:
		return admin.OperationUserDisable
	case accountStatusBanned:
		return admin.OperationUserBan
	case accountStatusNormal:
		if current == accountStatusBanned {
			return admin.OperationUserUnban
		}
		return admin.OperationUserEnable
	default:
		return "USER_STATUS_CHANGE"
	}
}

func validateRelatedRecordTx(ctx context.Context, tx *sql.Tx, relatedType string, relatedID, userID uint64) error {
	var query string
	var args []interface{}
	switch relatedType {
	case "REPORT":
		query = "SELECT COUNT(*) FROM reports WHERE id = ? AND target_type = 'USER' AND target_id = ?"
		args = []interface{}{relatedID, userID}
	case "APPEAL":
		query = `SELECT COUNT(*) FROM appeals
			WHERE id = ? AND appellant_id = ? AND target_type = 'USER' AND target_id = ? AND status = 'APPROVED'`
		args = []interface{}{relatedID, userID, userID}
	default:
		return fmt.Errorf("relatedType 只能是 REPORT 或 APPEAL")
	}
	var count int
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return fmt.Errorf("check related record: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("操作依据不存在")
	}
	return nil
}

func nullableUint64(value uint64) interface{} {
	if value == 0 {
		return nil
	}
	return value
}

func blockingOrderMessage(sellerPending, sellerWait, buyerPending, buyerWait int) string {
	var parts []string
	if sellerPending > 0 {
		parts = append(parts, fmt.Sprintf("作为卖家存在%d个待确认订单", sellerPending))
	}
	if sellerWait > 0 {
		parts = append(parts, fmt.Sprintf("作为卖家存在%d个待面交订单", sellerWait))
	}
	if buyerPending > 0 {
		parts = append(parts, fmt.Sprintf("作为买家存在%d个待确认订单", buyerPending))
	}
	if buyerWait > 0 {
		parts = append(parts, fmt.Sprintf("作为买家存在%d个待面交订单", buyerWait))
	}
	if len(parts) == 0 {
		return ""
	}
	return "该用户" + strings.Join(parts, "，")
}

func waitMeetOrderMessage(sellerWait, buyerWait int) string {
	return blockingOrderMessage(0, sellerWait, 0, buyerWait)
}

func optionalStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
