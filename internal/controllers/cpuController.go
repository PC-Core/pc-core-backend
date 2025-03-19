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

type CpuController struct {
	engine          *gin.Engine
	db              database.DbController
	auth_middleware gin.HandlerFunc
	caster          helpers.RoleCastFunc
}

func NewCpuController(engine *gin.Engine, db database.DbController, auth_middleware gin.HandlerFunc, caster helpers.RoleCastFunc) *CpuController {
	return &CpuController{
		engine, db, auth_middleware, caster,
	}
}

func (c *CpuController) ApplyRoutes() {
	c.engine.POST("/cpus/add", c.auth_middleware, middlewares.RoleCheck(models.Admin, c.db, c.caster), c.addCpu)
}

// Add laptop
// @Summary      Add a new laptop
// @Tags         laptops
// @Accept       json
// @Produce      json
// @Param 		 laptop 		body	inputs.AddCpuInput	true	"Laptop data"
// @Param		 Authorization  header	string					true	"access token for user with Admin role"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /laptops/add [post]
func (c *CpuController) addCpu(ctx *gin.Context) {
	var input inputs.AddCpuInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, chars, err := c.db.AddCpu(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"product": product,
		"chars":   chars,
	})
}
