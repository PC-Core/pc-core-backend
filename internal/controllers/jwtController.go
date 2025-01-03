package controllers

import (
	"fmt"
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/auth/jwt"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/gin-gonic/gin"
)

type JWTController struct {
	engine   *gin.Engine
	db       *database.DbController
	jwt_auth *jwt.JWTAuth
}

type SingleAccessToken struct {
	AccessToken string `json:"access_token"`
}

func NewJWTController(engine *gin.Engine, db *database.DbController, jwt_auth *jwt.JWTAuth) *JWTController {
	return &JWTController{
		engine, db, jwt_auth,
	}
}

func (c *JWTController) ApplyRoutes() {
	c.engine.POST("/auth/jwt/update", c.updateAccessToken)
}

// Update Access JWT token
// @Summary      Update Access JWT token
// @Tags         jwt
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header string	true	"User's refresh token"
// @Success      200  {object}  SingleAccessToken
// @Failure      401  {object}  map[string]interface{}
// @Failure		 400  {object}	map[string]interface{}
// @Router       /auth/jwt/update [post]
func (c *JWTController) updateAccessToken(ctx *gin.Context) {
	id, err := c.getUserIDFromRefreshHeader(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	user, err := c.db.GetUserByID(id)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	new_token, err := c.jwt_auth.CreateAccessToken(user)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, &SingleAccessToken{
		AccessToken: new_token,
	})

}

func (c *JWTController) getUserIDFromRefreshHeader(ctx *gin.Context) (int, error) {
	str_token, err := helpers.GetAutorizationToken(ctx, helpers.BearerPrefix)

	if err != nil {
		return -1, err
	}

	token, err := c.jwt_auth.ValidateRefreshJWT(str_token)

	if err != nil {
		return -1, err
	}

	claims, ok := token.Claims.(*jwt.JWTRefreshAuthClaims)

	if !ok {
		return -1, fmt.Errorf("wrong JWT claims type provided")
	}

	return claims.UserID, nil
}
