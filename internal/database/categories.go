package database

import (
	"github.com/Core-Mouse/cm-backend/internal/database/dberrors"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

func (c *DPostgresDbController) GetCategories() ([]models.Category, errors.PCCError) {
	cats := make([]models.Category, 0, 5)

	res, err := c.db.Query("SELECT * FROM Categories")

	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer res.Close()

	for res.Next() {
		var cat models.Category

		if err := res.Scan(&cat.ID, &cat.Title, &cat.Description, &cat.Icon, &cat.Slug); err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}

		cats = append(cats, cat)
	}

	return cats, nil
}
