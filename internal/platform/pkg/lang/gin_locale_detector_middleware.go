package lang

import "github.com/gin-gonic/gin"

func GinLocaleDetectorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lang := ctx.GetHeader("Accept-Language")
		ctx.Set("locale", lang)
	}
}
