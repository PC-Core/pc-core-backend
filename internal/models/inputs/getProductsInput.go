package inputs

type GetProductsInput struct {
	Page  uint64 `json:"page" form:"page"`
	Count uint64 `json:"count" form:"count"`
}
