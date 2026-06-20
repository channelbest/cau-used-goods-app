package admin

import "time"

const (
	TargetTypeUser     = "USER"
	TargetTypeProduct  = "PRODUCT"
	TargetTypeOrder    = "ORDER"
	TargetTypeReport   = "REPORT"
	TargetTypeNotice   = "NOTICE"
	TargetTypeWord     = "WORD"
	TargetTypeAppeal   = "APPEAL"
	TargetTypeCategory = "CATEGORY"
)

const (
	AnnouncementStatusDraft     = "DRAFT"
	AnnouncementStatusPublished = "PUBLISHED"
	AnnouncementStatusOffline   = "OFFLINE"
)

const (
	OperationCreateNotice         = "CREATE_NOTICE"
	OperationUpdateNotice         = "UPDATE_NOTICE"
	OperationStatusNotice         = "STATUS_NOTICE"
	OperationDeleteNotice         = "DELETE_NOTICE"
	OperationCreateWord           = "CREATE_WORD"
	OperationUpdateWord           = "UPDATE_WORD"
	OperationDeleteWord           = "DELETE_WORD"
	OperationCreateCategory       = "CREATE_CATEGORY"
	OperationUpdateCategory       = "UPDATE_CATEGORY"
	OperationStatusCategory       = "STATUS_CATEGORY"
	OperationDeleteCategory       = "DELETE_CATEGORY"
	OperationMarkReportProcessing = "MARK_REPORT_PROCESSING"
	OperationApproveReport        = "APPROVE_REPORT"
	OperationRejectReport         = "REJECT_REPORT"
	OperationMarkAppealProcessing = "MARK_APPEAL_PROCESSING"
	OperationApproveAppeal        = "APPROVE_APPEAL"
	OperationRejectAppeal         = "REJECT_APPEAL"
	OperationUserDisable          = "USER_DISABLE"
	OperationUserEnable           = "USER_ENABLE"
	OperationUserBan              = "USER_BAN"
	OperationUserUnban            = "USER_UNBAN"
	OperationUserRoleChange       = "USER_ROLE_CHANGE"
	OperationStudentVerifyApprove = "STUDENT_VERIFY_APPROVE"
	OperationStudentVerifyReject  = "STUDENT_VERIFY_REJECT"
	OperationUpdateProductStatus  = "UPDATE_PRODUCT_STATUS"
	OperationOrderExceptionClose  = "ORDER_EXCEPTION_CLOSE"
	OperationUpdateOrderStatus    = "UPDATE_ORDER_STATUS"
)

type AdminLog struct {
	ID            uint64    `json:"id"`
	AdminID       uint64    `json:"adminId"`
	OperationType string    `json:"operationType"`
	TargetType    string    `json:"targetType"`
	TargetID      uint64    `json:"targetId"`
	Description   *string   `json:"description,omitempty"`
	IPAddress     *string   `json:"ipAddress,omitempty"`
	RelatedType   *string   `json:"relatedType,omitempty"`
	RelatedID     *uint64   `json:"relatedId,omitempty"`
	CreateTime    time.Time `json:"createTime"`
	AdminName     string    `json:"adminName"`
	AdminAvatar   string    `json:"adminAvatar"`
	TargetName    string    `json:"targetName"`
	TargetCollege string    `json:"targetCollege"`
	TargetTitle   string    `json:"targetTitle"`
	ReportReason  string    `json:"reportReason"`
	ReporterName  string    `json:"reporterName"`
}

type LogActionInput struct {
	AdminID       uint64
	OperationType string
	TargetType    string
	TargetID      uint64
	Description   *string
	IPAddress     *string
	RelatedType   *string
	RelatedID     *uint64
}

type LogQuery struct {
	AdminID       uint64
	OperationType string
	TargetType    string
	TargetID      uint64
	StartTime     string
	EndTime       string
	Page          int
	PageSize      int
}

type Announcement struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	Content     *string    `json:"content,omitempty"`
	CoverURL    *string    `json:"coverUrl,omitempty"`
	Status      string     `json:"status"`
	PublishTime *time.Time `json:"publishTime,omitempty"`
	CreateBy    uint64     `json:"createBy"`
	CreateTime  time.Time  `json:"createTime"`
	UpdateTime  time.Time  `json:"updateTime"`
}

type AnnouncementQuery struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

type CreateAnnouncementInput struct {
	AdminID   uint64
	Title     string
	Content   *string
	CoverURL  *string
	Status    string
	IPAddress *string
}

type UpdateAnnouncementInput struct {
	AdminID   uint64
	ID        uint64
	Title     string
	Content   *string
	CoverURL  *string
	IPAddress *string
}

type UpdateAnnouncementStatusInput struct {
	AdminID   uint64
	ID        uint64
	Status    string
	IPAddress *string
}
