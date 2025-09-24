package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPostgresController struct {
	db *gorm.DB
}

// AddGpu implements database.DbController.
func (c *GormPostgresController) AddGpu(gpu *inputs.AddGpuInput) (*models.Gpu, *models.Product, errors.PCCError) {
	panic("unimplemented")
}

func NewGormPostgresController(conn string) (*GormPostgresController, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &GormPostgresController{db}, nil
}
