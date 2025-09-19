package outputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type ProductWithChars = struct {
	*models.Product `json:"product"`
	Chars           *RestCharsObject `json:"chars"`
}

func NewProductWithChars(product *models.Product, chars *RestCharsObject) *ProductWithChars {
	return &ProductWithChars{
		product, chars,
	}
}
