package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/controllers/conerrors"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	engine          *gin.Engine
	pucaster        helpers.PublicUserCaster
	auth_middleware gin.HandlerFunc
}

func NewProfileController(engine *gin.Engine, pucaster helpers.PublicUserCaster, auth_middleware gin.HandlerFunc) *ProfileController {
	return &ProfileController{
		engine, pucaster, auth_middleware,
	}
}

func (c *ProfileController) GetPubUser(ctx *gin.Context) (*models.PublicUser, errors.PCCError) {
	userdata, exists := ctx.Get(helpers.UserDataKey)

	if !exists {
		return nil, conerrors.GetUserDataFromContextError()
	}

	pu, err := c.pucaster(userdata)

	if err != nil {
		return nil, err
	}

	return pu, err
}

func (c *ProfileController) ApplyRoutes() {
	group := c.engine.Group("/profile")
	{
		group.GET("/", c.auth_middleware, c.getProfile)
	}
}

// Get user profile
// @Summary      Get user profile
// @Tags         profile
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string	true	"access token for authorization"
// @Success      200  {object}  models.Profile
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /profile/ [get]
func (c *ProfileController) getProfile(ctx *gin.Context) {
	user, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, models.NewProfile(user))
}
