package jwt

type JWTTokenWithType interface {
	GetType() TokenType
}
