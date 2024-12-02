package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/middlewares"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

type LaptopController struct {
	engine *gin.Engine
	db     *database.DbController
}

func NewLaptopController(engine *gin.Engine, db *database.DbController) *LaptopController {
	return &LaptopController{
		engine, db,
	}
}

func (c *LaptopController) ApplyRoutes() {
	c.engine.POST("/laptops/add", middlewares.RoleCheck(models.Admin, c.db), c.addLaptop)
}

func (c *LaptopController) addLaptop(ctx *gin.Context) {
	var input inputs.AddLaptopInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	laptop, err := c.db.AddLaptop(input.Name, input.Cpu, input.Ram, input.Gpu, input.Price, input.Discount)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, laptop)
}
