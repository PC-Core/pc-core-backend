package jwt

import (
	"fmt"
	"time"

	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/Core-Mouse/cm-backend/internal/models/outputs"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAuth struct {
	key []byte
}

func NewJWTAuth(key []byte) *JWTAuth {
	return &JWTAuth{
		key,
	}
}

func (a *JWTAuth) CreateRefreshToken(id int, rdur time.Duration) (string, error) {
	token := jwt.NewWithClaims(JWTTokenCryptoMethod, NewJWTRefreshClaimsFromID(id, rdur))
	jwt, err := token.SignedString(a.key)

	if err != nil {
		return "", err
	}

	return jwt, err
}

func (a *JWTAuth) CreateAccessToken(data *models.PublicUser, adur time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewJWTAccessClaimsFromUser(data, adur))
	jwt, err := token.SignedString(a.key)

	if err != nil {
		return "", err
	}

	return jwt, err
}

func (a *JWTAuth) Authentificate(data *models.PublicUser) (interface{}, error) {
	return a.AuthentificateWithDur(data, time.Duration(JWTAccessLifeTimeHours * time.Hour), time.Duration(JWTRefreshLifeTimeHours * time.Hour))
}

func (a *JWTAuth) AuthentificateWithDur(data *models.PublicUser, adur time.Duration, rdur time.Duration) (interface{}, error) {
	access, err := a.CreateAccessToken(data, adur)

	if err != nil {
		return nil, err
	}

	refresh, err := a.CreateRefreshToken(data.ID, rdur)

	if err != nil {
		return nil, err
	}

	return outputs.NewJWTPair(access, refresh), nil
}

func (a *JWTAuth) parsePairType(data interface{}) (*outputs.JWTPair, error) {
	pair, ok := data.(*outputs.JWTPair)

	if !ok {
		return nil, fmt.Errorf("wrong auth type provided")
	}

	return pair, nil
}

func (a *JWTAuth) parseJWT(token string, claims jwt.Claims) (*jwt.Token, error) {
	result, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	}, jwt.WithValidMethods([]string{JWTTokenCryptoMethod.Name}))

	if err != nil {
		return nil, err
	}

	return result, err
}

func validateJWT[T JWTTokenWithType](token string, req TokenType, a *JWTAuth, input_claims jwt.Claims) (*jwt.Token, error) {
	res, err := a.parseJWT(token, input_claims)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("jwt parser returned nil")
	}

	claims, ok := res.Claims.(T)

	if !ok {
		return nil, fmt.Errorf("wrong access params")
	}

	if claims.GetType() != req {
		return nil, fmt.Errorf("error token type: %s", claims.GetType())
	}

	return res, nil
}

func (a *JWTAuth) ValidateAccessJWT(access string) (*jwt.Token, error) {
	return validateJWT[*JWTAccessAuthClaims](access, AccessToken, a, &JWTAccessAuthClaims{})
}

func (a *JWTAuth) ValidateRefreshJWT(refresh string) (*jwt.Token, error) {
	return validateJWT[*JWTRefreshAuthClaims](refresh, RefreshToken, a, &JWTRefreshAuthClaims{})
}

func (a *JWTAuth) Authorize(data string) (interface{}, error) {
	tk, err := a.ValidateAccessJWT(data)

	if err != nil {
		return nil, err
	}

	access_claims, ok := tk.Claims.(*JWTAccessAuthClaims)

	if !ok {
		return nil, fmt.Errorf("wrong access token")
	}

	return access_claims, nil
}
