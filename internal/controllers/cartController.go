package controllers

import (
	"fmt"
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/Core-Mouse/cm-backend/internal/redis"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	engine          *gin.Engine
	db              *database.DbController
	rctrl           *redis.RedisController
	pucaster        helpers.PublicUserCaster
	auth_middleware gin.HandlerFunc
}

func NewCartController(engine *gin.Engine, db *database.DbController, rctrl *redis.RedisController, pucaster helpers.PublicUserCaster, auth_middleware gin.HandlerFunc) *CartController {
	return &CartController{
		engine, db, rctrl, pucaster, auth_middleware,
	}
}

func (c *CartController) GetPUCaster() helpers.PublicUserCaster {
	return c.pucaster
}

func (c *CartController) GetPubUser(ctx *gin.Context) (*models.PublicUser, error) {
	userdata, exists := ctx.Get(helpers.UserDataKey)

	if !exists {
		return nil, fmt.Errorf("no user data found")
	}

	pu, err := c.pucaster(userdata)

	if err != nil {
		return nil, err
	}

	return pu, err
}

func (c *CartController) ApplyRoutes() {
	gr := c.engine.Group("/cart", c.auth_middleware)
	{
		gr.GET("/", c.getCart)
		gr.POST("/", c.addToCart)
		gr.DELETE("/", c.removeFromCart)
		gr.PUT("/", c.changeQuantity)
	}
}

func (c *CartController) getTempCart(pu *models.PublicUser) (*models.Cart, error) {
	arr, err := c.rctrl.GetCart(uint64(pu.ID))

	if err != nil {
		return nil, err
	}

	products, err := c.db.LoadProductsRangeAsCartItem(arr)

	if err != nil {
		return nil, err
	}

	return models.NewCart(0, uint64(pu.ID), products), nil
}

func (c *CartController) getDefaultCart(pu *models.PublicUser) (*models.Cart, error) {
	return c.db.GetCartByUserID(uint64(pu.ID))
}

func (c *CartController) addToTempCart(pu *models.PublicUser, input inputs.AddToCartInput) (uint64, error) {
	id, err := c.rctrl.AddToCart(uint64(pu.ID), input.ProductID, uint(input.Quantity))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CartController) addToDefaultCart(pu *models.PublicUser, input inputs.AddToCartInput) (uint64, error) {
	return c.db.AddToCart(input.ProductID, uint64(pu.ID), uint64(input.Quantity))
}

// Get user's cart
// @Summary      Get user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string	true	"access token for authorization"
// @Success      200  {object}  models.Cart
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /cart/ [get]
func (c *CartController) getCart(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	var cart *models.Cart

	switch pu.Role {
	case models.Temporary:
		cart, err = c.getTempCart(pu)
	default:
		cart, err = c.getDefaultCart(pu)
	}

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

// Add product to a cart
// @Summary      Add product to a cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string					true	"access token for authorization"
// @Param		 obj  			body  	inputs.AddToCartInput	true	"info about a product to add"
// @Success      201  {object}  uint64
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /cart/ [post]
func (c *CartController) addToCart(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	var input inputs.AddToCartInput

	err = ctx.ShouldBindBodyWithJSON(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	var product_id uint64

	switch pu.Role {
	case models.Temporary:
		product_id, err = c.addToTempCart(pu, input)
	default:
		product_id, err = c.addToDefaultCart(pu, input)
	}

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"product_id": product_id})
}

// Remove the product from the cart
// @Summary      Remove the product from the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string						true	"access token for authorization"
// @Param		 obj  			body  	inputs.RemoveFromCartInput	true	"info about a product to remove"
// @Success      200  {object}  uint64
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /cart/ [delete]
func (c *CartController) removeFromCart(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	var input inputs.RemoveFromCartInput

	err = ctx.ShouldBindBodyWithJSON(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	product_id, err := c.db.RemoveFromCart(input.ProductID, uint64(pu.ID))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_id": product_id})
}

// Update product's quantity
// @Summary      Update product's quantity
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string					true	"access token for authorization"
// @Param		 obj  			body  	inputs.AddToCartInput	true	"info about a product to add"
// @Success      200  {object}  uint64
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /cart/ [put]
func (c *CartController) changeQuantity(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	var input inputs.AddToCartInput

	err = ctx.ShouldBindBodyWithJSON(&input)

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	product_id, err := c.db.ChangeQuantity(input.ProductID, uint64(pu.ID), int64(input.Quantity))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_id": product_id})
}
