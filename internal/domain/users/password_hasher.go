package users

type PasswordHasher interface {
	Hash(password string) (string, error)
	ValidatePassword(hashedPassword, password string) bool
}
