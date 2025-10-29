package models

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
)

// ===========================
// AUDIT LOG
// ===========================

type AuditLog struct {
	AuditID            string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"audit_id"`
	TableNameReference string         `gorm:"type:varchar(100);index:idx_audit_table_record,priority:1;not null" json:"table_name"`
	RecordID           string         `gorm:"type:varchar(100);index:idx_audit_table_record,priority:2;not null" json:"record_id"`
	OperationType      string         `gorm:"type:varchar(20);not null" json:"operation_type"`
	UserID             sql.NullString `gorm:"type:uuid;index:idx_audit_user_time,priority:1" json:"user_id"`
	OperationTimestamp time.Time      `gorm:"autoCreateTime;index:idx_audit_timestamp;index:idx_audit_table_record,priority:3;index:idx_audit_user_time,priority:2;not null" json:"operation_timestamp"`
	IPAddress          sql.NullString `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent          sql.NullString `gorm:"type:varchar(500)" json:"user_agent"`
	OldValues          datatypes.JSON `gorm:"type:jsonb" json:"old_values"`
	NewValues          datatypes.JSON `gorm:"type:jsonb" json:"new_values"`
	OperationResult    sql.NullString `gorm:"type:varchar(20)" json:"operation_result"`
	ErrorMessage       sql.NullString `json:"error_message"`
	SessionID          sql.NullString `gorm:"type:varchar(255)" json:"session_id"`
	RequestID          sql.NullString `gorm:"type:varchar(100)" json:"request_id"`
	AdditionalContext  datatypes.JSON `gorm:"type:jsonb" json:"additional_context"`

	// Relations
	User *User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
