package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ===========================
// TRANSACTIONS
// ===========================

type Transaction struct {
	TransactionID       string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"transaction_id"`
	AccountIDOrigin     string          `gorm:"type:uuid;index:idx_trans_account_origin,priority:1;index:idx_trans_account_date,priority:1;not null" json:"account_id_origin"`
	AccountIDDest       sql.NullString  `gorm:"type:uuid" json:"account_id_dest"`
	TransactionTypeCode string          `gorm:"type:varchar(20);not null" json:"transaction_type_code"`
	TransactionAmount   float64         `gorm:"type:decimal(15,2);not null" json:"transaction_amount"`
	TransactionStatus   string          `gorm:"type:varchar(20);default:'PENDING';index:idx_trans_status;not null" json:"transaction_status"`
	TransactionDate     time.Time       `gorm:"autoCreateTime;index:idx_trans_date;index:idx_trans_account_date,priority:2;not null" json:"transaction_date"`
	CompletedAt         sql.NullTime    `json:"completed_at"`
	Description         sql.NullString  `gorm:"type:varchar(500)" json:"description"`
	BalanceAfter        sql.NullFloat64 `gorm:"type:decimal(15,2)" json:"balance_after"`
	CreatedByUserID     sql.NullString  `gorm:"type:uuid" json:"created_by_user_id"`
	CreatedBySource     sql.NullString  `gorm:"type:varchar(20)" json:"created_by_source"`
	IdempotencyKey      string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"idempotency_key"`
	ExternalReference   sql.NullString  `gorm:"type:varchar(100)" json:"external_reference"`
	Metadata            datatypes.JSON  `gorm:"type:jsonb" json:"metadata"`

	// Relations
	AccountOrigin      *Account            `gorm:"foreignKey:AccountIDOrigin;references:AccountID;constraint:OnDelete:RESTRICT" json:"account_origin,omitempty"`
	AccountDest        *Account            `gorm:"foreignKey:AccountIDDest;references:AccountID;constraint:OnDelete:RESTRICT" json:"account_dest,omitempty"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByUserID;references:UserID;constraint:OnDelete:SET NULL" json:"created_by_user,omitempty"`
	RefTransactionType *RefTransactionType `gorm:"foreignKey:TransactionTypeCode;references:TransactionTypeCode" json:"ref_transaction_type,omitempty"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.TransactionID == "" {
		t.TransactionID = uuid.New().String()
	}
	return nil
}

func (Transaction) TableName() string {
	return "transactions"
}

type RefTransactionType struct {
	TransactionTypeCode string          `gorm:"type:varchar(20);primaryKey" json:"transaction_type_code"`
	Description         string          `gorm:"type:varchar(100);not null" json:"description"`
	RequiresDestination bool            `gorm:"default:false;not null" json:"requires_destination"`
	MaxDailyAmount      sql.NullFloat64 `gorm:"type:decimal(15,2)" json:"max_daily_amount"`

	// Relations
	Transactions []Transaction `gorm:"foreignKey:TransactionTypeCode" json:"transactions,omitempty"`
}

func (RefTransactionType) TableName() string {
	return "ref_transaction_types"
}
