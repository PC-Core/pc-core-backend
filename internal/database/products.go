package database

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/models"
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
