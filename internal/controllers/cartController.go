package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/redis"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/gin-gonic/gin"
)

type CartController struct {
	engine          *gin.Engine
	db              database.DbController
	rctrl           *redis.RedisController
	pucaster        helpers.PublicUserCaster
	auth_middleware gin.HandlerFunc
}

func NewCartController(engine *gin.Engine, db database.DbController, rctrl *redis.RedisController, pucaster helpers.PublicUserCaster, auth_middleware gin.HandlerFunc) *CartController {
	return &CartController{
		engine, db, rctrl, pucaster, auth_middleware,
	}
}

func (c *CartController) GetPUCaster() helpers.PublicUserCaster {
	return c.pucaster
}

func (c *CartController) GetPubUser(ctx *gin.Context) (*models.PublicUser, errors.PCCError) {
	userdata, exists := ctx.Get(helpers.UserDataKey)

	if !exists {
		return nil, conerrors.GetUserDataFromContextError()
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
		gr.POST("/item/:id", func(ctx *gin.Context) { c.setToCart(ctx) })
		gr.DELETE("/:id", c.removeFromCart)
		gr.DELETE("/item/:id", c.removeQuantity)
		gr.PUT("/item/:id", c.addQuantity)
	}
}

func (c *CartController) getTempCart(pu *models.PublicUser) (*models.Cart, errors.PCCError) {
	arr, err := c.rctrl.GetCart(uint64(pu.ID))

	if err != nil {
		return nil, err
	}

	products, err := c.db.LoadProductsRangeAsCartItem(arr)

	if err != nil {
		return nil, err
	}

	return models.NewCart(uint64(pu.ID), products), nil
}

func (c *CartController) getDefaultCart(pu *models.PublicUser) (*models.Cart, errors.PCCError) {
	return c.db.GetCartByUserID(uint64(pu.ID))
}

func (c *CartController) addToTempCart(pu *models.PublicUser, productID int, input inputs.AddToCartInput) (uint64, errors.PCCError) {
	id, err := c.rctrl.AddToCart(uint64(pu.ID), uint64(productID), uint(input.Quantity))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CartController) addToDefaultCart(pu *models.PublicUser, productID int, input inputs.AddToCartInput) (uint64, errors.PCCError) {
	return c.db.AddToCart(uint64(productID), uint64(pu.ID), uint64(input.Quantity))
}

func (c *CartController) setToDefaultCart(pu *models.PublicUser, productID int, input inputs.AddToCartInput) (uint64, errors.PCCError) {
	return c.db.SetToCart(uint64(productID), uint64(pu.ID), uint64(input.Quantity))
}

// Get user's cart
// @Summary      Get user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string	true	"access token for authorization"
// @Success      200  {object}  models.Cart
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
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

// Set the product with the provided quantity to a cart
// @Summary      Set the product with the provided quantity to a cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string					true	"access token for authorization"
// @Param		 obj  			body  	inputs.AddToCartInput	true	"info about a product to add"
// @Param		 Product ID		query	int						true	"id of the product to add to cart"
// @Success      201  {object}  uint64
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /cart/item/:id [post]
func (c *CartController) setToCart(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	idStr := ctx.Param("id")
	productID, eerr := strconv.Atoi(idStr)
	if err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(eerr))
		return
	}

	var input inputs.AddToCartInput

	if berr := ctx.ShouldBindBodyWithJSON(&input); berr != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(berr))
		return
	}

	var product_id uint64

	switch pu.Role {
	case models.Temporary:
		product_id, err = c.addToTempCart(pu, productID, input)
	default:
		product_id, err = c.setToDefaultCart(pu, productID, input)
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
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /cart/ [delete]
func (c *CartController) removeFromCart(ctx *gin.Context) {
	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	idStr := ctx.Param("id")
	id, eerr := strconv.Atoi(idStr)
	if err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(eerr))
		return
	}

	product_id, err := c.db.RemoveFromCart(uint64(id), uint64(pu.ID))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_id": product_id})
}

func (c *CartController) changeQuantity(productID int, pu *models.PublicUser, quantity int) (uint64, errors.PCCError) {
	return c.db.ChangeQuantity(uint64(productID), uint64(pu.ID), int64(quantity))
}

func (c *CartController) reqChangeQuantity(ctx *gin.Context, sign int) {
	if sign != 1 && sign != -1 {
		ctx.JSON(http.StatusInternalServerError, errors.NewInternalSecretError())
	}

	pu, err := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	idStr := ctx.Param("id")
	id, eerr := strconv.Atoi(idStr)
	if err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(eerr))
		return
	}

	var input inputs.AddToCartInput

	if berr := ctx.ShouldBindBodyWithJSON(&input); berr != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(berr))
		return
	}

	product_id, err := c.changeQuantity(id, pu, sign*int(math.Abs(float64(input.Quantity))))

	if CheckErrorAndWriteBadRequest(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"product_id": product_id})
}

// Remove the requested quantity of items from the cart
// @Summary      Remove the requested quantity of items from the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string					true	"access token for authorization"
// @Param		 obj  			body  	inputs.AddToCartInput	true	"info about a product to add"
// @Success      200  {object}  uint64
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /cart/item/:id [delete]
func (c *CartController) removeQuantity(ctx *gin.Context) {
	c.reqChangeQuantity(ctx, -1)
}

// Add the requested quantity of items to the cart
// @Summary      Add the requested quantity of items to the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param 		 Authorization	header	string					true	"access token for authorization"
// @Param		 obj  			body  	inputs.AddToCartInput	true	"info about a product to add"
// @Success      200  {object}  uint64
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      401  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /cart/item/:id [put]
func (c *CartController) addQuantity(ctx *gin.Context) {
	c.reqChangeQuantity(ctx, 1)
}
