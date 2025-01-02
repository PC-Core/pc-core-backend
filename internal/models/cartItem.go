package models

import "time"

type CartItem struct {
	ID       uint64    `json:"id"`
	Product  Product   `json:"product"`
	Quantity uint      `json:"quantity"`
	AddedAt  time.Time `json:"added_at"`
}
