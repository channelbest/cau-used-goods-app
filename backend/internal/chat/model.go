package chat

const (
	ConversationStatusActive = "ACTIVE"
	ConversationStatusClosed = "CLOSED"
)

const (
	MessageTypeText       = "TEXT"
	MessageTypeOrderEvent = "ORDER_EVENT"
)

const (
	ActorTypeUser   = "USER"
	ActorTypeSystem = "SYSTEM"
)

const (
	EventTypeOrderCreated         = "ORDER_CREATED"
	EventTypeOrderConfirmed       = "ORDER_CONFIRMED"
	EventTypeOrderCanceled        = "ORDER_CANCELED"
	EventTypeOrderCompleted       = "ORDER_COMPLETED"
	EventTypeOrderTimeout         = "ORDER_TIMEOUT"
	EventTypeOrderExceptionClosed = "ORDER_EXCEPTION_CLOSED"
	EventTypeOrderStatusUpdated   = "ORDER_STATUS_UPDATED"
)

const (
	ReadStatusUnread = "UNREAD"
	ReadStatusRead   = "READ"
)

type Conversation struct {
	ID                 uint64  `json:"id"`
	ProductID          uint64  `json:"productId"`
	BuyerID            uint64  `json:"buyerId"`
	SellerID           uint64  `json:"sellerId"`
	LastMessageID      *uint64 `json:"lastMessageId,omitempty"`
	LastMessageContent *string `json:"lastMessageContent,omitempty"`
	LastMessageTime    *string `json:"lastMessageTime,omitempty"`
	BuyerUnreadCount   int     `json:"buyerUnreadCount"`
	SellerUnreadCount  int     `json:"sellerUnreadCount"`
	Status             string  `json:"status"`
	CreateTime         string  `json:"createTime"`
	UpdateTime         string  `json:"updateTime"`
}

type ConversationDetail struct {
	Conversation
	ProductTitle   string  `json:"productTitle"`
	ProductImage   *string `json:"productImage,omitempty"`
	TargetUserID   uint64  `json:"targetUserId"`
	TargetNickname *string `json:"targetNickname,omitempty"`
	UnreadCount    int     `json:"unreadCount"`
}

type Message struct {
	ID             uint64  `json:"id"`
	ConversationID uint64  `json:"conversationId"`
	SenderID       *uint64 `json:"senderId,omitempty"`
	ReceiverID     *uint64 `json:"receiverId,omitempty"`
	Content        string  `json:"content"`
	MessageType    string  `json:"messageType"`
	ActorType      string  `json:"actorType"`
	OrderID        *uint64 `json:"orderId,omitempty"`
	EventType      *string `json:"eventType,omitempty"`
	ReadStatus     string  `json:"readStatus"`
	CreateTime     string  `json:"createTime"`
}

type OrderEventInput struct {
	ProductID uint64
	BuyerID   uint64
	SellerID  uint64
	OrderID   uint64
	ActorID   *uint64
	ActorType string
	EventType string
	Content   string
}

type ProductForChat struct {
	ID       uint64
	SellerID uint64
	Title    string
	Status   string
}

type ConversationProduct struct {
	ID       uint64   `json:"id"`
	Title    string   `json:"title"`
	Price    float64  `json:"price"`
	Status   string   `json:"status"`
	Images   []string `json:"images"`
	SellerID uint64   `json:"sellerId"`
}
