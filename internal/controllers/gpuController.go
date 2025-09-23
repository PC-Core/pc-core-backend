package controllers

import (
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/gin-gonic/gin"
)

type dbController struct {
	DB database.DbController
}

func (c *dbController) GetGpus(ctx *gin.Context) {
	gpus, err := c.DB.GetGpus()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gpus)
}