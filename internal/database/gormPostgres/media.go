package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"gorm.io/gorm"
)

func (c *GormPostgresController) AddMedias(tx *gorm.DB, imedias []models.InputMedia) (models.Medias, errors.PCCError) {
	var medias []models.Media

	for _, im := range imedias {
		medias = append(medias, models.Media{Url: im.Url, Type: im.Type})
	}

	if err := tx.Create(&medias).Error; err != nil {
		return nil, errors.NewInternalSecretError()
	}

	return medias, nil
}
