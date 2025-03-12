package actions

import (
	"blog/internal/domain/users"
	"blog/internal/platform/dtos"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func PostLoginAction(interactor users.Interactor) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}
		token, err := interactor.Login(ctx, req.Identifier, req.Password)
		if err != nil {
			log.WithError(err).Error("invalid credentials")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Credentials"})
			return
		}
		ctx.JSON(200, gin.H{
			"token": struct {
				AccessToken string `json:"access_token"`
			}{AccessToken: token},
			"message": "The User Successfully Logged In",
		})
	}
}
