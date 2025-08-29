package inputs

type AddToCartInput struct {
	ProductID uint64 `json:"product_id"`
	Quantity  int64  `json:"quantity"`
}
