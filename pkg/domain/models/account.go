package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ===========================
// ACCOUNTS
// ===========================

type Account struct {
	AccountID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"account_id"`
	CustomerID       string         `gorm:"type:uuid;uniqueIndex;index:idx_accounts_customer_id;not null" json:"customer_id"`
	AccountTypeCode  string         `gorm:"type:varchar(20);not null" json:"account_type_code"`
	AccountNumber    string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"account_number"`
	AgencyNumber     string         `gorm:"type:varchar(10);not null" json:"agency_number"`
	CurrentBalance   float64        `gorm:"type:decimal(15,2);not null" json:"current_balance"`
	AvailableBalance float64        `gorm:"type:decimal(15,2);not null" json:"available_balance"`
	OverdraftLimit   float64        `gorm:"type:decimal(15,2);default:0;not null" json:"overdraft_limit"`
	DateOpened       time.Time      `gorm:"type:date;not null" json:"date_opened"`
	DateClosed       sql.NullTime   `json:"date_closed"`
	AccountStatus    string         `gorm:"type:varchar(20);default:'ACTIVE';index:idx_accounts_status;not null" json:"account_status"`
	BlockedReason    sql.NullString `gorm:"type:varchar(200)" json:"blocked_reason"`
	CreatedAt        time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	Customer           *Customer              `gorm:"foreignKey:CustomerID;references:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	RefAccountType     *RefAccountType        `gorm:"foreignKey:AccountTypeCode;references:AccountTypeCode" json:"ref_account_type,omitempty"`
	StatusHistory      []AccountStatusHistory `gorm:"foreignKey:AccountID;constraint:OnDelete:RESTRICT" json:"status_history,omitempty"`
	TransactionsOrigin []Transaction          `gorm:"foreignKey:AccountIDOrigin;constraint:OnDelete:RESTRICT" json:"transactions_origin,omitempty"`
	TransactionsDest   []Transaction          `gorm:"foreignKey:AccountIDDest;constraint:OnDelete:RESTRICT" json:"transactions_dest,omitempty"`
}

func (acc *Account) BeforeCreate(tx *gorm.DB) error {
	if acc.AccountID == "" {
		acc.AccountID = uuid.New().String()
	}
	return nil
}

func (Account) TableName() string {
	return "accounts"
}

type AccountStatusHistory struct {
	HistoryID       string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"history_id"`
	AccountID       string         `gorm:"type:uuid;index:idx_status_account_id;not null" json:"account_id"`
	PreviousStatus  sql.NullString `gorm:"type:varchar(20)" json:"previous_status"`
	NewStatus       string         `gorm:"type:varchar(20);not null" json:"new_status"`
	ChangeReason    sql.NullString `gorm:"type:varchar(500)" json:"change_reason"`
	ChangedByUserID sql.NullString `gorm:"type:uuid" json:"changed_by_user_id"`
	ChangedAt       time.Time      `gorm:"autoCreateTime;not null" json:"changed_at"`
	IPAddress       sql.NullString `gorm:"type:varchar(45)" json:"ip_address"`
	AdditionalInfo  datatypes.JSON `gorm:"type:jsonb" json:"additional_info"`

	// Relations
	Account       *Account `gorm:"foreignKey:AccountID;references:AccountID;constraint:OnDelete:RESTRICT" json:"account,omitempty"`
	ChangedByUser *User    `gorm:"foreignKey:ChangedByUserID;references:UserID;constraint:OnDelete:SET NULL" json:"changed_by_user,omitempty"`
}

func (ash *AccountStatusHistory) BeforeCreate(tx *gorm.DB) error {
	if ash.HistoryID == "" {
		ash.HistoryID = uuid.New().String()
	}
	return nil
}

func (AccountStatusHistory) TableName() string {
	return "account_status_history"
}

type RefAccountType struct {
	AccountTypeCode string          `gorm:"type:varchar(20);primaryKey" json:"account_type_code"`
	Description     string          `gorm:"type:varchar(100);not null" json:"description"`
	AllowsOverdraft bool            `gorm:"default:false;not null" json:"allows_overdraft"`
	MonthlyFee      sql.NullFloat64 `gorm:"type:decimal(10,2)" json:"monthly_fee"`

	// Relations
	Accounts []Account `gorm:"foreignKey:AccountTypeCode" json:"accounts,omitempty"`
}

func (RefAccountType) TableName() string {
	return "ref_account_types"
}
