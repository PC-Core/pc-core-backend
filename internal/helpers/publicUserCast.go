package helpers

import (
	"github.com/PC-Core/pc-core-backend/internal/auth"
	"github.com/PC-Core/pc-core-backend/internal/auth/jwt"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
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
