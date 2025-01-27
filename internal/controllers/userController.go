package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/Core-Mouse/cm-backend/internal/redis"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	engine *gin.Engine
	db     database.DbController
	rctrl  *redis.RedisController
	auth   auth.Auth
}

func NewUserController(engine *gin.Engine, db database.DbController, rctrl *redis.RedisController, auth auth.Auth) *UserController {
	return &UserController{
		engine, db, rctrl, auth,
	}
}

func (c *UserController) ApplyRoutes() {
	c.engine.POST("/users/register", c.registerUser)
	c.engine.GET("/users/login", c.loginUser)
	c.engine.POST("/users/temp/new", c.createTempUser)
}

// Register a new User
// @Summary      Register a new User
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 user	body inputs.RegisterUserInput	true	"User data to register"
// @Success      200  {object}  models.User
// @Failure      400  {object}  map[string]interface{}
// @Router       /users/register [post]
func (c *UserController) registerUser(ctx *gin.Context) {
	var input inputs.RegisterUserInput

	err := ctx.ShouldBindJSON(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	user, err := c.db.RegisterUser(input.Name, input.Email, input.Password)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

// Login
// @Summary      Login
// @Tags         users
// @Accept       json
// @Produce      json
// @Param 		 user	body inputs.LoginUserInput	true	"User data to login"
// @Success      200  {object}  outputs.JWTPair
// @Failure      400  {object}  map[string]interface{}
// @Router       /users/login [get]
func (c *UserController) loginUser(ctx *gin.Context) {
	var input inputs.LoginUserInput

	err := ctx.ShouldBindQuery(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	user, err := c.db.LoginUser(input.Email, input.Password)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	res, err := c.auth.Authentificate(models.NewPublicUserFromUser(user))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *UserController) createTempUser(ctx *gin.Context) {
	res, err := c.rctrl.CreateTempUser(c.auth)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
