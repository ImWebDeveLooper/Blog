package middlewares

import (
	"blog/assets/locales"
	"blog/internal/platform/pkg/lang"
	"blog/internal/platform/pkg/rbac"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Authorize(obj rbac.Object, act rbac.Action, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sub, existed := ctx.Get("userID")
		if !existed {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.NotLoggedIn),
			})
			return
		}
		err := enforcer.LoadPolicy()
		if err != nil {
			log.WithError(err).Error("failed to load policy from database")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InternalServerError),
			})
			return
		}
		ok, err := enforcer.Enforce(fmt.Sprint(sub), obj.String(), act.String())
		if err != nil {
			log.WithError(err).Error("failed to authorize user")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.InternalServerError),
			})
			return
		}
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": lang.TryBy(ctx.GetString("locale"), locales.ForbiddenError),
			})
			return
		}
		ctx.Next()
	}
}
