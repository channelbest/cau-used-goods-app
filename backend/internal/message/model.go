package message

import "time"

const (
	ReadStatusUnread = "UNREAD"
	ReadStatusRead   = "READ"
)

const (
	MessageTypeOrderCreated   = "ORDER_CREATED"
	MessageTypeOrderConfirmed = "ORDER_CONFIRMED"
	MessageTypeOrderCanceled  = "ORDER_CANCELED"
	MessageTypeOrderTimeout   = "ORDER_TIMEOUT"
	MessageTypeReportHandled  = "REPORT_HANDLED"
	MessageTypeSystemNotice   = "SYSTEM_NOTICE"
)

const (
	RelatedTypeOrder   = "ORDER"
	RelatedTypeProduct = "PRODUCT"
	RelatedTypeReport  = "REPORT"
	RelatedTypeNotice  = "NOTICE"
)

type Message struct {
	ID          uint64     `json:"id"`
	ReceiverID  uint64     `json:"receiverId"`
	SenderID    *uint64    `json:"senderId,omitempty"`
	MessageType string     `json:"messageType"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	RelatedType *string    `json:"relatedType,omitempty"`
	RelatedID   *uint64    `json:"relatedId,omitempty"`
	ReadStatus  string     `json:"readStatus"`
	CreateTime  time.Time  `json:"createTime"`
	ReadTime    *time.Time `json:"readTime,omitempty"`
}

type CreateMessageInput struct {
	ReceiverID  uint64
	SenderID    *uint64
	MessageType string
	Title       string
	Content     string
	RelatedType *string
	RelatedID   *uint64
}
