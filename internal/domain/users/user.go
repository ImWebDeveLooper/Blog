package users

import "errors"

var (
	ErrUserNotValid   = errors.New("user is not valid")
	ErrUserNotAllowed = errors.New("user is not allowed")
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Username  string
	Email     string
	Password  string
	Bio       string
	Avatar    string
}
