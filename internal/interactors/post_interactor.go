package interactors

import (
	"blog/internal/domain/posts"
	"blog/internal/platform/dtos"
	"context"
	"github.com/google/uuid"
)

type PostInteractor struct {
	postRepository posts.Repository
}

func NewPostInteractor(r posts.Repository) posts.Interactor {
	return &PostInteractor{
		postRepository: r,
	}
}
func (i *PostInteractor) Save(ctx context.Context, req dtos.CreatePostRequest) error {
	postInfo := posts.Post{
		ID:      uuid.New().String(),
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}
	err := i.postRepository.Save(ctx, postInfo)
	if err != nil {
		return err
	}
	return nil
}

//func (i *PostInteractor) Update(ctx context.Context, req dtos.UpdatePostRequest) error {
//	postInfo := posts.Post{
//		Title:   req.Title,
//		Content: req.Content,
//		Author:  req.Author,
//	}
//	err := i.postRepository.Update(ctx, postInfo)
//	if err != nil {
//		return err
//	}
//	return nil
//}
