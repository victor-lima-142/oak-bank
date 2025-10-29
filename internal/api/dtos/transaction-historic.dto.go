package dtos

import "time"

type TransactionHistoryRequestDTO struct {
	AccountID string     `json:"account_id" validate:"required,uuid4"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Status    string     `json:"status,omitempty" validate:"omitempty,oneof=PENDING COMPLETED FAILED CANCELLED"`
	TypeCode  string     `json:"type_code,omitempty"`
	Limit     int        `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset    int        `json:"offset,omitempty" validate:"omitempty,min=0"`
}

type TransactionHistoryResponseDTO struct {
	AccountID    string                   `json:"account_id"`
	Transactions []TransactionListItemDTO `json:"transactions"`
	TotalCount   int                      `json:"total_count"`
	Limit        int                      `json:"limit"`
	Offset       int                      `json:"offset"`
}

type TransactionListItemDTO struct {
	TransactionID     string    `json:"transaction_id"`
	TypeCode          string    `json:"type_code"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	Date              time.Time `json:"date"`
	Description       string    `json:"description,omitempty"`
	AccountDestNumber string    `json:"account_dest_number,omitempty"`
	AccountDestName   string    `json:"account_dest_name,omitempty"`
}
