package controllers

import (
	"net/http"
	"strconv"

	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/gin-gonic/gin"
)

type ReactionsController struct {
	engine          *gin.Engine
	db              database.DbController
	auth_middleware gin.HandlerFunc
	pucaster        helpers.PublicUserCaster
}

func NewReactionsController(engine *gin.Engine, db database.DbController, auth_middleware gin.HandlerFunc, pucaster helpers.PublicUserCaster) *ReactionsController {
	return &ReactionsController{
		engine,
		db,
		auth_middleware,
		pucaster,
	}
}

func (c *ReactionsController) ApplyRoutes() {
	gr := c.engine.Group("/reactions")
	{
		gr.POST("/:id", c.auth_middleware, c.setReaction)
	}
}

// Add, change or delete reaction from a comment
// @Summary      Add, change or delete reaction from a comment
// @Tags         reactions
// @Accept       json
// @Produce      json
// @Param 		 comment_id 	query	int						true	"ID of the comment"
// @param		 input			body	inputs.SetReactionInput	true	"Input"
// @Param		 Authorization  header	string					true	"access token for user is used to check your reaction, is not required"
// @Success      200  {object} 	int
// @Failure      400  {object}  errors.PublicPCCError
// @Failure      403  {object}  errors.PublicPCCError
// @Router       /reactions/:id [post]
func (c *ReactionsController) setReaction(ctx *gin.Context) {
	pu, err := GetPubUser(ctx, c.pucaster)

	if CheckErrorAndWriteUnauthorized(ctx, err) {
		return
	}

	idStr := ctx.Param("id")

	id, perr := strconv.ParseInt(idStr, 10, 64)

	if perr != nil && CheckErrorAndWriteBadRequest(ctx, errors.NewAtoiError(err)) {
		return
	}

	var input inputs.SetReactionInput

	if err := ctx.ShouldBindBodyWithJSON(&input); err != nil {
		CheckErrorAndWriteBadRequest(ctx, conerrors.BindErrorCast(err))
		return
	}

	result, err := c.db.SetReaction(id, int64(pu.ID), input.Type)

	if err != nil {
		CheckErrorAndWriteBadRequest(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}
