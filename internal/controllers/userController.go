package controllers

import (
	"net/http"

	"github.com/PC-Core/pc-core-backend/internal/auth"
	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/redis"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/PC-Core/pc-core-backend/pkg/models/outputs"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	engine *gin.Engine
	db     database.DbController
	rctrl  *redis.RedisController
	auth   auth.Auth
}

const CookieUseHttps = true

func NewUserController(engine *gin.Engine, db database.DbController, rctrl *redis.RedisController, auth auth.Auth) *UserController {
	return &UserController{
		engine, db, rctrl, auth,
	}
}

func (c *UserController) ApplyRoutes() {
	c.engine.POST("/users/register", c.registerUser)
	c.engine.POST("/users/login", c.loginUser)
	c.engine.POST("/users/temp/new", c.createTempUser)
	c.engine.GET("/users/logout", c.logoutUser)
}

// Register a new User
// @Summary      Register a new User
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 user	body inputs.RegisterUserInput	true	"User data to register"
// @Success      200  {object}  outputs.LoginResult
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /users/register [post]
func (c *UserController) registerUser(ctx *gin.Context) {
	var input inputs.RegisterUserInput

	if berr := ctx.ShouldBindJSON(&input); berr != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(berr))
		return
	}

	user, err := c.db.RegisterUser(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	res, err := c.auth.Authentificate(models.NewPublicUserFromUser(user))

	if err != nil {
		CheckErrorAndWriteBadRequest(ctx, errors.NewInternalSecretError())
	}

	sendAuthData(ctx, res, http.StatusCreated, models.NewPublicUserFromUser(user), input.Remember)
}

// Login
// @Summary      Login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 user	body inputs.LoginUserInput	true	"User data to login"
// @Success      200  {object}  outputs.LoginResult
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /users/login [post]
func (c *UserController) loginUser(ctx *gin.Context) {
	var input inputs.LoginUserInput
	var err errors.PCCError

	if berr := ctx.ShouldBindJSON(&input); berr != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(berr))
		return
	}

	user, err := c.db.LoginUser(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	res, err := c.auth.Authentificate(models.NewPublicUserFromUser(user))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	sendAuthData(ctx, res, http.StatusOK, models.NewPublicUserFromUser(user), input.Remember)
}

func sendAuthData(ctx *gin.Context, ad *models.AuthData, status int, user *models.PublicUser, remember *bool) {
	setRefreshCookie(ctx, ad.GetPrivate().String(), remember, int(auth.AuthPrivateCookieLifetime.Seconds()))
	ctx.JSON(status, outputs.NewLoginResult(user, outputs.TokensMap{"access": ad.GetPublic().String()}))
}

func (c *UserController) createTempUser(ctx *gin.Context) {
	res, err := c.rctrl.CreateTempUser(c.auth)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func setRefreshCookie(ctx *gin.Context, refresh string, remember *bool, maxtime int) {
	if remember != nil && *remember {
		ctx.SetSameSite(http.SameSiteNoneMode)
		ctx.SetCookie(helpers.RefreshCookieName, refresh, maxtime, "/", "", CookieUseHttps, true)
	}
}

// Logout
// @Summary      Logout
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {string}	ok
// @Router       /users/logout [get]
func (c *UserController) logoutUser(ctx *gin.Context) {
	remember := true
	setRefreshCookie(ctx, "", &remember, -1)
	// ctx.SetSameSite(http.SameSiteNoneMode)
	// ctx.SetCookie(helpers.RefreshCookieName, "", -1, "/", "", CookieUseHttps, true)
	ctx.JSON(http.StatusOK, "ok")
}
