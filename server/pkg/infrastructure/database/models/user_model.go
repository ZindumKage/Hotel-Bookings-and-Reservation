package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	ID 	uint
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string
	IsActive bool
}
