package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)


var JWTTokenCryptoMethod = jwt.SigningMethodHS256

type TokenType string

const (
	RefreshToken TokenType = "Refresh"
	AccessToken  TokenType = "Access"
)
