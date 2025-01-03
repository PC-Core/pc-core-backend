package models

type Cart struct {
	ID     uint64     `json:"id"`
	UserID uint64     `json:"user_id"`
	Items  []CartItem `json:"items"`
}
