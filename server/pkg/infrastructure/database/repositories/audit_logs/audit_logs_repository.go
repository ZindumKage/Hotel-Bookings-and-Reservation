package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
)

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Save(ctx context.Context, log *audit.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *AuditRepository) FindWithFilter(
	ctx context.Context,
	filter audit.AuditFilter,
	page int,
	limit int,
) ([]audit.AuditLog, int64, error) {

	var modelsLogs []models.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&models.AuditLog{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&modelsLogs).
		Error

	if err != nil {
		return nil, 0, err
	}

	// Convert models -> domain
	var logs []audit.AuditLog
	for _, m := range modelsLogs {
		logs = append(logs, audit.AuditLog{
			ID:        m.ID,
			UserID:    m.UserID,
			Action:    m.Action,
			Entity:    m.Entity,
			EntityID:  m.EntityID,
			CreatedAt: m.CreatedAt,
		})
	}

	return logs, total, nil
}