package outputs

import "github.com/Core-Mouse/cm-backend/internal/models"

type AuthMap map[string]string

// Result of the login
type LoginResult struct {
	User *models.PublicUser `json:"user"`
	Auth AuthMap            `json:"auth"`
}

func NewLoginResult(user *models.PublicUser, auth AuthMap) *LoginResult {
	return &LoginResult{
		user, auth,
	}
}
