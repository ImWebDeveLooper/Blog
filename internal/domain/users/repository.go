package users

import "context"

type Repository interface {
	Save(ctx context.Context, user User) error
}
