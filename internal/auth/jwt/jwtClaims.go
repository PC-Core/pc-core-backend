package jwt

import (
	"strconv"
	"time"

	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAccessAuthClaims struct {
	ID    int
	Name  string
	Email string
	Role  models.UserRole
	Type  TokenType
	jwt.RegisteredClaims
}

func NewJWTAccessClaimsFromUser(data *models.PublicUser) *JWTAccessAuthClaims {
	return &JWTAccessAuthClaims{
		ID:    data.ID,
		Name:  data.Name,
		Email: data.Email,
		Role:  data.Role,
		Type:  AccessToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Subject:   strconv.Itoa(data.ID),
		},
	}
}

func (c *JWTAccessAuthClaims) IntoPublicUser() *models.PublicUser {
	return models.NewPublicUser(c.ID, c.Name, c.Email, c.Role)
}

func (t *JWTAccessAuthClaims) GetType() TokenType {
	return t.Type
}

type JWTRefreshAuthClaims struct {
	UserID int
	Type   TokenType
	jwt.RegisteredClaims
}

func NewJWTRefreshClaimsFromID(id int) *JWTRefreshAuthClaims {
	return &JWTRefreshAuthClaims{
		UserID: id,
		Type:   RefreshToken,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTRefreshLifeTimeHours * time.Hour)),
		},
	}
}

func (t *JWTRefreshAuthClaims) GetType() TokenType {
	return t.Type
}
