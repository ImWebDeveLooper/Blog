package application

import (
	"blog/internal/ui/actions"
)

func (a *App) RegisterRoutes() {
	v1 := a.Router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/signup", actions.PostSignUpAction(a.Interactors.UserInteractor))
			auth.POST("/login", actions.PostLoginAction(a.Interactors.UserInteractor))
		}

		post := v1.Group("/post")
		{
			post.POST("/posts/create", actions.PostCreatePostsAction(a.Interactors.PostInteractor))
			post.PATCH("/posts/:author")
		}
	}
}
