package models

import "time"

type UserSessionModel struct {
	ID uint `gorm:"primaryKey"`

	UserID   uint   `gorm:"index:idx_user_device;uniqueIndex:uid_device"`
DeviceID string `gorm:"index:idx_user_device;uniqueIndex:uid_device"`
	RefreshTokenHash string 
	IPAddress string
	UserAgent string
	SessionToken string `gorm:"uniqueIndex"`

	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}