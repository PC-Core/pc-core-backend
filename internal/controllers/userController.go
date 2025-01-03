package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	engine *gin.Engine
	db     *database.DbController
	auth   auth.Auth
}

func NewUserController(engine *gin.Engine, db *database.DbController, auth auth.Auth) *UserController {
	return &UserController{
		engine, db, auth,
	}
}

func (c *UserController) ApplyRoutes() {
	c.engine.POST("/users/register", c.registerUser)
	c.engine.GET("/users/login", c.loginUser)
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

	res, err := c.auth.Authentificate(user)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, res)
}
