package models

import "gorm.io/gorm"



type Room struct { 
	gorm.Model
	Name        string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:text"`
	Status      string  `gorm:"type:varchar(50);not null"`
	Price       int64  `gorm:"type:decimal(10,2);not null"`
	Amenities   string  `gorm:"type:text"` // Comma-separated list of amenities
	RoomNumber  string  `gorm:"type:varchar(50);not null;unique"`
}

func (Room) TableName() string {
	return "rooms"
}