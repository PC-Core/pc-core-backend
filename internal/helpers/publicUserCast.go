package helpers

import (
	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/auth/jwt"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

const (
	UserDataKey = "user_data"
)

type PublicUserCaster func(from interface{}) (*models.PublicUser, errors.PCCError)

func JWTPublicUserCaster(auth auth.Auth) PublicUserCaster {
	return func(from interface{}) (*models.PublicUser, errors.PCCError) {
		claims, ok := from.(*jwt.JWTAccessAuthClaims)

		if !ok {
			return nil, errors.NewInternalSecretError()
		}

		return claims.IntoPublicUser(), nil
	}
}
