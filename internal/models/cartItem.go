package models

import "time"

type CartItem struct {
	ID       uint64    `json:"id"`
	Product  Product   `json:"product"`
	Quantity uint      `json:"quantity"`
	AddedAt  time.Time `json:"added_at"`
}

func NewCartItem(id uint64, product Product, quantity uint, AddedAt time.Time) *CartItem {
	return &CartItem{
		id, product, quantity, AddedAt,
	}
}
