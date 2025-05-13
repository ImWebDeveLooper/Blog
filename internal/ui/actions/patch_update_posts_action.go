package actions

import (
	"blog/assets/locales"
	"blog/internal/domain/posts"
	"blog/internal/platform/dtos"
	"blog/internal/platform/pkg/lang"
	"blog/internal/platform/pkg/validators"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// PatchUpdatePostAction handles HTTP PATCH requests to update a post, validating input and ensuring user authorization.
func PatchUpdatePostAction(interactor posts.Interactor, validator validators.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetString("locale")
		authorID := ctx.GetString("userID")
		postID := ctx.Param("postID")
		_, err := uuid.Parse(postID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.ResourceNotFound),
			})
			return
		}
		var req dtos.UpdatePostRequest
		if err = ctx.ShouldBindJSON(&req); err != nil {
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
		post, err := interactor.GetPost(ctx, postID)
		if err != nil {
			if errors.Is(err, posts.ErrPostNotFound) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"message": lang.TryBy(ctx.GetString("locale"), locales.ResourceNotFound),
				})
				return
			}
		}
		if authorID != post.Author {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.ForbiddenError),
			})
			return
		}
		err = interactor.UpdatePost(ctx, postID, req)
		if err != nil {
			log.WithError(err).Error("error saving post")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InternalServerError),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": lang.TryBy(ctx.GetString("locale"), locales.PostUpdated),
		})
	}
}
