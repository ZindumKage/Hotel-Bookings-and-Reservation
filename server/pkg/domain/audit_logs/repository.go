package audit

import (
	"context"
	"time"
)

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
		Save(ctx context.Context, log *AuditLog) error

	FindWithFilter(
		ctx context.Context,
		filter AuditFilter,
		page int,
		limit int,
	) ([]AuditLog, int64, error)
}