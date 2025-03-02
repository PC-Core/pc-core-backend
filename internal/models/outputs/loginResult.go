package outputs

import "github.com/PC-Core/pc-core-backend/internal/models"

type TokensMap map[string]string

// Result of the login
type LoginResult struct {
	User   *models.PublicUser `json:"user"`
	Tokens TokensMap          `json:"tokens"`
}

func NewLoginResult(user *models.PublicUser, tokens TokensMap) *LoginResult {
	return &LoginResult{
		user, tokens,
	}
}
