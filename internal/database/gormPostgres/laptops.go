package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

func (c *GormPostgresController) GetLaptopChars(charId uint64) (*models.LaptopChars, errors.PCCError) {
	var chars DbLaptopChars

	err := c.db.
		Preload("Cpu").
		Where("id = ?", charId).
		First(&chars).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return chars.IntoLaptopChars(), nil
}

func (c *GormPostgresController) AddLaptop(laptop *inputs.AddLaptopInput) (*models.Product, *models.LaptopChars, errors.PCCError) {
	tx := c.db.Begin()

	if tx.Error != nil {
		return nil, nil, gormerrors.GormErrorCast(tx.Error)
	}

	defer tx.Rollback()

	chars := DbLaptopChars{
		CpuID: laptop.CpuID,
		Ram:   laptop.Ram,
		Gpu:   laptop.Gpu,
	}

	err := tx.
		Create(&chars).
		Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	medias, err := c.AddMedias(tx, laptop.Medias)

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	product := DbProduct{
		Name:           laptop.Name,
		Price:          laptop.Price,
		Selled:         0,
		Stock:          laptop.Stock,
		CharsTableName: database.LaptopCharsTable,
		CharsID:        chars.ID,
	}

	err = tx.Create(&product).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	tx.Commit()

	return product.WithMediasIntoProduct(medias), chars.IntoLaptopChars(), nil
}
