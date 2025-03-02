package jwt

import (
	"strconv"
	"time"

	"github.com/PC-Core/pc-core-backend/internal/auth"
	"github.com/PC-Core/pc-core-backend/internal/auth/jwt/jerrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/PC-Core/pc-core-backend/internal/models/outputs"
	"github.com/golang-jwt/jwt/v5"
)

type StrWrapper string

func (s StrWrapper) String() string {
	return string(s)
}

type JWTAuth struct {
	key []byte
}

func NewJWTAuth(key []byte) *JWTAuth {
	return &JWTAuth{
		key,
	}
}

func (a *JWTAuth) CreateRefreshToken(id int, rdur time.Duration) (string, errors.PCCError) {
	token := jwt.NewWithClaims(JWTTokenCryptoMethod, NewJWTRefreshClaimsFromID(id, rdur))
	jwt, err := token.SignedString(a.key)

	if err != nil {
		return "", jerrors.JwtErrorCaster(err)
	}

	return jwt, nil
}

func (a *JWTAuth) CreateAccessToken(data *models.PublicUser, adur time.Duration) (string, errors.PCCError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, NewJWTAccessClaimsFromUser(data, adur))
	jwt, err := token.SignedString(a.key)

	if err != nil {
		return "", jerrors.JwtErrorCaster(err)
	}

	return jwt, nil
}

func (a *JWTAuth) Authentificate(data *models.PublicUser) (*models.AuthData, errors.PCCError) {
	return a.AuthentificateWithDur(data, time.Duration(auth.AuthPublicLifetime), time.Duration(auth.AuthPrivateCookieLifetime))
}

func (a *JWTAuth) AuthentificateWithDur(data *models.PublicUser, adur time.Duration, rdur time.Duration) (*models.AuthData, errors.PCCError) {
	access, err := a.CreateAccessToken(data, adur)

	if err != nil {
		return nil, err
	}

	refresh, err := a.CreateRefreshToken(data.ID, rdur)

	if err != nil {
		return nil, err
	}

	return models.NewAuthData(StrWrapper(access), StrWrapper(refresh)), nil
}

func (a *JWTAuth) parsePairType(data interface{}) (*outputs.JWTPair, errors.PCCError) {
	pair, ok := data.(*outputs.JWTPair)

	if !ok {
		return nil, errors.NewInternalSecretError()
	}

	return pair, nil
}

func (a *JWTAuth) parseJWT(token string, claims jwt.Claims) (*jwt.Token, errors.PCCError) {
	result, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return a.key, nil
	}, jwt.WithValidMethods([]string{JWTTokenCryptoMethod.Name}))

	if err != nil {
		return nil, jerrors.JwtErrorCaster(err)
	}

	return result, nil
}

func validateJWT[T JWTTokenWithType](token string, req TokenType, a *JWTAuth, input_claims jwt.Claims) (*jwt.Token, errors.PCCError) {
	res, err := a.parseJWT(token, input_claims)

	if err != nil {
		return nil, err
	}

	claims, ok := res.Claims.(T)

	if !ok {
		return nil, errors.NewInternalSecretError()
	}

	if claims.GetType() != req {
		return nil, jerrors.NewJwtTokenTypeError(claims.GetType())
	}

	return res, nil
}

func (a *JWTAuth) ValidateAccessJWT(access string) (*jwt.Token, errors.PCCError) {
	return validateJWT[*JWTAccessAuthClaims](access, AccessToken, a, &JWTAccessAuthClaims{})
}

func (a *JWTAuth) ValidateRefreshJWT(refresh string) (*jwt.Token, errors.PCCError) {
	return validateJWT[*JWTRefreshAuthClaims](refresh, RefreshToken, a, &JWTRefreshAuthClaims{})
}

func (a *JWTAuth) Authorize(data string) (interface{}, errors.PCCError) {
	tk, err := a.ValidateAccessJWT(data)

	if err != nil {
		return nil, err
	}

	access_claims, ok := tk.Claims.(*JWTAccessAuthClaims)

	if !ok {
		return nil, errors.NewInternalSecretError()
	}

	return access_claims, nil
}

func (a *JWTAuth) CheckAndReissue(token string) (string, errors.PCCError) {
	tk, err := a.ValidateRefreshJWT(token)

	if err != nil {
		return "", err
	}

	exp, ierr := tk.Claims.GetExpirationTime()

	if ierr != nil {
		return "", jerrors.JwtErrorCaster(ierr)
	}

	if time.Until(exp.Time) > 1*time.Minute {
		return token, nil
	}

	sub, ierr := tk.Claims.GetSubject()

	if ierr != nil {
		return "", jerrors.JwtErrorCaster(ierr)
	}

	isub, ierr := strconv.Atoi(sub)

	if ierr != nil {
		return "", errors.NewAtoiError(ierr)
	}

	return a.CreateRefreshToken(isub, auth.AuthPrivateCookieLifetime)
}
