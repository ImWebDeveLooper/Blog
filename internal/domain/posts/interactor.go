package posts

import (
	"blog/internal/platform/dtos"
	"context"
)

type Interactor interface {
	CreatePost(ctx context.Context, author string, req dtos.CreatePostRequest) error
	GetPost(ctx context.Context, postID string) (*Post, error)
	UpdatePost(ctx context.Context, postID string, req dtos.UpdatePostRequest) error
}
