package gormpostgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPostgresController struct {
	db *gorm.DB
}

func NewGormPostgresController(conn string) (*GormPostgresController, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return &GormPostgresController{db}, nil
}
