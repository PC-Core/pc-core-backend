package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckErrorAndWrite(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return true
	}

	return false
}
