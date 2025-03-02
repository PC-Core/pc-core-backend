package helpers

import (
	"fmt"

	"github.com/PC-Core/pc-core-backend/internal/auth/jwt"
	"github.com/PC-Core/pc-core-backend/internal/models"
)

type RoleCastFunc = func(data interface{}) (models.UserRole, error)

func JWTRoleCast(data interface{}) (models.UserRole, error) {
	user, ok := data.(*jwt.JWTAccessAuthClaims)

	if !ok {
		return "", fmt.Errorf("jwt user data has a wrong type")
	}

	return user.Role, nil
}
