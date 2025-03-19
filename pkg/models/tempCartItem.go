package models

type TempCartItem struct {
	ProductID uint64
	Quantity  uint
}

func NewTempCartItem(product_id uint64, quantity uint) *TempCartItem {
	return &TempCartItem{
		product_id, quantity,
	}
}
