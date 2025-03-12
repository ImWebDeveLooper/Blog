package users

import "errors"

var (
	ErrUserNotValid = errors.New("user is not valid")
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
