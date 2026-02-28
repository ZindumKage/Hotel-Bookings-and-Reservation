// package models



// import "gorm.io/gorm"

// type User struct {
// 	gorm.Model
//  Name       string    `json:"name"`
//     Email      string    `gorm:"unique;not null" json:"email"`
//     Password   string    `gorm:"-" json:"-"` // The '-' tag ignores the field in JSON output
//   Role     Role `gorm:"type:varchar(20);default:'USER'"`
// 	IsActive bool `gorm:"default:true"`

// 	Bookings []Booking
// 	Reviews  []Review
// }


package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string
	IsActive bool
}