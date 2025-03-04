package users

import (
	"blog/internal/platform/dtos"
	"context"
)

type Interactor interface {
	SignUp(ctx context.Context, req dtos.CreateUserRequest) error
}
