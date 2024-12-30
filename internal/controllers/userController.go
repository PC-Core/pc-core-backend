package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	engine *gin.Engine
	db     *database.DbController
}

func NewUserController(engine *gin.Engine, db *database.DbController) *UserController {
	return &UserController{
		engine, db,
	}
}

func (c *UserController) ApplyRoutes() {
	c.engine.POST("/users/register", c.registerUser)
	c.engine.GET("/users/login", c.loginUser)
}

func (c *UserController) registerUser(ctx *gin.Context) {
	var input inputs.RegisterUserInput

	err := ctx.ShouldBindJSON(&input)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	user, err := c.db.RegisterUser(input.Name, input.Email, input.Password)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) loginUser(ctx *gin.Context) {
	var input inputs.LoginUserInput

	err := ctx.ShouldBindQuery(&input)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	user, err := c.db.LoginUser(input.Email, input.Password)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, user)
}
