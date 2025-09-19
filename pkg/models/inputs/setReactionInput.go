package inputs

import "github.com/PC-Core/pc-core-backend/pkg/models"

type SetReactionInput struct {
	Type models.ReactionType `json:"type"`
}
