package posts

import (
	"blog/internal/platform/dtos"
	"context"
)

type Interactor interface {
	Save(ctx context.Context, req dtos.CreatePostRequest) error
}
