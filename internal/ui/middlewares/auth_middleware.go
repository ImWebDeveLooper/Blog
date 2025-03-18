package middlewares

import (
	"blog/assets/locales"
	"blog/internal/platform/pkg/jwt"
	"blog/internal/platform/pkg/lang"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const BearerSchema = "Bearer "

func AuthMiddleware(jwtService jwt.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.UnauthorizedUserError),
			})
			return
		}
		if !strings.HasPrefix(authHeader, BearerSchema) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.UnauthorizedUserError),
			})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, BearerSchema))
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.UnauthorizedUserError),
			})
			return
		}
		claims, ok := token.Claims.(*jwt.CustomClaims)
		if ok && token.Valid {
			ctx.Set("email", claims.Email)
			ctx.Set("userID", claims.UserID)
			ctx.Set("role", claims.Role)
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": lang.TryBy(ctx.GetString("locale"), locales.UnauthorizedUserError),
		})
		return
	}
}
