package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

func (c *GormPostgresController) GetProducts(start uint64, count uint64) ([]models.Product, errors.PCCError) {
	var dbproducts []DbProductWithMedias

	err := c.db.
		Select("p.*, COALESCE(json_agg(json_build_object('id', m.id, 'url', m.url, 'type', m.type)), '[]') AS medias").
		Table("products AS p").
		Joins("LEFT JOIN medias m ON m.id = ANY(p.medias)").
		Group("p.id").
		Order("p.id").
		Limit(int(count)).
		Offset(int(start)).
		Scan(&dbproducts).Error

	if err != nil {
		return nil, errors.NewInternalSecretError()
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
		Select("p.*, COALESCE(json_agg(json_build_object('id', m.id, 'url', m.url, 'type', m.type)), '[]') AS medias").
		Table("products AS p").
		Joins("LEFT JOIN medias m ON m.id = ANY(p.medias)").
		Where("p.id = ?", id).
		Group("p.id").
		Scan(&dbproduct).Error

	if err != nil {
		return nil, errors.NewInternalSecretError()
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
		return nil, errors.NewInternalSecretError()
	}
}
