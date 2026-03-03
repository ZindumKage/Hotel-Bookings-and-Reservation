package models

import (
	"time"

	
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model

	UserID uint `gorm:"not null"`
	RoomID uint `gorm:"not null"`

	Rating  int    `gorm:"not null;check:rating >= 1 AND <= 5"`
	Comment string	`gorm:"text"`

	CreatedAt time.Time

	User UserModel `gorm:"foreignKey:UserID"`
	Room Room`gorm:"foreignKey:RoomID"`
}