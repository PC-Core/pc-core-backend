package gormpostgres

import (
	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"gorm.io/gorm"
)

func (c *GormPostgresController) AddMedias(tx *gorm.DB, imedias []models.InputMedia) (models.Medias, errors.PCCError) {
	var medias []models.Media

	if (len(imedias) == 0) || (tx == nil) {
		return medias, nil
	}

	for _, im := range imedias {
		medias = append(medias, models.Media{Url: im.Url, Type: im.Type})
	}

	if err := tx.Create(&medias).Error; err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return medias, nil
}
