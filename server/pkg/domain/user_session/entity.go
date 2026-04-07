package session

import "time"

type Session struct {
	ID uint

	UserID uint

	DeviceID string

	RefreshTokenHash string

	IPAddress string

	SessionToken string

	UserAgent string

	ExpiresAt time.Time

	CreatedAt time.Time
}