package dtos

import "time"

type TransactionRequestDTO struct {
	TransactionTypeCode string  `json:"transaction_type_code" validate:"required,oneof=PIX TED TRANSFER INTERNAL"`
	AccountIDOrigin     string  `json:"account_id_origin" validate:"required,uuid4"`
	AccountIDDest       string  `json:"account_id_dest,omitempty" validate:"omitempty,uuid4"`
	Amount              float64 `json:"amount" validate:"required,gt=0"`
	Description         string  `json:"description,omitempty" validate:"max=500"`
	IdempotencyKey      string  `json:"idempotency_key" validate:"required,max=100"`
	ExternalReference   string  `json:"external_reference,omitempty"`
	Metadata            any     `json:"metadata,omitempty"`
}

type TransactionResponseDTO struct {
	TransactionID     string          `json:"transaction_id"`
	TransactionType   string          `json:"transaction_type"`
	TransactionStatus string          `json:"transaction_status"`
	TransactionDate   time.Time       `json:"transaction_date"`
	CompletedAt       *time.Time      `json:"completed_at,omitempty"`
	Amount            float64         `json:"amount"`
	Description       string          `json:"description,omitempty"`
	AccountOrigin     AccountMiniDTO  `json:"account_origin"`
	AccountDest       *AccountMiniDTO `json:"account_dest,omitempty"`
	BalanceAfter      *float64        `json:"balance_after,omitempty"`
}

type AccountMiniDTO struct {
	AccountID     string `json:"account_id"`
	AccountNumber string `json:"account_number"`
	AgencyNumber  string `json:"agency_number"`
	AccountType   string `json:"account_type"`
}
