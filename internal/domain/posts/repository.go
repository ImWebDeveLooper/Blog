package posts

import "context"

type Repository interface {
	Save(ctx context.Context, post Post) error
}
