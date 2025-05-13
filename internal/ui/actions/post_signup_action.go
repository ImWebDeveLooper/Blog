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

// PostSignUpAction handles HTTP POST requests for signing up the users and validate the inputs.
func PostSignUpAction(interactor users.Interactor, validator validators.Validator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := ctx.GetString("locale")
		var req dtos.CreateUserRequest
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
		err := interactor.SignUp(ctx, req)
		if err != nil {
			log.WithError(err).Error("error while saving user")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InternalServerError),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": lang.TryBy(ctx.GetString("locale"), locales.UserRegistered),
		})
	}
}
