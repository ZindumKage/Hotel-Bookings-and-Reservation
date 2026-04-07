package models

import "time"

type RefreshTokenModel struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"index"`
	Token string `gorm:"uniqueIndex"`

	ExpiresAt time.Time
	CreatedAt time.Time
}

