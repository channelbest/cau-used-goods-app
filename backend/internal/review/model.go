package review

type Review struct {
	ID         uint64  `json:"id"`
	OrderID    uint64  `json:"orderId"`
	ProductID  uint64  `json:"productId"`
	ReviewerID uint64  `json:"reviewerId"`
	SellerID   uint64  `json:"sellerId"`
	Rating     int     `json:"rating"`
	Content    *string `json:"content,omitempty"`
	Status     string  `json:"status"`
	CreateTime string  `json:"createTime"`
	IsDeleted  bool    `json:"isDeleted"`
}

type ReviewDetail struct {
	Review
	ReviewerNickname *string `json:"reviewerNickname,omitempty"`
	ProductTitle     string  `json:"productTitle"`
	ProductImage     *string `json:"productImage,omitempty"`
}
