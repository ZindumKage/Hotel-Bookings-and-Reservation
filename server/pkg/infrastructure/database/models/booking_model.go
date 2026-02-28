 package models

// import (
// 	"time"
// 	"gorm.io/gorm"
// )

// type Booking struct {
// 	gorm.Model

// 	UserID uint `gorm:"not null;index"`
// 	RoomID uint `gorm:"not null;index"`

// 	CheckInDate  time.Time `gorm:"not null"`
// 	CheckOutDate time.Time `gorm:"not null"`

// 	Status        BookingStatus `gorm:"type:varchar(20);default:'PENDING'"`
// 	PaymentStatus PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`

// 	TotalPrice int64 `gorm:"not null"`

// 	User User `gorm:"foreignKey:UserID"`
// 	Room Room `gorm:"foreignKey:RoomID"`

// 	Payment Payment
// }