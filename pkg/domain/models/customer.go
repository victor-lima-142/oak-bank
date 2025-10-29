package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ===========================
// CUSTOMERS
// ===========================

type Customer struct {
	CustomerID     string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"customer_id"`
	TaxID          string         `gorm:"type:varchar(11);uniqueIndex:idx_customers_taxId;not null" json:"taxId"`
	TaxIDHash      string         `gorm:"type:varchar(64);uniqueIndex:idx_customers_taxId_hash;not null" json:"taxId_hash"`
	CustomerName   string         `gorm:"type:varchar(200);not null" json:"customer_name"`
	CustomerPhone  sql.NullString `gorm:"type:varchar(20)" json:"customer_phone"`
	CustomerEmail  sql.NullString `gorm:"type:varchar(100)" json:"customer_email"`
	DateOfBirth    time.Time      `gorm:"not null" json:"date_of_birth"`
	DateRegistered time.Time      `gorm:"autoCreateTime;not null" json:"date_registered"`
	IsPep          bool           `gorm:"default:false;not null" json:"is_pep"`
	KYCStatus      string         `gorm:"type:varchar(20);default:'PENDING';not null" json:"kyc_status"`
	KYCVerifiedAt  sql.NullTime   `json:"kyc_verified_at"`
	RiskRating     sql.NullString `gorm:"type:varchar(10)" json:"risk_rating"`
	OtherDetails   datatypes.JSON `gorm:"type:jsonb" json:"other_details"`

	// Relations
	Users             []User            `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"users,omitempty"`
	Account           *Account          `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"account,omitempty"`
	CustomerAddresses []CustomerAddress `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"customer_addresses,omitempty"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.CustomerID == "" {
		c.CustomerID = uuid.New().String()
	}
	return nil
}

func (Customer) TableName() string {
	return "customers"
}

type CustomerAddress struct {
	CustomerID  string    `gorm:"type:uuid;primaryKey" json:"customer_id"`
	AddressID   string    `gorm:"type:uuid;primaryKey" json:"address_id"`
	AddressType string    `gorm:"type:varchar(20);not null" json:"address_type"`
	CreatedAt   time.Time `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	Customer       *Customer       `gorm:"foreignKey:CustomerID;references:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	Address        *Address        `gorm:"foreignKey:AddressID;references:AddressID;constraint:OnDelete:CASCADE" json:"address,omitempty"`
	RefAddressType *RefAddressType `gorm:"foreignKey:AddressType;references:AddressTypeCode" json:"ref_address_type,omitempty"`
}

func (CustomerAddress) TableName() string {
	return "customer_addresses"
}
