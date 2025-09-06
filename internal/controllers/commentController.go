package controllers

import (
	"net/http"
	"strconv"

	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	engine                       *gin.Engine
	db                           database.DbController
	auth_middleware              gin.HandlerFunc
	auth_not_required_middleware gin.HandlerFunc
	pucaster                     helpers.PublicUserCaster
}

func NewCommentController(engine *gin.Engine, db database.DbController, auth_middleware gin.HandlerFunc, auth_not_req_middleware gin.HandlerFunc, pucaster helpers.PublicUserCaster) *CommentController {
	return &CommentController{
		engine,
		db,
		auth_middleware,
		auth_not_req_middleware,
		pucaster,
	}
}

func (c *CommentController) ApplyRoutes() {
	g := c.engine.Group("/comment")
	{
		g.GET("/product/:id", c.auth_not_required_middleware, c.getRootComments)
		g.GET("/parent/:id", c.auth_not_required_middleware, c.getAnswers)
		g.POST("/product/:id", c.auth_middleware, c.addComment)
	}
}

func (c *CommentController) GetNotRequiredUserID(ctx *gin.Context) *int64 {
	var userID *int64 = nil

	userData, exists := ctx.Get(helpers.UserDataKey)

	if exists {
		data, err := c.pucaster(userData)

		if err == nil && data.Role != models.Temporary {
			i := int64(data.ID)
			userID = &i
		}
	}

	return userID
}

func (c *CommentController) GetPubUser(ctx *gin.Context) (*models.PublicUser, errors.PCCError) {
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

// Get root comments
// @Summary      Get root comments
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param 		 product_id 	query	int			true	"ID of the product"
// @Param		 Authorization  header	string		false	"access token for user is used to check your reaction, is not required"
// @Success      200  {array}  	models.Comment
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /comment/product/:id [get]
func (c *CommentController) getRootComments(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil && CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	userID := c.GetNotRequiredUserID(ctx)

	comments, perr := c.db.GetRootCommentsForProduct(int64(id), userID)

	if CheckErrorAndWriteBadRequest(ctx, perr) {
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

// Get answers on comment
// @Summary      Get answers on comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param 		 comment_id 	query	int						true	"ID of the comment"
// @param		 input			body	inputs.GetAnswersInput	true	"Input"
// @Param		 Authorization  header	string					false	"access token for user is used to check your reaction, is not required"
// @Success      200  {array}  	models.Comment
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /comment/parent/:id [get]
func (c *CommentController) getAnswers(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil && CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	var input inputs.GetAnswersInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(err))
		return
	}

	userID := c.GetNotRequiredUserID(ctx)

	ans, perr := c.db.GetAnswersOnComment(input.ProductID, userID, int64(id))

	if CheckErrorAndWriteBadRequest(ctx, perr) {
		return
	}

	ctx.JSON(http.StatusOK, ans)
}

// Add comment
// @Summary      Add comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param 		 product_id 	query	int						true	"ID of the product"
// @Param		 Authorization  header	string					true	"access token"
// @Param		 input			body	inputs.AddCommentInput	true	"input"
// @Success      201  {object} 	int
// @Failure      400  {object}  errors.PublicPCCError
// @Router       /comment/product/:id [post]
func (c *CommentController) addComment(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil && CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	data, perr := c.GetPubUser(ctx)

	if CheckErrorAndWriteUnauthorized(ctx, perr) {
		return
	}

	var input inputs.AddCommentInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(err))
		return
	}

	newID, perr := c.db.AddComment(&input, int64(data.ID), int64(id))

	if CheckErrorAndWriteBadRequest(ctx, perr) {
		return
	}

	ctx.JSON(http.StatusCreated, newID)
}
