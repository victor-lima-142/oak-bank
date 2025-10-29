package dtos

type AccountOnboardingDTO struct {
	User     *UserRegistrationDTO `json:"user,omitempty"`
	Customer *CreateCustomerDTO   `json:"customer,omitempty"`
	Address  *CreateAddressDTO    `json:"address,omitempty"`
	Account  *CreateAccountDTO    `json:"account,omitempty"`
	KYC      *CustomerKycDTO      `json:"kyc,omitempty"`
}

type AccountOnboardingResponseDTO struct {
	UserID      string `json:"user_id"`
	CustomerID  string `json:"customer_id"`
	AccountID   string `json:"account_id"`
	KYCStatus   string `json:"kyc_status"`
	AccountType string `json:"account_type"`
	AccountNum  string `json:"account_number"`
}
