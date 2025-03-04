package actions

import (
	"blog/internal/domain/users"
	"blog/internal/platform/dtos"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func PostSignUpAction(interactor users.Interactor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.CreateUserRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		err := interactor.SignUp(ctx, req)
		if err != nil {
			log.WithError(err).Error("error while saving user")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		ctx.JSON(200, gin.H{
			"message": "The User Registered",
		})
	}
}
