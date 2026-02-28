package models

import (
	"time"
	"gorm.io/gorm"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
)

type Payment struct {
	gorm.Model

	BookingID uint `gorm:"uniqueIndex;not null"`

	Amount int64 `gorm:"not null"` // store in kobo

	Status common.PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`

	TransactionRef string `gorm:"uniqueIndex"`

	PaymentDate time.Time

	Booking common.BookingStatus `gorm:"foreignKey:BookingID"`
}