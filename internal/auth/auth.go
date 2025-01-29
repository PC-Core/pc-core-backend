package auth

import (
	"time"

	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

type Auth interface {
	Authentificate(data *models.PublicUser) (interface{}, errors.PCCError)
	AuthentificateWithDur(data *models.PublicUser, adur time.Duration, rdur time.Duration) (interface{}, errors.PCCError)
	Authorize(data string) (interface{}, errors.PCCError)
}
