package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===========================
// ADDRESSES
// ===========================

type Address struct {
	AddressID    string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"address_id"`
	AddressLine1 string         `gorm:"type:varchar(200);not null" json:"address_line_1"`
	AddressLine2 sql.NullString `gorm:"type:varchar(200)" json:"address_line_2"`
	City         string         `gorm:"type:varchar(100);not null" json:"city"`
	State        string         `gorm:"type:varchar(50);not null" json:"state"`
	Country      string         `gorm:"type:varchar(50);not null" json:"country"`
	PostalCode   string         `gorm:"type:varchar(20);not null" json:"postal_code"`
	CreatedAt    time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`

	// Relations
	CustomerAddresses []CustomerAddress `gorm:"foreignKey:AddressID;constraint:OnDelete:CASCADE" json:"customer_addresses,omitempty"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.AddressID == "" {
		a.AddressID = uuid.New().String()
	}
	return nil
}

func (Address) TableName() string {
	return "addresses"
}

type RefAddressType struct {
	AddressTypeCode string `gorm:"type:varchar(20);primaryKey" json:"address_type_code"`
	Description     string `gorm:"type:varchar(100);not null" json:"description"`

	// Relations
	CustomerAddresses []CustomerAddress `gorm:"foreignKey:AddressType" json:"customer_addresses,omitempty"`
}

func (RefAddressType) TableName() string {
	return "ref_address_types"
}
