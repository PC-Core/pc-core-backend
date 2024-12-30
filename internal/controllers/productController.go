package controllers

import (
	"net/http"
	"strconv"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	engine *gin.Engine
	db     *database.DbController
}

func NewProductController(engine *gin.Engine, db *database.DbController) *ProductController {
	return &ProductController{
		engine, db,
	}
}

func (c *ProductController) ApplyRoutes() {
	c.engine.GET("/products/", c.getProducts)
	c.engine.GET("/products/:id", c.getProductById)
}

func (c *ProductController) getProducts(ctx *gin.Context) {
	var input inputs.GetProductsInput

	err := ctx.ShouldBindQuery(&input)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	start := (input.Page * input.Count) - input.Count

	products, err := c.db.GetProducts(start, input.Count)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func (c *ProductController) getProductById(ctx *gin.Context) {
	ids := ctx.Param("id")

	id, err := strconv.ParseUint(ids, 10, 64)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	product, err := c.db.GetProductById(id)

	if CheckErrorAndWrite(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}
