package database

import (
	"time"

	"github.com/Core-Mouse/cm-backend/internal/database/dberrors"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/lib/pq"
)

func (c *DPostgresDbController) GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError) {
	rows, err := c.db.Query("SELECT * FROM Products OFFSET $1 LIMIT $2", start, count)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer rows.Close()
	products := make([]models.Product, 0, count)

	for rows.Next() {

		p, err := c.ScanProduct(rows)

		if err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		products = append(
			products,
			*p,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return products, nil
}

func (c *DPostgresDbController) ScanProduct(rows RowLike) (*models.Product, errors.PCCError) {
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
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewProduct(rid, name, price, selled, stock, charTableName, charId), nil
}

func (c *DPostgresDbController) GetProductById(id uint64) (*models.Product, errors.PCCError) {
	row := c.db.QueryRow("SELECT * FROM Products WHERE id = $1", id)

	p, err := c.ScanProduct(row)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return p, nil

}

func (c *DPostgresDbController) GetProductCharsByProductID(productId uint64) (ProductChars, errors.PCCError) {
	p, err := c.GetProductById(productId)

	if err != nil {
		return nil, err
	}

	switch p.CharTableName {
	case LaptopCharsTable:
		return c.GetLaptopChars(p.CharId)
	default:
		return nil, errors.NewInternalSecretError()
	}
}

func (c *DPostgresDbController) LoadProductsRangeAsCartItem(tempCart []models.TempCartItem) ([]models.CartItem, errors.PCCError) {
	productIDs := make([]uint64, len(tempCart))
	quantityMap := make(map[uint64]uint)
	for i, item := range tempCart {
		productIDs[i] = item.ProductID
		quantityMap[item.ProductID] = item.Quantity
	}

	query := `
		SELECT id, name, price, selled, stock, chars_table_name, chars_id
		FROM Products
		WHERE id = ANY($1::bigint[])
	`
	rows, err := c.db.Query(query, pq.Array(productIDs))
	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}
	defer rows.Close()

	cartItems := make([]models.CartItem, 0)
	for rows.Next() {
		var (
			id            uint64
			name          string
			price         float64
			selled        uint64
			stock         uint64
			chartablename string
			charid        uint64
		)

		err := rows.Scan(&id, &name, &price, &selled, &stock, &chartablename, &charid)
		if err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		quantity, exists := quantityMap[id]
		if !exists {
			return nil, errors.NewInternalSecretError()
		}

		product := models.NewProduct(id, name, price, selled, stock, chartablename, charid)
		cartItem := models.NewCartItem(*product, quantity, time.Now())

		cartItems = append(cartItems, *cartItem)
	}

	return cartItems, nil
}

// func (c *DbController) LoadProductsRangeAsCartItem(rng []models.TempCartItem) ([]models.CartItem, error) {
// 	cartItems := make([]models.CartItem, 0)

// 	query := `
// 		SELECT * FROM Products
// 		WHERE id = ANY($1::bigint[])
// 		GROUP BY id, name, price, selled, stock, chars_table_name, chars_id
// 	`

// 	rows, err := c.db.Query(query, pq.Array(rng))
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

// 		product := models.NewProduct(id, name, price, selled, stock, chartablename, charid)
// 		cartItem := models.NewCartItem(0, *product, quantity, time.Now())

// 		cartItems = append(cartItems, *cartItem)
// 	}

// 	return cartItems, nil
// }

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
