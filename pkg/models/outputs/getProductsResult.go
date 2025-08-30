package outputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type GetProductsResult struct {
	Products []models.Product `json:"products"`
	Amount   uint64           `json:"amount"`
	Page     uint64           `json:"page"`
}

func NewGetProductsResult(products []models.Product, amount uint64, page uint64) *GetProductsResult {
	return &GetProductsResult{
		products,
		amount,
		page,
	}
}
