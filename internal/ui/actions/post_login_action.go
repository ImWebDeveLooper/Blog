package actions

import (
	"blog/assets/locales"
	"blog/internal/domain/users"
	"blog/internal/platform/dtos"
	"blog/internal/platform/pkg/lang"
	"blog/internal/platform/pkg/validators"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func PostLoginAction(interactor users.Interactor, validator validators.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetString("locale")
		var req dtos.LoginRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": lang.TryBy(ctx.GetString("locale"), locales.InvalidSchemaError),
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
		token, err := interactor.Login(ctx, req.Identifier, req.Password)
		if err != nil {
			log.WithError(err).Error("invalid credentials")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InvalidCredential),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": lang.TryBy(ctx.GetString("locale"), locales.LoggedIn),
			"data": struct {
				AccessToken string `json:"access_token"`
			}{AccessToken: token},
		})
	}
}
