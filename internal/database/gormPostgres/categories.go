package gormpostgres

import (
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
)

func (c *GormPostgresController) GetCategories() ([]models.Category, errors.PCCError) {
	var dbcats []DbCategories

	err := c.db.Find(&dbcats).Error
	if err != nil {
		// TODO: error type
		return nil, errors.NewInternalSecretError()
	}

	cats := make([]models.Category, len(dbcats))

	for i, cat := range dbcats {
		cats[i] = models.Category{
			ID:          cat.ID,
			Title:       cat.Title,
			Description: cat.Description,
			Icon:        cat.Icon,
			Slug:        cat.Slug,
		}
	}

	return cats, nil
}
