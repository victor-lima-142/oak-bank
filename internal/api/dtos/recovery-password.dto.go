package dtos

type ForgotPasswordRequestDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponseDTO struct {
	Message string `json:"message"`
}

type ResetPasswordDTO struct {
	ResetToken      string `json:"reset_token" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
