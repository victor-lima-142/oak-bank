package dtos

type CreateAccountDTO struct {
	AccountTypeCode string  `json:"account_type_code" validate:"required"`
	AgencyNumber    string  `json:"agency_number" validate:"required"`
	InitialDeposit  float64 `json:"initial_deposit" validate:"gte=0"`
	OverdraftLimit  float64 `json:"overdraft_limit,omitempty"`
}
