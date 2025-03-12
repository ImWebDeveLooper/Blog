package interactors

import (
	"blog/internal/domain/users"
	"blog/internal/platform/dtos"
	"blog/internal/platform/pkg/jwt"
	"context"
	"github.com/google/uuid"
)

type UserInteractor struct {
	userRepository users.Repository
	passwordHasher users.PasswordHasher
	jwtService     jwt.Service
}

func NewUserInteractor(r users.Repository, p users.PasswordHasher, jwt jwt.Service) users.Interactor {
	return &UserInteractor{
		userRepository: r,
		passwordHasher: p,
		jwtService:     jwt,
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

func (i *UserInteractor) Login(ctx context.Context, identifier, password string) (string, error) {
	user, err := i.userRepository.FindByEmailOrUsername(ctx, identifier)
	if err != nil {
		return "", err
	}

	storedPass := user.Password
	if storedPass == "" {
		return "", users.ErrUserNotValid
	}
	isValid := i.passwordHasher.ValidatePassword(storedPass, password)
	if !isValid {
		return "", users.ErrUserNotValid
	}

	token, err := i.jwtService.GenerateToken(user.Email, user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
