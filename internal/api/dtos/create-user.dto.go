package dtos

type UserRegistrationDTO struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	UserRole string `json:"user_role" validate:"required,oneof=admin customer"`
}
