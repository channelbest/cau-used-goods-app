package user

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"cau-used-goods-app/backend/internal/admin"
	"cau-used-goods-app/backend/internal/message"
	"cau-used-goods-app/backend/internal/order"
)

const (
	accountStatusNormal   = "NORMAL"
	accountStatusDisabled = "DISABLED"
	accountStatusBanned   = "BANNED"
	accountStatusCanceled = "CANCELED"
	authStatusPending     = "PENDING"
	authStatusVerified    = "VERIFIED"
	authStatusRejected    = "REJECTED"
	roleUser              = "USER"
	roleAdmin             = "ADMIN"
	roleSuperAdmin        = "SUPER_ADMIN"
)

var (
	phonePattern     = regexp.MustCompile(`^1[3-9]\d{9}$`)
	studentIDPattern = regexp.MustCompile(`^[A-Za-z0-9]{6,30}$`)
)

type Service struct {
	repo     *Repository
	products ProductLifecycle
	orders   OrderLifecycle
	messages MessageNotifier
	admin    *admin.Service
}

type ProductLifecycle interface {
	OffShelfOnSaleBySellerTx(ctx context.Context, tx *sql.Tx, sellerID uint64, reason string) ([]uint64, error)
}

type OrderLifecycle interface {
	CountBlockingOrdersTx(ctx context.Context, tx *sql.Tx, userID uint64) (sellerPendingConfirm, sellerWaitMeet, buyerPendingConfirm, buyerWaitMeet int, err error)
	AutoExceptionClosePendingConfirmByUserTx(ctx context.Context, tx *sql.Tx, userID, adminID uint64, reason string, ipAddress *string, relatedType *string, relatedID *uint64) ([]order.AccountStatusClosedOrder, error)
}

type MessageNotifier interface {
	Create(ctx context.Context, input message.CreateMessageInput) (uint64, error)
}

type UpdateProfileInput struct {
	Nickname  *string
	AvatarURL *string
	Phone     *string
}

type SubmitStudentVerificationInput struct {
	StudentID string
	RealName  string
	College   string
}

type ReviewStudentVerificationInput struct {
	UserID      uint64
	AuthStatus  string
	Description string
}

type UpdateAccountStatusInput struct {
	UserID        uint64
	AccountStatus string
	Reason        string
	IPAddress     string
	RelatedType   string
	RelatedID     uint64
}

type AccountStatusChangeEffects struct {
	ClosedOrders []order.AccountStatusClosedOrder
	OffShelfProducts []uint64
}

type UpdateRoleInput struct {
	UserID    uint64
	Role      string
	Reason    string
	IPAddress string
}

func NewService(repo *Repository, products ProductLifecycle, orders OrderLifecycle, messageService MessageNotifier, adminService *admin.Service) *Service {
	return &Service{repo: repo, products: products, orders: orders, messages: messageService, admin: adminService}
}

func (s *Service) Me(ctx context.Context, userID uint64) (*User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Service) ListAdminUsers(ctx context.Context, query AdminUserQuery, includeOpenID bool) (*PagedResult[AdminUser], error) {
	query.Keyword = strings.TrimSpace(query.Keyword)
	query.AuthStatus = strings.ToUpper(strings.TrimSpace(query.AuthStatus))
	query.AccountStatus = strings.ToUpper(strings.TrimSpace(query.AccountStatus))
	query.Role = strings.ToUpper(strings.TrimSpace(query.Role))
	query.Page, query.PageSize = normalizePage(query.Page, query.PageSize)
	query.IncludeOpenID = includeOpenID

	if query.AuthStatus != "" && !isValidAuthStatus(query.AuthStatus) {
		return nil, fmt.Errorf("invalid authStatus")
	}
	if query.AccountStatus != "" && !isValidAccountStatus(query.AccountStatus) {
		return nil, fmt.Errorf("invalid accountStatus")
	}
	if query.Role != "" && !isValidRole(query.Role) {
		return nil, fmt.Errorf("invalid role")
	}

	items, total, err := s.repo.ListAdminUsers(ctx, query)
	if err != nil {
		return nil, err
	}
	if !includeOpenID {
		for i := range items {
			items[i].OpenID = nil
		}
	}
	return &PagedResult[AdminUser]{Items: items, Total: total, Page: query.Page, PageSize: query.PageSize}, nil
}

func (s *Service) AdminUserDetail(ctx context.Context, userID uint64, includeOpenID, includeIP bool) (*AdminUserDetail, error) {
	if userID == 0 {
		return nil, fmt.Errorf("userId is required")
	}
	user, err := s.repo.FindAdminUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if !includeOpenID {
		user.OpenID = nil
	}
	stats, err := s.repo.UserStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	reports, _, err := s.repo.ListUserReports(ctx, userID, 1, 5)
	if err != nil {
		return nil, err
	}
	appeals, _, err := s.repo.ListUserAppeals(ctx, userID, 1, 5)
	if err != nil {
		return nil, err
	}
	logs, _, err := s.repo.ListUserLogs(ctx, userID, 1, 5)
	if err != nil {
		return nil, err
	}
	if !includeIP {
		hideLogIPs(logs)
	}
	return &AdminUserDetail{
		User:            user,
		Stats:           stats,
		RecentReports:   reports,
		RecentAppeals:   appeals,
		RecentAdminLogs: logs,
	}, nil
}

func (s *Service) UpdateAccountStatus(ctx context.Context, adminID uint64, input UpdateAccountStatusInput) (*AdminUser, error) {
	input.AccountStatus = strings.ToUpper(strings.TrimSpace(input.AccountStatus))
	input.Reason = strings.TrimSpace(input.Reason)
	input.IPAddress = strings.TrimSpace(input.IPAddress)
	input.RelatedType = strings.ToUpper(strings.TrimSpace(input.RelatedType))
	if input.UserID == 0 {
		return nil, fmt.Errorf("userId is required")
	}
	if input.UserID == adminID {
		return nil, fmt.Errorf("管理员不能修改自己的账号状态")
	}
	if input.AccountStatus != accountStatusNormal && input.AccountStatus != accountStatusDisabled && input.AccountStatus != accountStatusBanned {
		return nil, fmt.Errorf("accountStatus 只能是 NORMAL、DISABLED 或 BANNED")
	}
	if input.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}
	if len([]rune(input.Reason)) > 500 {
		return nil, fmt.Errorf("reason cannot exceed 500 characters")
	}
	if (input.RelatedType == "") != (input.RelatedID == 0) {
		return nil, fmt.Errorf("relatedType 和 relatedId 必须同时提供")
	}
	if input.RelatedType != "" && input.RelatedType != "REPORT" && input.RelatedType != "APPEAL" {
		return nil, fmt.Errorf("relatedType 只能是 REPORT 或 APPEAL")
	}
	adminUser, err := s.repo.FindAdminUserByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if adminUser == nil {
		return nil, fmt.Errorf("需要管理员权限")
	}
	effects, err := s.repo.ChangeAccountStatus(ctx, adminID, adminUser.Role, input, s.products, s.orders, s.admin)
	if err != nil {
		return nil, err
	}
	s.notifyAccountStatusChange(ctx, adminID, input, effects)
	result, err := s.repo.FindAdminUserByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if adminUser.Role != roleSuperAdmin && result != nil {
		result.OpenID = nil
	}
	return result, nil
}

func (s *Service) notifyAccountStatusChange(ctx context.Context, adminID uint64, input UpdateAccountStatusInput, effects *AccountStatusChangeEffects) {
	if s.messages == nil {
		return
	}
	if effects != nil {
		relatedType := message.RelatedTypeOrder
		for _, item := range effects.ClosedOrders {
			relatedID := item.ID
			content := fmt.Sprintf("订单「%s」因相关账号状态异常，已由管理员异常关闭。", item.ProductTitleSnapshot)
			for _, receiverID := range []uint64{item.BuyerID, item.SellerID} {
				_, _ = s.messages.Create(ctx, message.CreateMessageInput{
					ReceiverID:  receiverID,
					SenderID:    &adminID,
					MessageType: message.MessageTypeSystemNotice,
					Title:       "订单异常关闭",
					Content:     content,
					RelatedType: &relatedType,
					RelatedID:   &relatedID,
				})
			}
		}
	}

	relatedType := message.RelatedTypeUser
	relatedID := input.UserID
	_, _ = s.messages.Create(ctx, message.CreateMessageInput{
		ReceiverID:  input.UserID,
		SenderID:    &adminID,
		MessageType: message.MessageTypeSystemNotice,
		Title:       "账号状态变更",
		Content:     buildAccountStatusMessage(input.AccountStatus, input.Reason, effects),
		RelatedType: &relatedType,
		RelatedID:   &relatedID,
	})
}

func buildAccountStatusMessage(status string, reason string, effects *AccountStatusChangeEffects) string {
	var content string
	switch status {
	case accountStatusDisabled:
		content = fmt.Sprintf("你的账号已被临时禁用，原因：%s", reason)
	case accountStatusBanned:
		content = fmt.Sprintf("你的账号已被永久封禁，原因：%s", reason)
	case accountStatusNormal:
		content = fmt.Sprintf("你的账号状态已恢复正常，原因：%s", reason)
	default:
		content = fmt.Sprintf("你的账号状态已变更为 %s，原因：%s", status, reason)
	}
	if effects != nil {
		if len(effects.OffShelfProducts) > 0 {
			content += fmt.Sprintf("；你发布的 %d 件商品已被下架", len(effects.OffShelfProducts))
		}
		if len(effects.ClosedOrders) > 0 {
			content += fmt.Sprintf("；%d 个进行中的订单已被关闭", len(effects.ClosedOrders))
		}
	}
	return limitRunes(content, 500)
}

func limitRunes(value string, limit int) string {
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit])
}

func (s *Service) UpdateRole(ctx context.Context, adminID uint64, input UpdateRoleInput) (*AdminUser, error) {
	input.Role = strings.ToUpper(strings.TrimSpace(input.Role))
	input.Reason = strings.TrimSpace(input.Reason)
	input.IPAddress = strings.TrimSpace(input.IPAddress)
	if input.UserID == 0 {
		return nil, fmt.Errorf("userId is required")
	}
	if input.UserID == adminID {
		return nil, fmt.Errorf("管理员不能修改自己的角色")
	}
	adminUser, err := s.repo.FindAdminUserByID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if adminUser == nil || adminUser.Role != roleSuperAdmin {
		return nil, fmt.Errorf("需要超级管理员权限")
	}
	if input.Role != roleUser && input.Role != roleAdmin {
		return nil, fmt.Errorf("role 只能是 USER 或 ADMIN")
	}
	if input.Reason == "" {
		return nil, fmt.Errorf("reason is required")
	}
	if len([]rune(input.Reason)) > 500 {
		return nil, fmt.Errorf("reason cannot exceed 500 characters")
	}
	if err := s.repo.ChangeRole(ctx, adminID, input.UserID, input.Role, input.Reason, input.IPAddress, s.admin); err != nil {
		return nil, err
	}
	return s.repo.FindAdminUserByID(ctx, input.UserID)
}

func (s *Service) ListUserProducts(ctx context.Context, userID uint64, page, pageSize int, dataScope string) (*PagedResult[UserProductItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	dataScope, err := normalizeDataScope(dataScope)
	if err != nil {
		return nil, err
	}
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserProducts(ctx, userID, page, pageSize, dataScope)
	if err != nil {
		return nil, err
	}
	return &PagedResult[UserProductItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) ListUserOrders(ctx context.Context, userID uint64, page, pageSize int) (*PagedResult[UserOrderItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserOrders(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PagedResult[UserOrderItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) ListUserReports(ctx context.Context, userID uint64, page, pageSize int) (*PagedResult[UserReportItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserReports(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PagedResult[UserReportItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) ListUserAppeals(ctx context.Context, userID uint64, page, pageSize int) (*PagedResult[UserAppealItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserAppeals(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PagedResult[UserAppealItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) ListUserReviews(ctx context.Context, userID uint64, page, pageSize int, dataScope string) (*PagedResult[UserReviewItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	dataScope, err := normalizeDataScope(dataScope)
	if err != nil {
		return nil, err
	}
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserReviews(ctx, userID, page, pageSize, dataScope)
	if err != nil {
		return nil, err
	}
	return &PagedResult[UserReviewItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) ListUserLogs(ctx context.Context, userID uint64, page, pageSize int, includeIP bool) (*PagedResult[UserLogItem], error) {
	page, pageSize = normalizePage(page, pageSize)
	if err := s.ensureAdminUserExists(ctx, userID); err != nil {
		return nil, err
	}
	items, total, err := s.repo.ListUserLogs(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	if !includeIP {
		hideLogIPs(items)
	}
	return &PagedResult[UserLogItem]{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) PublicProfile(ctx context.Context, userID uint64) (*PublicProfile, error) {
	item, err := s.repo.FindPublicProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return item, nil
}

func (s *Service) PublicHomepage(ctx context.Context, userID uint64, page, pageSize int) (*PublicHomepage, error) {
	profile, err := s.PublicProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	page, pageSize = normalizePage(page, pageSize)

	stats, err := s.repo.PublicHomepageStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	products, total, err := s.repo.ListPublicHomepageProducts(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	return &PublicHomepage{
		Profile: profile,
		Stats:   stats,
		Products: PagedResult[UserProductItem]{
			Items:    products,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	}, nil
}

func (s *Service) Restriction(ctx context.Context, userID uint64) (*Restriction, error) {
	item, err := s.repo.FindRestriction(ctx, userID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	return item, nil
}

func (s *Service) CancelAccount(ctx context.Context, userID uint64, confirmed bool) error {
	if !confirmed {
		return fmt.Errorf("必须确认注销账号")
	}
	return s.repo.CancelAccount(ctx, userID, s.products, s.orders)
}

func (s *Service) ReactivateAccount(ctx context.Context, userID uint64) error {
	return s.repo.ReactivateAccount(ctx, userID)
}

func (s *Service) EnsureAccountNormal(ctx context.Context, userID uint64) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}
	if user.AccountStatus != accountStatusNormal {
		return fmt.Errorf("当前账号状态不可操作")
	}
	return nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID uint64, input UpdateProfileInput) (*User, error) {
	if input.Nickname == nil && input.AvatarURL == nil && input.Phone == nil {
		return nil, fmt.Errorf("nothing to update")
	}
	if err := s.EnsureAccountNormal(ctx, userID); err != nil {
		return nil, err
	}

	if input.Nickname != nil {
		trimmed := strings.TrimSpace(*input.Nickname)
		if err := validateStringLength("昵称", trimmed, 1, 50); err != nil {
			return nil, err
		}
		input.Nickname = &trimmed
	}
	if input.AvatarURL != nil {
		trimmed := strings.TrimSpace(*input.AvatarURL)
		if err := validateStringLength("头像地址", trimmed, 1, 255); err != nil {
			return nil, err
		}
		if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") && !strings.HasPrefix(trimmed, "/uploads/avatar/") {
			return nil, fmt.Errorf("头像地址必须是 http(s) URL 或 /uploads/avatar/ 路径")
		}
		input.AvatarURL = &trimmed
	}
	if input.Phone != nil {
		trimmed := strings.TrimSpace(*input.Phone)
		if !phonePattern.MatchString(trimmed) {
			return nil, fmt.Errorf("手机号格式不正确")
		}
		input.Phone = &trimmed
	}

	if err := s.repo.UpdateProfile(ctx, userID, input.Nickname, input.AvatarURL, input.Phone); err != nil {
		return nil, err
	}
	return s.Me(ctx, userID)
}

func (s *Service) SubmitStudentVerification(ctx context.Context, userID uint64, input SubmitStudentVerificationInput) (*StudentVerification, error) {
	input.StudentID = strings.TrimSpace(input.StudentID)
	input.RealName = strings.TrimSpace(input.RealName)
	input.College = strings.TrimSpace(input.College)

	if !studentIDPattern.MatchString(input.StudentID) {
		return nil, fmt.Errorf("学号格式不正确")
	}
	if err := validateStringLength("真实姓名", input.RealName, 2, 30); err != nil {
		return nil, err
	}
	if err := validateStringLength("学院名称", input.College, 1, 50); err != nil {
		return nil, err
	}

	current, err := s.repo.FindStudentVerification(ctx, userID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if current.AccountStatus != accountStatusNormal {
		return nil, fmt.Errorf("当前账号状态不可操作")
	}
	switch current.AuthStatus {
	case authStatusVerified:
		return nil, fmt.Errorf("学生认证已通过，不能重复提交")
	case authStatusPending:
		return nil, fmt.Errorf("学生认证正在审核中，请勿重复提交")
	case "UNVERIFIED", authStatusRejected:
	default:
		return nil, fmt.Errorf("当前认证状态不可提交")
	}

	used, err := s.repo.StudentIDUsedByOther(ctx, input.StudentID, userID)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, fmt.Errorf("学号已被使用")
	}

	if err := s.repo.SubmitStudentVerification(ctx, userID, input.StudentID, input.RealName, input.College); err != nil {
		return nil, err
	}
	return s.StudentVerification(ctx, userID)
}

func (s *Service) StudentVerification(ctx context.Context, userID uint64) (*StudentVerification, error) {
	verification, err := s.repo.FindStudentVerification(ctx, userID)
	if err != nil {
		return nil, err
	}
	if verification == nil {
		return nil, fmt.Errorf("user not found")
	}
	return verification, nil
}

func (s *Service) ListStudentVerifications(ctx context.Context, status string) ([]StudentVerification, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		status = "PENDING"
	}
	switch status {
	case "PENDING", "VERIFIED", "REJECTED", "UNVERIFIED":
	default:
		return nil, fmt.Errorf("invalid authStatus")
	}
	return s.repo.ListStudentVerifications(ctx, status)
}

func (s *Service) ReviewStudentVerification(ctx context.Context, adminID uint64, input ReviewStudentVerificationInput) (*StudentVerification, error) {
	if input.UserID == 0 {
		return nil, fmt.Errorf("userId is required")
	}
	if input.UserID == adminID {
		return nil, fmt.Errorf("管理员不能审核自己的认证")
	}
	input.AuthStatus = strings.TrimSpace(input.AuthStatus)
	input.Description = strings.TrimSpace(input.Description)

	if input.AuthStatus != authStatusVerified && input.AuthStatus != authStatusRejected {
		return nil, fmt.Errorf("authStatus must be VERIFIED or REJECTED")
	}
	if len(input.Description) > 500 {
		return nil, fmt.Errorf("审核说明不能超过 500 个字符")
	}
	if input.AuthStatus == authStatusRejected && input.Description == "" {
		return nil, fmt.Errorf("驳回时必须填写审核说明")
	}
	if input.AuthStatus == authStatusVerified && input.Description == "" {
		input.Description = "学生认证审核通过"
	}

	target, err := s.repo.FindStudentVerification(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, fmt.Errorf("用户不存在")
	}
	if target.AccountStatus != accountStatusNormal {
		return nil, fmt.Errorf("目标用户账号状态不可审核")
	}
	if target.AuthStatus != authStatusPending {
		return nil, fmt.Errorf("认证状态已变化，请刷新后重试")
	}

	if err := s.repo.ReviewStudentVerification(ctx, adminID, input.UserID, input.AuthStatus, input.Description, s.admin); err != nil {
		return nil, err
	}
	s.notifyStudentVerificationResult(ctx, adminID, input)
	return s.StudentVerification(ctx, input.UserID)
}

func (s *Service) notifyStudentVerificationResult(ctx context.Context, adminID uint64, input ReviewStudentVerificationInput) {
	if s.messages == nil {
		return
	}

	title := "学生认证审核结果"
	content := buildStudentVerificationResultContent(input.AuthStatus, input.Description)

	relatedType := message.RelatedTypeUser
	relatedID := input.UserID
	_, _ = s.messages.Create(ctx, message.CreateMessageInput{
		ReceiverID:  input.UserID,
		SenderID:    &adminID,
		MessageType: message.MessageTypeSystemNotice,
		Title:       title,
		Content:     limitRunes(content, 500),
		RelatedType: &relatedType,
		RelatedID:   &relatedID,
	})
}

func buildStudentVerificationResultContent(authStatus string, description string) string {
	description = strings.TrimSpace(description)
	if authStatus == authStatusVerified {
		if description == "" || description == "学生认证审核通过" {
			return "你的学生认证已通过，现在可以正常使用发布、收藏、预约和举报等功能。"
		}
		return "你的学生认证已通过。" + description
	}
	if description == "" {
		return "你的学生认证未通过，请核对资料后重新提交。"
	}
	return "你的学生认证未通过。" + description
}

func validateStringLength(field string, value string, min int, max int) error {
	length := len([]rune(value))
	if length < min || length > max {
		if min == max {
			return fmt.Errorf("%s长度应为 %d 个字符", field, min)
		}
		return fmt.Errorf("%s长度应为 %d-%d 个字符", field, min, max)
	}
	return nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func isValidAuthStatus(status string) bool {
	return status == "UNVERIFIED" ||
		status == authStatusPending ||
		status == authStatusVerified ||
		status == authStatusRejected
}

func isValidAccountStatus(status string) bool {
	return status == accountStatusNormal ||
		status == accountStatusDisabled ||
		status == accountStatusBanned ||
		status == accountStatusCanceled
}

func isValidRole(role string) bool {
	return role == roleUser || role == roleAdmin || role == roleSuperAdmin
}

func (s *Service) ensureAdminUserExists(ctx context.Context, userID uint64) error {
	item, err := s.repo.FindAdminUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if item == nil {
		return fmt.Errorf("用户不存在")
	}
	return nil
}

func hideLogIPs(items []UserLogItem) {
	for i := range items {
		items[i].IPAddress = nil
	}
}

func normalizeDataScope(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "ACTIVE", nil
	}
	if value != "ACTIVE" && value != "DELETED" && value != "ALL" {
		return "", fmt.Errorf("dataScope 只能是 ACTIVE、DELETED 或 ALL")
	}
	return value, nil
}
