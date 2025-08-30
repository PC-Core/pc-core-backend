package gormpostgres

import (
	"time"

	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

func (c *GormPostgresController) GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError) {
	var dbproducts []DbProductWithMedias

	err := c.db.
		Preload("Medias").
		Order("id").
		Limit(int(count)).
		Offset(int(start)).
		Find(&dbproducts).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	products := make([]models.Product, 0, len(dbproducts))

	for _, product := range dbproducts {
		products = append(products, *product.IntoProduct())
	}

	return products, nil
}

func (c *GormPostgresController) GetProductById(id uint64) (*models.Product, errors.PCCError) {
	var dbproduct DbProductWithMedias

	err := c.db.
		Preload("Medias").
		Where("id = ?", id).
		First(&dbproduct).
		Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return dbproduct.IntoProduct(), nil
}

func (c *GormPostgresController) GetProductCharsByProductID(productId uint64) (database.ProductChars, errors.PCCError) {
	p, err := c.GetProductById(productId)

	if err != nil {
		return nil, err
	}

	switch p.CharTableName {
	case database.LaptopCharsTable:
		return c.GetLaptopChars(p.CharId)
	case database.CpuCharsTable:
		return c.GetCpuChars(p.CharId)
	default:
		return nil, gormerrors.GormErrorCast(err)
	}
}

func (c *GormPostgresController) LoadProductsRangeAsCartItem(tempCart []models.TempCartItem) ([]models.CartItem, errors.PCCError) {
	productIDs := make([]uint64, len(tempCart))
	quantityMap := make(map[uint64]uint)

	for i, item := range tempCart {
		productIDs[i] = item.ProductID
		quantityMap[item.ProductID] = item.Quantity
	}

	var products []DbProductWithMedias

	err := c.db.
		Preload("Medias").
		Where("id IN ?", productIDs).
		Find(&products).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	cartItems := make([]models.CartItem, 0, len(products))
	for _, p := range products {
		quantity, exists := quantityMap[p.ID]
		if !exists {
			return nil, gormerrors.GormErrorCast(err)
		}
		cartItem := models.NewCartItem(*p.IntoProduct(), quantity, time.Now())
		cartItems = append(cartItems, *cartItem)
	}

	return cartItems, nil
}
