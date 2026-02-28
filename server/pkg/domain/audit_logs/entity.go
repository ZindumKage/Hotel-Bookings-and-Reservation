package audit_logs

import (
	"encoding/json"
	"time"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "LOW"
	RiskMedium   RiskLevel = "MEDIUM"
	RiskHigh     RiskLevel = "HIGH"
	RiskCritical RiskLevel = "CRITICAL"
)

type AuditLog struct {
	ID uint

	// Actor
	UserID    uint
	UserEmail string
	UserRole  string

	// Action
	Action   string
	Entity   string
	EntityID string

	// State changes
	BeforeState   json.RawMessage
	AfterState    json.RawMessage
	ChangedFields json.RawMessage

	// Risk
	RiskLevel  RiskLevel
	Suspicious bool
	Reason     string

	// Metadata
	CreatedAt time.Time
}