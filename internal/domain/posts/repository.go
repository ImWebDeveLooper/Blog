package posts

import "context"

type Repository interface {
	Save(ctx context.Context, post Post) error
	FindByID(ctx context.Context, id string) (*Post, error)
	Update(ctx context.Context, id string, update map[string]interface{}) error
}
