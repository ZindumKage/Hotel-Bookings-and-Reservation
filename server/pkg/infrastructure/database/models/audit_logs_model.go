package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID uint `gorm:"primaryKey"`

	UserID    uint
	UserEmail string
	UserRole  string

	Action   string
	Entity   string
	EntityID string

	BeforeState   datatypes.JSON
	AfterState    datatypes.JSON
	ChangedFields datatypes.JSON

	RiskLevel  string
	Suspicious bool
	Reason     string

	IPAddress string
	UserAgent string

	CreatedAt time.Time
}