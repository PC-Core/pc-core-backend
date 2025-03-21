package models

type Cart struct {
	UserID uint64     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

func NewCart(user_id uint64, items []CartItem) *Cart {
	return &Cart{
		user_id, items,
	}
}
