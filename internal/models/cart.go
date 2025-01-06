package models

type Cart struct {
	ID     uint64     `json:"id"`
	UserID uint64     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

func NewCart(id uint64, user_id uint64, items []CartItem) *Cart {
	return &Cart{
		id, user_id, items,
	}
}
