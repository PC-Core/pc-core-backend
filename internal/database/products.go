package database

import (
	"fmt"
	"time"

	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/lib/pq"
)

type RowLike interface {
	Scan(...any) error
}

func (c *DbController) GetProducts(start uint64, count uint64) ([]models.Product, error) {
	rows, err := c.db.Query("SELECT * FROM Products OFFSET $1 LIMIT $2", start, count)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	products := make([]models.Product, 0, count)

	for rows.Next() {

		p, err := c.ScanProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(
			products,
			*p,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (c *DbController) ScanProduct(rows RowLike) (*models.Product, error) {
	var (
		rid           uint64
		name          string
		price         float64
		selled        uint64
		stock         uint64
		charTableName string
		charId        uint64
	)

	if err := rows.Scan(&rid, &name, &price, &selled, &stock, &charTableName, &charId); err != nil {
		return nil, err
	}

	return models.NewProduct(rid, name, price, selled, stock, charTableName, charId), nil
}

func (c *DbController) GetProductById(id uint64) (*models.Product, error) {
	row := c.db.QueryRow("SELECT * FROM Products WHERE id = $1", id)

	p, err := c.ScanProduct(row)

	if err != nil {
		return nil, err
	}

	return p, nil

}

func (c *DbController) GetProductCharsByProductID(productId uint64) (ProductChars, error) {
	p, err := c.GetProductById(productId)

	if err != nil {
		return nil, err
	}

	switch p.CharTableName {
	case LaptopCharsTable:
		return c.GetLaptopChars(p.CharId)
	default:
		return nil, fmt.Errorf("unknown chars table name: %s", p.CharTableName)
	}
}

func (c *DbController) LoadProductsRangeAsCartItem(rng []uint64) ([]models.CartItem, error) {
	cartItems := make([]models.CartItem, 0)

	query := `
		SELECT id, name, price, selled, stock, chars_table_name, chars_id, COUNT(*) AS quantity
		FROM Products
		WHERE id = ANY($1::bigint[])
		GROUP BY id, name, price, selled, stock, chars_table_name, chars_id
	`

	rows, err := c.db.Query(query, pq.Array(rng))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            uint64
			name          string
			price         float64
			selled        uint64
			stock         uint64
			chartablename string
			charid        uint64
			quantity      uint
		)

		err := rows.Scan(&id, &name, &price, &selled, &stock, &chartablename, &charid, &quantity)
		if err != nil {
			return nil, err
		}

		product := models.NewProduct(id, name, price, selled, stock, chartablename, charid)
		cartItem := models.NewCartItem(0, *product, quantity, time.Now())

		cartItems = append(cartItems, *cartItem)
	}

	return cartItems, nil
}

// func (c *DbController) LoadProductsRangeAsCartItem(rng []uint64) ([]models.CartItem, error) {
// 	products := make([]models.CartItem, 0)

// 	rows, err := c.db.Query("SELECT * FROM Products WHERE id = ANY($1::integer[])", pq.Array(rng))

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var (
// 			id            uint64
// 			name          string
// 			price         float64
// 			selled        uint64
// 			stock         uint64
// 			chartablename string
// 			charid        uint64
// 		)

// 		err := rows.Scan(&id, &name, &price, &selled, &stock, &chartablename, &charid)

// 		if err != nil {
// 			return nil, err
// 		}

// 		products = append(products, *models.NewCartItem(0, *models.NewProduct(id, name, price, selled, stock, chartablename, charid), 0, time.Now()))
// 	}

// 	return products, nil
// }
