package database

import (
	"github.com/Core-Mouse/cm-backend/internal/models"
)

func (c *DbController) GetCartByUserID(userID uint64) (*models.Cart, error) {
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
		return nil, err
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
			return nil, err
		}

		item.Product = product
		cart.Items = append(cart.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *DbController) SingleExecOrError(req string, params ...any) error {
	_, err := c.db.Exec(req, params)
	return err
}

func (c *DbController) AddToCart(product_id, user_id, quantity uint64) error {
	return c.SingleExecOrError("INSERT INTO Cart (user_id, product_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (user_id, product_id) DO UPDATE SET quantity = Cart.quantity + $3;", user_id, product_id, quantity)
}

func (c *DbController) RemoveFromCart(product_id, user_id uint64) error {
	return c.SingleExecOrError("DELETE FROM Cart WHERE user_id = $1 AND product_id = $2", user_id, product_id)
}

func (c *DbController) ChangeQuantity(product_id, user_id uint64, val int64) error {
	return c.SingleExecOrError("UPDATE Cart SET quantity = GREATEST(quantity + $1, 1) WHERE product_id = $2 AND user_id = $3", val, product_id, user_id)
}

// func (c *DbController) parseCart(itemRows *sql.Rows) (*models.Cart, error) {
// 	var items []models.CartItem

// 	for itemRows.Next() {
// 		var ()
// 		items = append(items)
// 	}
// }

// func (c *DbController) loadCartProducts()
