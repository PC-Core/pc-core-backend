package controllers

import (
	"net/http"
	"strconv"

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
	c.engine.GET("/laptops/chars/:id", c.getChars)
}

func (c *LaptopController) addLaptop(ctx *gin.Context) {
	var input inputs.AddLaptopInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, chars, err := c.db.AddLaptop(input.Name, input.Price, 0, input.Stock, input.Cpu, input.Ram, input.Gpu)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"product": product,
		"chars":   chars,
	})
}

func (c *LaptopController) getChars(ctx *gin.Context) {
	ids := ctx.Param("id")

	id, err := strconv.ParseUint(ids, 10, 64)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	chars, err := c.db.GetProductCharsByProductID(id)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chars": chars,
	})
}
