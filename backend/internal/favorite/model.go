package favorite

type Favorite struct {
	ID         uint64  `json:"id"`
	UserID     uint64  `json:"userId"`
	ProductID  uint64  `json:"productId"`
	CreateTime string  `json:"createTime"`
	IsDeleted  bool    `json:"isDeleted"`
}

type FavoriteDetail struct {
	Favorite
	ProductTitle   string  `json:"productTitle"`
	ProductPrice   float64 `json:"productPrice"`
	ProductStatus  string  `json:"productStatus"`
	ProductImage   *string `json:"productImage,omitempty"`
	SellerNickname *string `json:"sellerNickname,omitempty"`
}
