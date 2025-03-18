package dtos

type CreateUserRequest struct {
	FirstName string `json:"firstName" validate:"required,gt=2"`
	LastName  string `json:"lastName" validate:"required,gt=2"`
	Username  string `json:"username" validate:"required,gt=3,username,uniqueUsername"`
	Email     string `json:"email" validate:"required,email,uniqueEmail"`
	Password  string `json:"password" validate:"required,min=8,strongPassword"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,emailOrUsername"`
	Password   string `json:"password" validate:"required"`
}
