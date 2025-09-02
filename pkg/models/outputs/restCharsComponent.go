package outputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type RestCharsComponent struct {
	Type   string                    `json:"type"`
	Values []models.CharsDescription `json:"values"`
	Info   any                       `json:"info"`
}

// func NewRestCharsComponent(ty string, values []models.CharsDescription, info any) *RestCharsComponent {
// 	return &RestCharsComponent{
// 		ty, values, info,
// 	}
// }
