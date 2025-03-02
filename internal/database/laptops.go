package database

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/database/dberrors"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

func (c *DPostgresDbController) GetLaptopChars(charId uint64) (*models.LaptopChars, errors.PCCError) {
	var (
		id  uint64
		cpu string
		ram int16
		gpu string
	)

	row := c.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", LaptopCharsTable), charId)

	err := row.Scan(&id, &cpu, &ram, &gpu)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewLaptopChars(id, cpu, ram, gpu), nil
}

func (c *DPostgresDbController) AddLaptop(name string, price float64, selled uint64, stock uint64, cpu string, ram int16, gpu string, imedias []models.InputMedia) (*models.Product, *models.LaptopChars, errors.PCCError) {
	var (
		charId    uint64
		productId uint64
	)

	tx, err := c.db.Begin()

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer tx.Rollback()

	err = tx.QueryRow(fmt.Sprintf("INSERT INTO %s (cpu, ram, gpu) VALUES ($1, $2, $3) returning id", LaptopCharsTable), cpu, ram, gpu).Scan(&charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	productId, medias, err := c.AddProduct(tx, name, price, selled, stock, imedias, LaptopCharsTable, charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewProduct(productId, name, price, selled, stock, medias, LaptopCharsTable, charId),
		models.NewLaptopChars(charId, cpu, ram, gpu),
		nil
}
