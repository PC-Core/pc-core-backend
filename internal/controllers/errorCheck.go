package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkErrorAndWrite(ctx *gin.Context, err error, status int) bool {
	if err != nil {
		ctx.JSON(status, gin.H{
			"error": err.Error(),
		})
		return true
	}

	return false
}

func CheckErrorAndWriteBadRequest(ctx *gin.Context, err error) bool {
	return checkErrorAndWrite(ctx, err, http.StatusBadRequest)
}

func CheckErrorAndWriteUnauthorized(ctx *gin.Context, err error) bool {
	return checkErrorAndWrite(ctx, err, http.StatusUnauthorized)
}
