package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model

	ID        uint   `gorm:"primaryKey"`
	Reference string `gorm:"uniqueIndex;not null"`

	UserID     uint   `gorm:"not null"`
	RoomID     uint   `gorm:"index:idx_room;not null"`
	RoomNumber string `gorm:"not null"`

	CheckInDate  time.Time `gorm:"index:idx_room_dates;not null"`
	CheckOutDate time.Time `gorm:"index:idx_room_dates;not null"`

	NightCount int    `gorm:"not null"`
	UnitPrice  int64  `gorm:"not null"`
	TaxAmount  int64  `gorm:"not null"`
	Discount   int64  `gorm:"not null"`
	TotalPrice int64  `gorm:"not null"`
	Currency   string `gorm:"type:varchar(10);not null"`

	Status        string `gorm:"type:varchar(50);index;not null"`
	PaymentStatus string `gorm:"type:varchar(50);index;not null"`

	ExpiresAt   time.Time `gorm:"not null"`
	CancelledAt *time.Time

	
}
