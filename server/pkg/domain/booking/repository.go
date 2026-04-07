package booking

import (
	"context"
	
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(*Booking) error
	Update(*Booking) error
	Delete(uint) error

	FindByID(uint) (*Booking, error)

	List(ctx context.Context, page, limit int) ([]Booking, int64, error)

	FindAll(
		status *PaymentStatus,
		ctx context.Context,
		page int,
		limit int,
	) ([]Booking, int64, error)

	FindByUser(
		ctx context.Context,
		userID uint,
		page int,
		limit int,
	) ([]Booking, int64, error)

	FindByRoomID(roomID uint) ([]Booking, error)

	FindOverlappingBookings(
		roomID uint,
		checkIn time.Time,
		checkOut time.Time,
	) ([]Booking, error)

	UpdatePaymentStatus(id uint, status PaymentStatus) error

	UpdatePaymentStatusTx(
		tx *gorm.DB,
		id uint,
		status PaymentStatus,
	) error

	LockRoom(roomID uint) error

	WithTransaction(fn func(tx Repository) error) error

	

	FindExpiredBookings(now time.Time) ([]Booking, error)
}