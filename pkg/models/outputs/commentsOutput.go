package outputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type CommentsOutput struct {
	Comments           []models.Comment `json:"comments"`
	TotalCommentsCount int64            `json:"total_comments_count"`
}
