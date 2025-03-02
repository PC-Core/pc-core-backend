package middlewares

import (
	"net/http"

	"github.com/PC-Core/pc-core-backend/internal/auth"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/middlewares/merrors"
	"github.com/PC-Core/pc-core-backend/internal/models"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware func(auth auth.Auth) gin.HandlerFunc

func checkJWTRefresh(ctx *gin.Context) {
	// tk, err := ctx.Request.Cookie(helpers.REFRESH_COOKIE_NAME)

	// if err != nil {
	// 	return
	// }
}

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
