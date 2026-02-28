package models

// import (
// 	"time"

// 	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
// 	"gorm.io/gorm"
// )

// type Review struct {
// 	gorm.Model

// 	UserID uint `gorm:"not null"`
// 	RoomID uint `gorm:"not null"`

// 	Rating  int    `gorm:"not null"`
// 	Comment string

// 	CreatedAt time.Time

// 	User common.User `gorm:"foreignKey:UserID"`
// 	Room common.Room`gorm:"foreignKey:RoomID"`
// }