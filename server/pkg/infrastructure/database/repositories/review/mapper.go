package repositories

import (
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/review"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
)

func toModel(d *domain.Review) *models.Review {
	return &models.Review{
		Model:   gorm.Model{ID: d.ID},
		UserID:  d.UserID,
		RoomID:  d.RoomID,
		Rating:  d.Rating,
		Comment: d.Comment,
	}
}

func toDomain(m *models.Review) *domain.Review {
	return &domain.Review{
		ID:        m.ID,
		UserID:    m.UserID,
		RoomID:    m.RoomID,
		Rating:    m.Rating,
		Comment:   m.Comment,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}