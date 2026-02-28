package models

import (
	"time"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model

	UserID uint `gorm:"not null"`
	RoomID uint `gorm:"not null"`

	Rating  int    `gorm:"not null"`
	Comment string

	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
	Room Room `gorm:"foreignKey:RoomID"`
}