package dtos

type TwoFactorVerificationDTO struct {
	UserID string `json:"user_id" validate:"required,uuid4"`
	Code   string `json:"code" validate:"required,len=6"`
}

type TwoFactorVerificationSenderDTO struct {
	Email string `json:"email" validate:"required,email"`
}
