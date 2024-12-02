package models

import "math/big"

type Laptop struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Cpu      string    `json:"cpu"`
	Ram      int16     `json:"ram"`
	Gpu      string    `json:"gpu"`
	Price    big.Float `json:"price"`
	Discount int16     `json:"discount"`
}

func NewLaptop(id int, name string, cpu string, ram int16, gpu string, price big.Float, discount int16) *Laptop {
	return &Laptop{
		id, name, cpu, ram, gpu, price, discount,
	}
}
