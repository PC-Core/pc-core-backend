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

// Get products from page N in quantity M
// @Summary      Get products from page N in quantity M
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 product query	inputs.GetProductsInput	true	"Page and count"
// @Success      200  {array}  models.Product
// @Failure      400  {object}  map[string]interface{}
// @Router       /products/ [get]
func (c *ProductController) getProducts(ctx *gin.Context) {
	var input inputs.GetProductsInput

	err := ctx.ShouldBindQuery(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	start := (input.Page * input.Count) - input.Count

	products, err := c.db.GetProducts(start, input.Count)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// Get a single product by ID
// @Summary      Get a single product by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 id		path	uint64	true	"Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]interface{}
// @Router       /products/{id} [get]
func (c *ProductController) getProductById(ctx *gin.Context) {
	ids := ctx.Param("id")

	id, err := strconv.ParseUint(ids, 10, 64)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	product, err := c.db.GetProductById(id)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}
