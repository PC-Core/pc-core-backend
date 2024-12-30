package models

type Product struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Selled        uint64  `json:"selled"`
	Stock         uint64  `json:"stock"`
	CharTableName string  `json:"-"`
	CharId        uint64  `json:"-"`
}

func NewProduct(id uint64, name string, price float64, selled uint64, stock uint64, charTableName string, charId uint64) *Product {
	return &Product{
		id, name, price, selled, stock, charTableName, charId,
	}
}
