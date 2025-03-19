package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/lib/pq"
)

type TempProduct struct {
	models.Product
	JSONMedias string
}

func (c *DPostgresDbController) GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError) {
	query := `
	SELECT p.*, 
       COALESCE(json_agg(json_build_object('id', m.id, 'url', m.url, 'type', m.type))
	   		FILTER (WHERE m.id IS NOT NULL), '[]') AS medias
	FROM Products p
	LEFT JOIN Medias m ON m.id = ANY(p.medias)
	GROUP BY p.id
	ORDER BY p.id
	LIMIT $1 OFFSET $2;
	`

	rows, err := c.db.Query(query, count, start)

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
		ids           pq.Int64Array
		jsonMedia     string
		charTableName string
		charId        uint64
	)

	if err := rows.Scan(&rid, &name, &price, &selled, &stock, &charTableName, &charId, &ids, &jsonMedia); err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	medias, err := c.MediasFromJson(jsonMedia)

	if err != nil {
		return nil, err
	}

	return models.NewProduct(rid, name, price, selled, stock, medias, charTableName, charId), nil
}

func (c *DPostgresDbController) MediasFromJson(jss string) ([]models.Media, errors.PCCError) {
	var medias []models.Media
	err := json.Unmarshal([]byte(jss), &medias)

	if err != nil {
		return nil, errors.NewJsonUnmarshalError()
	}

	return medias, nil
}

func (c *DPostgresDbController) GetProductById(id uint64) (*models.Product, errors.PCCError) {
	query := `
	SELECT p.*, 
       COALESCE(json_agg(json_build_object('id', m.id, 'url', m.url, 'type', m.type))
	   		FILTER (WHERE m.id IS NOT NULL), '[]') AS medias
	FROM Products p
	LEFT JOIN Medias m ON m.id = ANY(p.medias)
	WHERE p.id = $1
	GROUP BY p.id;
	`

	row := c.db.QueryRow(query, id)

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
	case CpuCharsTable:
		return c.GetCpuChars(p.CharId)
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
		p, err := c.ScanProduct(rows)

		if err != nil {
			return nil, err
		}

		quantity, exists := quantityMap[p.ID]
		if !exists {
			return nil, errors.NewInternalSecretError()
		}
		cartItem := models.NewCartItem(*p, quantity, time.Now())

		cartItems = append(cartItems, *cartItem)
	}

	return cartItems, nil
}

func (c *DPostgresDbController) AddProduct(tx *sql.Tx, name string, price float64, selled uint64, stock uint64, medias []models.InputMedia, charTable string, charID uint64) (uint64, []models.Media, errors.PCCError) {
	meds, aerr := c.AddMedias(medias)

	if aerr != nil {
		return 0, nil, aerr
	}

	var productId uint64
	err := tx.QueryRow("INSERT INTO Products (name, price, selled, stock, medias, chars_table_name, chars_id) VALUES ($1, $2, $3, $4, $5, $6, $7) returning id", name, price, selled, stock, pq.Array(c.IDsFromMedias(meds)), charTable, charID).Scan(&productId)

	if err != nil {
		return 0, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return productId, meds, nil
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
