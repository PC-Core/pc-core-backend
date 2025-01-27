package database

import "github.com/Core-Mouse/cm-backend/internal/models"

func (c *DPostgresDbController) GetCategories() ([]models.Category, error) {
	cats := make([]models.Category, 0, 5)

	res, err := c.db.Query("SELECT * FROM Categories")

	if err != nil {
		return nil, err
	}

	defer res.Close()

	for res.Next() {
		var cat models.Category

		if err := res.Scan(&cat.ID, &cat.Title, &cat.Description, &cat.Icon, &cat.Slug); err != nil {
			return nil, err
		}

		cats = append(cats, cat)
	}

	return cats, nil
}
