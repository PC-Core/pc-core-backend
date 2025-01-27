package controllers

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	engine *gin.Engine
	db     database.DbController
}

func NewCategoryController(engine *gin.Engine, db database.DbController) *CategoryController {
	return &CategoryController{
		engine, db,
	}
}

func (c *CategoryController) ApplyRoutes() {
	category := c.engine.Group("/categories")
	{
		category.GET("/", c.getAll)
	}
}

// Get all categories
// @Summary      Get all categories
// @Tags         categories
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Category
// @Failure      400  {object}  map[string]interface{}
// @Router       /categories/ [get]
func (c *CategoryController) getAll(ctx *gin.Context) {
	cats, err := c.db.GetCategories()

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, cats)
}
