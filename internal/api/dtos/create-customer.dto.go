package dtos

import "time"

type CreateCustomerDTO struct {
	TaxID         string    `json:"tax_id" validate:"required,len=11"`
	CustomerName  string    `json:"customer_name" validate:"required"`
	DateOfBirth   time.Time `json:"date_of_birth" validate:"required"`
	CustomerPhone string    `json:"customer_phone,omitempty"`
	CustomerEmail string    `json:"customer_email,omitempty"`
	IsPep         bool      `json:"is_pep"`
}

type CreateAddressDTO struct {
	AddressLine1 string `json:"address_line_1" validate:"required"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	Country      string `json:"country" validate:"required"`
	PostalCode   string `json:"postal_code" validate:"required"`
	AddressType  string `json:"address_type" validate:"required,oneof=residential commercial billing"`
}

type CustomerKycDTO struct {
	KYCStatus    string `json:"kyc_status" validate:"required,oneof=PENDING APPROVED REJECTED"`
	RiskRating   string `json:"risk_rating,omitempty"`
	OtherDetails any    `json:"other_details,omitempty"`
}
