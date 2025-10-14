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

type KeyBoardController struct {
	engine          *gin.Engine
	db              database.DbController
	auth_middleware gin.HandlerFunc
	caster          helpers.RoleCastFunc
}

func NewKeyBoardController(engine *gin.Engine, db database.DbController, auth_middleware gin.HandlerFunc, caster helpers.RoleCastFunc) *KeyBoardController{
	return &KeyBoardController{
		engine, db, auth_middleware, caster,
	}
}

func (c *KeyBoardController) ApplyRoutes(){ 
	c.engine.POST("/keyboards/add", c.auth_middleware, middlewares.RoleCheck(models.Admin, c.db, c.caster), c.addKeyBoard)
}

func (c *KeyBoardController) addKeyBoard(ctx *gin.Context){
	var input inputs.AddKeyBoardInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, keyboards, err := c.db.AddKeyBoard(&input)

	if CheckErrorAndWriteBadRequest(ctx, err){
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"Product": product,
		"Keyboard": keyboards,
	})
}