package controllers

import (
	"net/http"

	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/middlewares"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/gin-gonic/gin"
)

type MouseController struct {
	engine          *gin.Engine
	db              database.DbController
	auth_middleware gin.HandlerFunc
	caster          helpers.RoleCastFunc
}

func NewMouseController(engine *gin.Engine, db database.DbController, auth_middleware gin.HandlerFunc, caster helpers.RoleCastFunc) *MouseController{
	return &MouseController{
		engine, db, auth_middleware, caster,
	}
}

func (c *MouseController) ApplyRoutes(){ 
	c.engine.POST("/mouses/add", c.auth_middleware, middlewares.RoleCheck(models.Admin, c.db, c.caster), c.addMouse)
}

func (c *MouseController) addMouse(ctx *gin.Context){
	var input inputs.AddMouseInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, mouses, err := c.db.AddMouse(&input)

	if CheckErrorAndWriteBadRequest(ctx, err){
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"Product": product,
		"Mouse": mouses,
	})
}