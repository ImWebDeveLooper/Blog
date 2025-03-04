package actions

import (
	"blog/internal/domain/posts"
	"blog/internal/platform/dtos"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func PostCreatePostsAction(interactor posts.Interactor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.CreatePostRequest
		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		err = interactor.Save(ctx, req)
		if err != nil {
			log.WithError(err).Error("error while saving post")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "Post has been created successfully",
		})
	}
}
