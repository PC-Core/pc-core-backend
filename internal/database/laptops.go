package database

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

func (c *DPostgresDbController) GetLaptopChars(charId uint64) (*models.LaptopChars, errors.PCCError) {
	var (
		id  uint64
		ram int16
		gpu string
		cpu models.CpuChars
	)

	query := fmt.Sprintf(`
	SELECT l.id, c.*, l.gpu, l.ram
	FROM %s AS l
	JOIN %s AS c ON l.cpu_id = c.id
	WHERE l.id = $1
	`, LaptopCharsTable, CpuCharsTable)

	row := c.db.QueryRow(query, charId)

	err := row.Scan(&id, &cpu.ID, &cpu.Name, &cpu.PCores, &cpu.ECores, &cpu.Threads, &cpu.BasePFreqMHz, &cpu.MaxPFreqMHz, &cpu.BaseEFreqMHz, &cpu.MaxEFreqMHz, &cpu.Socket, &cpu.L1KB, &cpu.L2KB, &cpu.L3KB, &cpu.TecProcNM, &cpu.TDPWatt, &cpu.ReleaseYear, &gpu, &ram)

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	return models.NewLaptopChars(id, &cpu, ram, gpu), nil
}

func (c *DPostgresDbController) AddLaptop(laptop *inputs.AddLaptopInput) (*models.Product, *models.LaptopChars, errors.PCCError) {
	var (
		charId    uint64
		productId uint64
	)

	tx, err := c.db.Begin()

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer tx.Rollback()

	err = tx.QueryRow(fmt.Sprintf("INSERT INTO %s (cpu_id, ram, gpu) VALUES ($1, $2, $3) returning id", LaptopCharsTable), laptop.CpuID, laptop.Ram, laptop.Gpu).Scan(&charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	productId, medias, err := c.AddProduct(tx, laptop.Name, laptop.Price, 0, laptop.Stock, laptop.Medias, LaptopCharsTable, charId)

	if err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	cpu, cerr := c.GetCpuChars(laptop.CpuID)

	if cerr != nil {
		return nil, nil, cerr
	}

	return models.NewProduct(productId, laptop.Name, laptop.Price, 0, laptop.Stock, medias, LaptopCharsTable, charId),
		models.NewLaptopChars(charId, cpu, laptop.Ram, laptop.Gpu),
		nil
}
