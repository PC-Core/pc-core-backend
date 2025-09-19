package models

import "time"

type CartItem struct {
	Product  Product   `json:"product"`
	Quantity uint      `json:"quantity"`
	AddedAt  time.Time `json:"added_at"`
}

func NewCartItem(product Product, quantity uint, AddedAt time.Time) *CartItem {
	return &CartItem{
		product, quantity, AddedAt,
	}
}
