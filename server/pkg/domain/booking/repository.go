package booking

import "time"

type Repository interface {
	Create(booking *Booking) error
	Update(booking *Booking) error
	Delete(id uint) error
	UpdatePaymentStatus(id uint, status PaymentStatus) error
	List(page, limit int) ([]Booking, int64, error)
	FindByUser(userID uint) ([]Booking, error)
	FindByID(id uint) (*Booking, error)
	FindByRoomID(roomID uint) ([]Booking, error)
	FindOverlappingBookings(roomID uint, checkIn, checkOut time.Time) ([]Booking, error)
	FindAll(status *PaymentStatus, page, limit int) ([]Booking, int64, error)
	WithTransaction(fn func(tx Repository) error) error
}
