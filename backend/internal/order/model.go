package order

type Order struct {
	ID                   uint64  `json:"id"`
	OrderNo              string  `json:"orderNo"`
	ProductID            uint64  `json:"productId"`
	BuyerID              uint64  `json:"buyerId"`
	SellerID             uint64  `json:"sellerId"`
	ProductTitleSnapshot string  `json:"productTitleSnapshot"`
	ProductPriceSnapshot float64 `json:"productPriceSnapshot"`
	Status               string  `json:"status"`
	Remark               *string `json:"remark,omitempty"`
	MeetTime             *string `json:"meetTime,omitempty"`
	MeetLocation         *string `json:"meetLocation,omitempty"`
	CancelReason         *string `json:"cancelReason,omitempty"`
	CancelBy             *uint64 `json:"cancelBy,omitempty"`
	ExpireTime           string  `json:"expireTime"`
	ConfirmTime          *string `json:"confirmTime,omitempty"`
	FinishTime           *string `json:"finishTime,omitempty"`
	CloseTime            *string `json:"closeTime,omitempty"`
	CreateTime           string  `json:"createTime"`
	UpdateTime           string  `json:"updateTime"`
}

type OrderDetail struct {
	Order
	BuyerNickname  *string `json:"buyerNickname,omitempty"`
	SellerNickname *string `json:"sellerNickname,omitempty"`
	ProductImage   *string `json:"productImage,omitempty"`
}
