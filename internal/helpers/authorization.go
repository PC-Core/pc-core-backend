package helpers

import (
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func GetAutorizationToken(ctx *gin.Context, prefix string) (string, errors.PCCError) {
	authHeader := ctx.GetHeader(AuthorizationHeader)

	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", errors.MissingHeader(AuthorizationHeader)
	}

	return authHeader[len(prefix):], nil
}
