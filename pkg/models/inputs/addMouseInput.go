package inputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type AddMouseInput struct {
	ID          uint64              `json:"id"`
	Price       float64             `json:"price"`
	Name        string              `json:"name"`
	Stock       uint64              `json:"stock"`
	TypeMouses  string              `json:"type_mouses"`
	Dpi         uint64              `json:"dpi"`
	ReleaseYear uint64              `json:"release_year"`
	Medias      []models.InputMedia `json:"medias"`
}
