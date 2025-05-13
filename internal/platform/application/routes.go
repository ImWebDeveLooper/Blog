package application

import (
	"blog/internal/platform/pkg/rbac"
	"blog/internal/ui/actions"
	"blog/internal/ui/middlewares"
)

func (a *App) RegisterRoutes() {
	v1 := a.Router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", actions.PostSignUpAction(a.Interactors.UserInteractor, *a.Validator))
			auth.POST("/login", actions.PostLoginAction(a.Interactors.UserInteractor, *a.Validator))
		}

		post := v1.Group("/posts", middlewares.AuthMiddleware(a.AuthService))
		{
			post.POST("/create",
				middlewares.Authorize(rbac.ObjectPost, rbac.ActionWrite, a.Enforcer),
				actions.PostCreatePostsAction(a.Interactors.PostInteractor, *a.Validator),
			)
			post.PATCH("/:postID",
				middlewares.Authorize(rbac.ObjectPost, rbac.ActionModify, a.Enforcer),
				actions.PatchUpdatePostAction(a.Interactors.PostInteractor, *a.Validator),
			)
		}
	}
}
