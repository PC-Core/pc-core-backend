package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (c *GormPostgresController) GetCartByUserID(userID uint64) (*models.Cart, errors.PCCError) {
	var cart []DbCart

	err := c.db.
		Joins("JOIN products ON products.id = cart.product_id").
		Where("cart.user_id = ?", userID).
		Find(&cart)

	if err != nil {
		// TODO: type
		return nil, errors.NewInternalSecretError()
	}

	return DbCartIntoCart(cart), nil
}

func (c *GormPostgresController) AddToCart(product_id, user_id, quantity uint64) (uint64, errors.PCCError) {
	cartItem := DbCart{
		UserID:    user_id,
		ProductID: product_id,
		Quantity:  uint(quantity),
	}

	err := c.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"quantity": gorm.Expr("quantity + ?", quantity),
			}),
		}).
		Create(&cartItem).Error

	if err != nil {
		// TODO: error type
		return product_id, errors.NewInternalSecretError()
	}

	return product_id, nil
}

func (c *GormPostgresController) RemoveFromCart(productID, userID uint64) (uint64, errors.PCCError) {
	err := c.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&models.CartItem{}).Error

	if err != nil {
		return productID, errors.NewInternalSecretError()
	}

	// TODO: type
	return productID, nil
}

func (c *GormPostgresController) ChangeQuantity(productID, userID uint64, val int64) (uint64, errors.PCCError) {
	err := c.db.Model(&models.CartItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", gorm.Expr("GREATEST(quantity + ?, 1)", val)).Error

	if err != nil {
		// TODO: error type
		return productID, errors.NewInternalSecretError()
	}

	return productID, nil
}
