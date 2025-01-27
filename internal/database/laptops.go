package database

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/models"
)

func (c *DPostgresDbController) GetLaptopChars(charId uint64) (*models.LaptopChars, error) {
	var (
		id  uint64
		cpu string
		ram int16
		gpu string
	)

	row := c.db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", LaptopCharsTable), charId)

	err := row.Scan(&id, &cpu, &ram, &gpu)

	if err != nil {
		return nil, err
	}

	return models.NewLaptopChars(id, cpu, ram, gpu), nil
}

func (c *DPostgresDbController) AddLaptop(name string, price float64, selled uint64, stock uint64, cpu string, ram int16, gpu string) (*models.Product, *models.LaptopChars, error) {
	var (
		charId    uint64
		productId uint64
	)

	tx, err := c.db.Begin()

	if err != nil {
		return nil, nil, err
	}

	defer tx.Rollback()

	err = tx.QueryRow(fmt.Sprintf("INSERT INTO %s (cpu, ram, gpu) VALUES ($1, $2, $3) returning id", LaptopCharsTable), cpu, ram, gpu).Scan(&charId)

	if err != nil {
		return nil, nil, err
	}

	err = tx.QueryRow("INSERT INTO Products (name, price, selled, stock, chars_table_name, chars_id) VALUES ($1, $2, $3, $4, $5, $6) returning id", name, price, selled, stock, LaptopCharsTable, charId).Scan(&productId)

	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("error on commiting operation: %s", err)
	}

	return models.NewProduct(productId, name, price, selled, stock, LaptopCharsTable, charId),
		models.NewLaptopChars(charId, cpu, ram, gpu),
		nil
}
