package database

import (
	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

func (c *DPostgresDbController) GetCartByUserID(userID uint64) (*models.Cart, errors.PCCError) {
	query := `
	SELECT
		Cart.quantity,
		Cart.added_at,
		Products.id AS product_id,
		Products.name,
		Products.price,
		Products.selled,
		Products.stock,
		Products.chars_table_name,
		Products.chars_id
	FROM
		Cart
	JOIN
		Products
	ON
		Cart.product_id = Products.id
	WHERE
		Cart.user_id = $1;`

	rows, err := c.db.Query(query, userID)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer rows.Close()

	cart := &models.Cart{
		UserID: userID,
		Items:  []models.CartItem{},
	}

	for rows.Next() {
		var item models.CartItem
		var product models.Product

		err := rows.Scan(&item.Quantity, &item.AddedAt,
			&product.ID, &product.Name, &product.Price, &product.Selled,
			&product.Stock, &product.CharTableName, &product.CharId)

		if err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		item.Product = product
		cart.Items = append(cart.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return cart, nil
}

// func (c *DbController) SingleExecOrError(req string, params ...any) error {

// }

func (c *DPostgresDbController) AddToCart(product_id, user_id, quantity uint64) (uint64, errors.PCCError) {
	_, err := c.db.Exec("INSERT INTO Cart (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = Cart.quantity + $3;", user_id, product_id, quantity)
	return product_id, dberrors.PQDbErrorCaster(c.db, err)
}

func (c *DPostgresDbController) RemoveFromCart(product_id, user_id uint64) (uint64, errors.PCCError) {
	_, err := c.db.Exec("DELETE FROM Cart WHERE user_id = $1 AND product_id = $2", user_id, product_id)
	return product_id, dberrors.PQDbErrorCaster(c.db, err)
}

func (c *DPostgresDbController) ChangeQuantity(product_id, user_id uint64, val int64) (uint64, errors.PCCError) {
	_, err := c.db.Exec("UPDATE Cart SET quantity = GREATEST(quantity + $1, 1) WHERE product_id = $2 AND user_id = $3", val, product_id, user_id)
	return product_id, dberrors.PQDbErrorCaster(c.db, err)
}

// func (c *DbController) parseCart(itemRows *sql.Rows) (*models.Cart, error) {
// 	var items []models.CartItem

// 	for itemRows.Next() {
// 		var ()
// 		items = append(items)
// 	}
// }

// func (c *DbController) loadCartProducts()
