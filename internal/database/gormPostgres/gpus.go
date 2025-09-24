package gormpostgres

import (
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

func (c *GormPostgresController) AddGpu(gpu *inputs.AddGpuInput) (*models.Gpu, *models.Product, errors.PCCError) {
	panic("unimplemented")
}
