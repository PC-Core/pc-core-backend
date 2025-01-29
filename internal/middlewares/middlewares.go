package middlewares

import (
	"net/http"

	"github.com/Core-Mouse/cm-backend/internal/auth"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/errors"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/middlewares/merrors"
	"github.com/Core-Mouse/cm-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware func(auth auth.Auth) gin.HandlerFunc

func JWTAuthorize(auth auth.Auth) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := helpers.GetAutorizationToken(ctx, helpers.BearerPrefix)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.IntoPublic()})
			ctx.Abort()
			return
		}

		data, err := auth.Authorize(token)

		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.IntoPublic()})
			ctx.Abort()
			return
		}

		ctx.Set(helpers.UserDataKey, data)

		ctx.Next()
	}
}

func RoleCheck(required models.UserRole, db database.DbController, caster helpers.RoleCastFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, exists := ctx.Get(helpers.UserDataKey)

		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.NewInternalSecretError()})
			ctx.Abort()
			return
		}

		role, err := caster(data)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.NewInternalSecretError()})
			ctx.Abort()
			return
		}

		if role != required {
			ctx.JSON(http.StatusForbidden, gin.H{"error": merrors.NewLowerRoleError(required, role)})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
