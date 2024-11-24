package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models"
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
	var input RegisterUserInput

	err := ctx.ShouldBindJSON(&input)

	if checkErrorAndWrite(ctx, err) {
		return
	}

	user, err := c.db.RegisterUser(input.Name, input.Email, models.UserRole(input.Role), input.Password)

	if checkErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) loginUser(ctx *gin.Context) {
	var input LoginUserInput

	err := ctx.ShouldBindQuery(&input)

	if checkErrorAndWrite(ctx, err) {
		return
	}

	user, err := c.db.LoginUser(input.Email, input.Password)

	if checkErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func checkErrorAndWrite(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return true
	}

	return false
}
