package auth

import "github.com/Core-Mouse/cm-backend/internal/models"

type Auth interface {
	Authentificate(data *models.User) (interface{}, error)
	Authorize(data string) (interface{}, error)
}
