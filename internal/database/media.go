package database

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/database/dberrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
)

func (c *DPostgresDbController) AddMedias(imedias []models.InputMedia) ([]models.Media, errors.PCCError) {
	var values []interface{}

	query := "INSERT INTO Medias (url, type) VALUES "

	for i, media := range imedias {
		query += fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2)
		if i < len(imedias)-1 {
			query += ", "
		}
		values = append(values, media.Url, media.Type)
	}

	query += " RETURNING id;"

	rows, err := c.db.Query(query, values...)
	if err != nil {
		return nil, dberrors.PQDbErrorCaster(c.db, err)
	}

	defer rows.Close()

	var medias []models.Media
	i := 0

	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, dberrors.PQDbErrorCaster(c.db, err)
		}
		medias = append(medias, *models.NewMediaFromInput(id, &imedias[i]))
		i += 1
	}

	return medias, nil
}

func (c *DPostgresDbController) IDsFromMedias(medias []models.Media) []uint64 {
	var ids []uint64

	for _, media := range medias {
		ids = append(ids, media.ID)
	}

	return ids
}
