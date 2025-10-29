package dtos

import "time"

type AccountFullSyncResponseDTO struct {
	Account  AccountSummaryDTO  `json:"account"`
	Customer CustomerProfileDTO `json:"customer"`
	Address  *AddressDTO        `json:"address,omitempty"`
	LastSync time.Time          `json:"last_sync"`
}

type AccountSummaryDTO struct {
	AccountID        string  `json:"account_id"`
	AccountNumber    string  `json:"account_number"`
	AgencyNumber     string  `json:"agency_number"`
	AccountTypeCode  string  `json:"account_type_code"`
	AccountStatus    string  `json:"account_status"`
	CurrentBalance   float64 `json:"current_balance"`
	AvailableBalance float64 `json:"available_balance"`
	OverdraftLimit   float64 `json:"overdraft_limit"`
}

type CustomerProfileDTO struct {
	CustomerID   string `json:"customer_id"`
	CustomerName string `json:"customer_name"`
	TaxID        string `json:"tax_id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	IsPep        bool   `json:"is_pep"`
	KYCStatus    string `json:"kyc_status"`
}

type AddressDTO struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	AddressType  string `json:"address_type"`
}
