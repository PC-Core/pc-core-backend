package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
)

func (c *GormPostgresController) GetCpuChars(charId uint64) (*models.CpuChars, errors.PCCError) {
	var chars DbCpuChars

	err := c.db.Where("id = ?", charId).First(&chars).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return chars.IntoCpuChars(), nil
}

func (c *GormPostgresController) AddCpu(cpu *inputs.AddCpuInput) (*models.Product, *models.CpuChars, errors.PCCError) {
	tx := c.db.Begin()

	if tx.Error != nil {
		return nil, nil, errors.NewInternalSecretError()
	}

	defer tx.Rollback()

	chars := DbCpuChars{
		Name:         cpu.CpuName,
		PCores:       cpu.PCores,
		ECores:       cpu.ECores,
		Threads:      cpu.Threads,
		BasePFreqMHz: cpu.BasePFreqMHz,
		MaxPFreqMHz:  cpu.MaxPFreqMHz,
		BaseEFreqMHz: cpu.BaseEFreqMHz,
		MaxEFreqMHz:  cpu.MaxEFreqMHz,
		Socket:       cpu.Socket,
		L1KB:         cpu.L1KB,
		L2KB:         cpu.L2KB,
		L3KB:         cpu.L3KB,
		TecProcNM:    cpu.TecProcNM,
		TDPWatt:      cpu.TDPWatt,
		ReleaseYear:  cpu.ReleaseYear,
	}

	err := tx.Create(chars).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	medias, err := c.AddMedias(tx, cpu.Medias)

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	product := DbProduct{
		Name:           cpu.Name,
		Price:          cpu.Price,
		Selled:         0,
		Stock:          cpu.Stock,
		CharsTableName: database.CpuCharsTable,
		CharsID:        chars.ID,
	}

	err = tx.Create(&product).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	tx.Commit()

	return product.WithMediasIntoProduct(medias), chars.IntoCpuChars(), nil
}
