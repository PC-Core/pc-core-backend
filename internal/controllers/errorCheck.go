package controllers

import (
	"net/http"

	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/gin-gonic/gin"
)

func checkErrorAndWrite(ctx *gin.Context, err errors.PCCError, status int) bool {
	if err != nil {
		ctx.JSON(status, gin.H{
			"error": *err.IntoPublic(),
		})
		return true
	}

	return false
}

func CheckErrorAndWriteBadRequest(ctx *gin.Context, err errors.PCCError) bool {
	return checkErrorAndWrite(ctx, err, http.StatusBadRequest)
}

func CheckErrorAndWriteUnauthorized(ctx *gin.Context, err errors.PCCError) bool {
	return checkErrorAndWrite(ctx, err, http.StatusUnauthorized)
}
