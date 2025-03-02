package inputs

import "github.com/Core-Mouse/cm-backend/internal/models"

type AddLaptopInput struct {
	Name   string              `json:"name"`
	Cpu    string              `json:"cpu"`
	Ram    int16               `json:"ram"`
	Gpu    string              `json:"gpu"`
	Price  float64             `json:"price"`
	Stock  uint64              `json:"stock"`
	Medias []models.InputMedia `json:"medias"`
}
