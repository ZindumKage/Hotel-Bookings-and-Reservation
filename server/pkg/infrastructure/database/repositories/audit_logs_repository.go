package repositories

import (
	

	"gorm.io/gorm"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Save(log *audit_logs.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditRepository) FindWithFilter(filter audit_logs.AuditFilter, page, limit int) ([]audit_logs.AuditLog, int64, error) {
	var logs []audit_logs.AuditLog
	var total int64

	query := r.db.Model(&audit_logs.AuditLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	if filter.Entity != nil {
		query = query.Where("entity = ?", *filter.Entity)
	}
	if filter.RiskLevel != nil {
		query = query.Where("risk_level = ?", *filter.RiskLevel)
	}
	if filter.Suspicious != nil {
		query = query.Where("suspicious = ?", *filter.Suspicious)
	}
	if filter.StartDate != nil && filter.EndDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *filter.StartDate, *filter.EndDate)
	}

	query.Count(&total)

	offset := (page - 1) * limit
	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}