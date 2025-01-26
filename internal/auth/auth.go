package auth

import (
	"time"

	"github.com/Core-Mouse/cm-backend/internal/models"
)

type Auth interface {
	Authentificate(data *models.PublicUser) (interface{}, error)
	AuthentificateWithDur(data *models.PublicUser, adur time.Duration, rdur time.Duration) (interface{}, error)
	Authorize(data string) (interface{}, error)
}
