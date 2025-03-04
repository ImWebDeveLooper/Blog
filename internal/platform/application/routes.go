package application

import (
	"blog/internal/ui/actions"
)

func (a *App) RegisterRoutes() {
	v1 := a.Router.Group("/api/v1")
	{
		post := v1.Group("/post")
		{
			post.POST("/posts/create", actions.PostCreatePostsAction(a.Interactors.PostInteractor))
		}
	}
}
