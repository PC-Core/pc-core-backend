package helpers

import (
	"fmt"

	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/auth/jwt"
	"github.com/Core-Mouse/cm-backend/internal/models"
)

const (
	UserDataKey = "user_data"
)

type PublicUserCaster func(from interface{}) (*models.PublicUser, error)

func JWTPublicUserCaster(auth auth.Auth) PublicUserCaster {
	return func(from interface{}) (*models.PublicUser, error) {
		claims, ok := from.(*jwt.JWTAccessAuthClaims)

		if !ok {
			return nil, fmt.Errorf("wrong claims type! Maybe you are using a wrong token type")
		}

		return claims.IntoPublicUser(), nil
	}
}
