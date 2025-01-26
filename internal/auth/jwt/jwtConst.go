package jwt

import "github.com/golang-jwt/jwt/v5"

const (
	JWTAccessLifeTimeHours   = 15 / 60
	JWTRefreshLifeTimeHours = 24 * 30
)

var JWTTokenCryptoMethod = jwt.SigningMethodHS256

type TokenType string

const (
	RefreshToken TokenType = "Refresh"
	AccessToken  TokenType = "Access"
)
