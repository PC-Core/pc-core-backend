package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTAccessLifeTime  = 15 * time.Minute
	JWTRefreshLifeTime = 24 * 30 * time.Hour
)

var JWTTokenCryptoMethod = jwt.SigningMethodHS256

type TokenType string

const (
	RefreshToken TokenType = "Refresh"
	AccessToken  TokenType = "Access"
)
