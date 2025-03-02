package controllers

import (
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	ApplyRoutes()
}

type RequiresAuth interface {
	GetPUCaster() helpers.PublicUserCaster
	GetPubUser(ctx *gin.Context) (*models.PublicUser, error)
}
