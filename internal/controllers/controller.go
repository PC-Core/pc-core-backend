package controllers

import (
	"github.com/PC-Core/pc-core-backend/internal/controllers/conerrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	ApplyRoutes()
}

type RequiresAuth interface {
	GetPUCaster() helpers.PublicUserCaster
	GetPubUser(ctx *gin.Context) (*models.PublicUser, error)
}

func GetPubUser(ctx *gin.Context, pucaster helpers.PublicUserCaster) (*models.PublicUser, errors.PCCError) {
	userdata, exists := ctx.Get(helpers.UserDataKey)

	if !exists {
		return nil, conerrors.GetUserDataFromContextError()
	}

	pu, err := pucaster(userdata)

	if err != nil {
		return nil, err
	}

	return pu, err
}

func GetNotRequiredUserID(ctx *gin.Context, pucaster helpers.PublicUserCaster) *int64 {
	var userID *int64 = nil

	userData, exists := ctx.Get(helpers.UserDataKey)

	if exists {
		data, err := pucaster(userData)

		if err == nil && data.Role != models.Temporary {
			i := int64(data.ID)
			userID = &i
		}
	}

	return userID
}
