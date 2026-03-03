package booking

import (
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/booking"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
)






func toDomain(m *models.Booking) *domain.Booking {
	return &domain.Booking{
		ID:            m.ID,
		Reference:     m.Reference,
		UserID:        m.UserID,
		RoomID:        m.RoomID,
		RoomNumber:    m.RoomNumber,
		CheckInDate:   m.CheckInDate,
		CheckOutDate:  m.CheckOutDate,
		NightCount:    m.NightCount,
		UnitPrice:     m.UnitPrice,
		TaxAmount:     m.TaxAmount,
		Discount:      m.Discount,
		TotalPrice:    m.TotalPrice,
		Currency:      m.Currency,
		Status:        domain.BookingStatus(m.Status),
		PaymentStatus: domain.PaymentStatus(m.PaymentStatus),
		ExpiresAt:     m.ExpiresAt,
		CancelledAt:   m.CancelledAt,
	}
}

func toModel(d *domain.Booking) *models.Booking {
	return &models.Booking{
		Model:         gorm.Model{ID: d.ID},
		Reference:     d.Reference,
		UserID:        d.UserID,
		RoomID:        d.RoomID,
		RoomNumber:    d.RoomNumber,
		CheckInDate:   d.CheckInDate,
		CheckOutDate:  d.CheckOutDate,
		NightCount:    d.NightCount,
		UnitPrice:     d.UnitPrice,
		TaxAmount:     d.TaxAmount,
		Discount:      d.Discount,
		TotalPrice:    d.TotalPrice,
		Currency:      d.Currency,
		Status:        string(d.Status),
		PaymentStatus: string(d.PaymentStatus),
		ExpiresAt:     d.ExpiresAt,
		CancelledAt:   d.CancelledAt,
	}
}