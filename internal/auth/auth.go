package auth

import (
	"time"

	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

type Auth interface {
	Authentificate(data *models.PublicUser) (*models.AuthData, errors.PCCError)
	AuthentificateWithDur(data *models.PublicUser, adur time.Duration, rdur time.Duration) (*models.AuthData, errors.PCCError)
	Authorize(data string) (interface{}, errors.PCCError)
}

const (
	AuthPublicLifetime  = 15 * time.Minute
	AuthPrivateCookieLifetime = 24 * 30 * time.Hour
)
