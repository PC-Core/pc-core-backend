package inputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type AddLaptopInput struct {
	Name   string              `json:"name"`
	CpuID  uint64              `json:"cpu"`
	Ram    int16               `json:"ram"`
	GpuID  uint64              `json:"gpu"`
	Price  float64             `json:"price"`
	Stock  uint64              `json:"stock"`
	Medias []models.InputMedia `json:"medias"`
}
