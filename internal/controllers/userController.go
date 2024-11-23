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
}

func (c *UserController) registerUser(ctx *gin.Context) {
	var params map[string]interface{}

	if err := ctx.BindJSON(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	name := params["name"].(string)
	email := params["email"].(string)
	role := params["role"].(models.UserRole)
	password := params["password"].(string)

	user, err := c.db.RegisterUser(name, email, role, password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	ctx.JSON(http.StatusCreated, user)
}
