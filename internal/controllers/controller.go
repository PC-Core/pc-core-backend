package controllers

import (
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	ApplyRoutes()
}

type RequiresAuth interface {
	GetPUCaster() helpers.PublicUserCaster
	GetPubUser(ctx *gin.Context) (*models.PublicUser, error)
}
