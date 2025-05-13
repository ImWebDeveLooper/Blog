package actions

import (
	"blog/assets/locales"
	"blog/internal/domain/posts"
	"blog/internal/platform/dtos"
	"blog/internal/platform/pkg/lang"
	"blog/internal/platform/pkg/validators"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// PostCreatePostsAction handles HTTP POST requests to create a post and validate the input
func PostCreatePostsAction(interactor posts.Interactor, validator validators.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetString("locale")
		authorID := ctx.GetString("userID")
		var req dtos.CreatePostRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InvalidSchemaError),
			})
			return
		}
		result := validator.Validate(req, locale)
		if result.Fails {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.ValidationError),
				"errors":  result.Messages.All(),
			})
			return
		}
		err := interactor.CreatePost(ctx, authorID, req)
		if err != nil {
			log.WithError(err).Error("error saving post")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InternalServerError),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": lang.TryBy(ctx.GetString("locale"), locales.PostCreated),
		})
	}
}
