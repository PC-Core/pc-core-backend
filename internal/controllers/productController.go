package controllers

import (
	"net/http"
	"strconv"

	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	engine *gin.Engine
	db     database.DbController
}

func NewProductController(engine *gin.Engine, db database.DbController) *ProductController {
	return &ProductController{
		engine, db,
	}
}

func (c *ProductController) ApplyRoutes() {
	c.engine.GET("/products/", c.getProducts)
	c.engine.GET("/products/:id", c.getProductById)
	c.engine.GET("/products/chars/:id", c.getProductChars)
}

// Get products from page N in quantity M
// @Summary      Get products from page N in quantity M
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 product query	inputs.GetProductsInput	true	"Page and count"
// @Success      200  {array}  models.Product
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /products/ [get]
func (c *ProductController) getProducts(ctx *gin.Context) {
	var input inputs.GetProductsInput

	berr := ctx.ShouldBindQuery(&input)

	if berr != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(berr))
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
// @Failure      400  {object} errors.PublicPCCError
// @Router       /products/{id} [get]
func (c *ProductController) getProductById(ctx *gin.Context) {
	ids := ctx.Param("id")

	id, err := strconv.ParseUint(ids, 10, 64)

	if CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	product, perr := c.db.GetProductById(id)

	if CheckErrorAndWriteBadRequest(ctx, perr) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// Get product 	characteristics
// @Summary 	Get product chars
// @Tags 		products
// @Accept 		json
// @Produce 	json
// @Param		id	path	 int	true		"Product id"
// @Success 	200	{object} interface{}
// @Failure 	400 {object}  errors.PublicPCCError
// @Failure 	400 {object}  errors.PublicPCCError
// @Failure		400 {object}  errors.PublicPCCError
// @Router		/products/chars/{id} [get]
func (c *ProductController) getProductChars(ctx *gin.Context) {
	ids := ctx.Param("id")

	id, err := strconv.ParseUint(ids, 10, 64)

	if CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	chars, perr := c.db.GetProductCharsByProductID(id)

	if CheckErrorAndWriteBadRequest(ctx, perr) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"chars": chars,
	})
}
