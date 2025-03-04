package interactors

import (
	"blog/internal/domain/users"
	"blog/internal/platform/dtos"
	"context"
	"github.com/google/uuid"
)

type UserInteractor struct {
	userRepository users.Repository
	passwordHasher users.PasswordHasher
}

func NewUserInteractor(r users.Repository, p users.PasswordHasher) users.Interactor {
	return &UserInteractor{
		userRepository: r,
		passwordHasher: p,
	}
}

func (i *UserInteractor) SignUp(ctx context.Context, req dtos.CreateUserRequest) error {
	hashedPassword, err := i.passwordHasher.Hash(req.Password)
	if err != nil {
		return err
	}
	userInfo := users.User{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
	}
	err = i.userRepository.Save(ctx, userInfo)
	if err != nil {
		return err
	}
	return nil
}
