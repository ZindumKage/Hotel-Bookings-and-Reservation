package audit_logs

import "time"

type AuditFilter struct {
	UserID     *uint
	Action     *string
	Entity     *string
	RiskLevel  *RiskLevel
	Suspicious *bool
	StartDate  *time.Time
	EndDate    *time.Time
}

type Repository interface {
	Save(log *AuditLog) error
	FindWithFilter(filter AuditFilter, page, limit int) ([]AuditLog, int64, error)
}