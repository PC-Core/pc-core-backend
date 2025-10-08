package inputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type AddKeyBoardInput struct {
	ID            uint64              `json:"id"`
	Price         float64             `json:"price"`
	Name          string              `json:"name"`
	Stock         uint64              `json:"stock"`
	TypeKeyBoards string              `json:"type_keyboards"`
	Switches      string            `json:"switches"`
	ReleaseYear   uint64              `json:"release_year"`
	Medias        []models.InputMedia `json:"medias"`
}
