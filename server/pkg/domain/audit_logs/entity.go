package audit

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

	UserID    uint
	UserEmail string
	UserRole  string

	Action   string
	Entity   string
	EntityID string

	BeforeState   json.RawMessage
	AfterState    json.RawMessage
	ChangedFields json.RawMessage

	RiskLevel  RiskLevel
	Suspicious bool
	Reason     string

	IPAddress string
	UserAgent string

	CreatedAt time.Time
}