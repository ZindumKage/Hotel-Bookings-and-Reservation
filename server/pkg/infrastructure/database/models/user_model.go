package models

import (
	"time"
	"gorm.io/gorm"
)

type UserModel struct {
	ID uint `gorm:"primaryKey"`

	Name string
	Email string `gorm:"uniqueIndex"`
	Password string
	Role string

	IsActive bool
	IsEmailVerified bool

	LastLoginIP   string
	LastUserAgent string

	FailedLoginAttempts int
	AccountLockedUntil *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}