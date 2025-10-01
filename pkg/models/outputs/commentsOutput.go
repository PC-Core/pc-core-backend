package outputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type CommentsOutput struct {
	Comments []models.Comment `json:"comments"`
	Amount   int64            `json:"amount"`
}
