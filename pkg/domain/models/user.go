package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ===========================
// USERS & AUTHENTICATION
// ===========================

type User struct {
	UserID              string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"user_id"`
	CustomerID          sql.NullString `gorm:"type:uuid;index:idx_users_customer_id" json:"customer_id"`
	Username            string         `gorm:"type:varchar(50);uniqueIndex:idx_users_username;not null" json:"username"`
	Email               string         `gorm:"type:varchar(100);uniqueIndex:idx_users_email;not null" json:"email"`
	PasswordHash        string         `gorm:"type:varchar(255);not null" json:"password_hash"`
	UserRole            string         `gorm:"type:varchar(20);not null" json:"user_role"`
	IsActive            bool           `gorm:"default:true;not null" json:"is_active"`
	IsLocked            bool           `gorm:"default:false;not null" json:"is_locked"`
	FailedLoginAttempts int            `gorm:"default:0;not null" json:"failed_login_attempts"`
	LastLogin           sql.NullTime   `json:"last_login"`
	CreatedAt           time.Time      `gorm:"autoCreateTime;not null" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime;not null" json:"updated_at"`
	TwoFactorEnabled    bool           `gorm:"default:false;not null" json:"two_factor_enabled"`
	TwoFactorSecret     sql.NullString `gorm:"type:varchar(255)" json:"two_factor_secret"`

	// Relations
	Customer             *Customer              `gorm:"foreignKey:CustomerID;references:CustomerID;constraint:OnDelete:CASCADE" json:"customer,omitempty"`
	AuthLogs             []UserAuthLog          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"auth_logs,omitempty"`
	TransactionsCreated  []Transaction          `gorm:"foreignKey:CreatedByUserID;constraint:OnDelete:SET NULL" json:"transactions_created,omitempty"`
	AccountStatusChanges []AccountStatusHistory `gorm:"foreignKey:ChangedByUserID;constraint:OnDelete:SET NULL" json:"account_status_changes,omitempty"`
	AuditLogs            []AuditLog             `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"audit_logs,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == "" {
		u.UserID = uuid.New().String()
	}
	return nil
}

func (User) TableName() string {
	return "users"
}

type UserAuthLog struct {
	LogID             string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"log_id"`
	UserID            sql.NullString `gorm:"type:uuid;index:idx_auth_user_id" json:"user_id"`
	AuthEventType     string         `gorm:"type:varchar(30);not null" json:"auth_event_type"`
	AuthTimestamp     time.Time      `gorm:"autoCreateTime;index:idx_auth_timestamp;not null" json:"auth_timestamp"`
	IPAddress         sql.NullString `gorm:"type:varchar(45);index:idx_auth_ip" json:"ip_address"`
	UserAgent         sql.NullString `gorm:"type:varchar(500)" json:"user_agent"`
	DeviceFingerprint sql.NullString `gorm:"type:varchar(255)" json:"device_fingerprint"`
	Success           bool           `gorm:"not null" json:"success"`
	FailureReason     sql.NullString `gorm:"type:varchar(200)" json:"failure_reason"`
	SessionID         sql.NullString `gorm:"type:varchar(255)" json:"session_id"`
	Geolocation       sql.NullString `gorm:"type:varchar(100)" json:"geolocation"`
	RiskScore         sql.NullInt32  `json:"risk_score"`

	// Relations
	User *User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

func (ual *UserAuthLog) BeforeCreate(tx *gorm.DB) error {
	if ual.LogID == "" {
		ual.LogID = uuid.New().String()
	}
	return nil
}

func (UserAuthLog) TableName() string {
	return "user_auth_log"
}
