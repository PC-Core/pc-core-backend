package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

func GetAutorizationToken(ctx *gin.Context, prefix string) (string, error) {
	authHeader := ctx.GetHeader(AuthorizationHeader)

	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", fmt.Errorf("authorization header missing or invalid")
	}

	return authHeader[len(prefix):], nil
}
