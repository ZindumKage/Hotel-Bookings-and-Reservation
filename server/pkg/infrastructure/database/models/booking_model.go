package models

import (
	"time"

	"gorm.io/gorm"
)



type Booking struct {
	gorm.Model

	ID        uint   `gorm:"primaryKey"`
	Reference string `gorm:"uniqueIndex;not null"`

	UserID uint `gorm:"index:idx_user_bookings;not null"`
	// Composite index for availability search
	RoomID uint `gorm:"index:idx_room_dates,priority:1;not null"`

	RoomNumber string `gorm:"not null"`

	CheckInDate time.Time `gorm:"index:idx_room_dates,priority:2;not null"`

	CheckOutDate time.Time `gorm:"index:idx_room_dates,priority:3;not null"`

	NightCount int   `gorm:"not null"`
	UnitPrice  int64 `gorm:"not null"`
	TaxAmount  int64 `gorm:"not null"`
	Discount   int64 `gorm:"not null"`
	TotalPrice int64 `gorm:"not null"`

	Currency string `gorm:"type:varchar(10);not null"`

	Status        string `gorm:"type:varchar(50);index;not null"`
	PaymentStatus string `gorm:"type:varchar(50);index;not null"`

	ExpiresAt   time.Time `gorm:"not null;index"`
	CancelledAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}