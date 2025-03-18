package application

import (
	"blog/internal/ui/actions"
)

func (a *App) RegisterRoutes() {
	v1 := a.Router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", actions.PostSignUpAction(a.Interactors.UserInteractor, *a.Validator))
			auth.POST("/login", actions.PostLoginAction(a.Interactors.UserInteractor, *a.Validator))
		}

		post := v1.Group("/posts")
		{
			post.POST("/create", actions.PostCreatePostsAction(a.Interactors.PostInteractor, *a.Validator))
			post.PATCH("/:author")
		}
	}
}
