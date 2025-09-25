package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"gorm.io/gorm"
)

func (c *GormPostgresController) GetGpuChars() ([]models.Gpu, errors.PCCError) {
	var gpus []models.Gpu

	err := c.db.Model(&models.Gpu{}).Find(&gpus).Error
	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return gpus, nil
}

func (c *GormPostgresController) GetGpuByID(id uint64) (*models.Gpu, errors.PCCError) {
	var gpu models.Gpu

	err := c.db.Model(&models.Gpu{}).Where("id = ?", id).First(&gpu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, gormerrors.GormErrorCast(err)
	}

	return &gpu, nil
}

func (c *GormPostgresController) AddGpu(gpu *inputs.AddGpuInput) (*models.GpuChars, *models.Product, errors.PCCError) {
	tx := c.db.Begin()

	if tx.Error != nil {
		return nil, nil, errors.NewInternalSecretError()
	}

	defer tx.Rollback()

	chars := DbGpuChars{
		ID:           uint64(gpu.ID),
		Name:         gpu.Name,
		MemoryGB:     gpu.MemoryGB,
		MemoryType:   gpu.MemoryType,
		BusWidthBit:  gpu.BusWidthBit,
		BaseFreqMHz:  gpu.BaseFreqMHz,
		BoostFreqMHz: gpu.BoostFreqMHz,
		TecprocNm:    gpu.TecprocNm,
		TDPWatt:      gpu.TDPWatt,
		ReleaseYear:  gpu.RealeseYear,
	}

	err := tx.Create(chars).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	medias, err := c.AddMedias(tx, gpu.Medias)

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	product := DbProduct{
		Name:           gpu.Name,
		Price:          gpu.Price,
		Selled:         0,
		Stock:          gpu.Stock,
		CharsTableName: database.GpuCharsTable,
		CharsID:        chars.ID,
	}

	err = tx.Create(&product).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	tx.Commit()

	return chars.IntoGpu(), product.WithMediasIntoProduct(medias), nil
}
