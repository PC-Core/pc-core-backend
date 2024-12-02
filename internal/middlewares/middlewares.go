package middlewares

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/Core-Mouse/cm-backend/internal/models/inputs"
	"github.com/gin-gonic/gin"
)

func RoleCheck(required models.UserRole, db *database.DbController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input inputs.LoginUserInput;

		if err := ctx.ShouldBindQuery(&input); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Wrong user auth data sent"})
			ctx.Abort()
			return
		}

		if err := db.AuthentificateWithRole(input.Email, input.Password, required); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have the required role"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}