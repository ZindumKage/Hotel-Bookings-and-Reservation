package session

import "time"

type Repository interface {
	CreateUserSession(userID uint, deviceID string, refreshTokenHash string, ip string, userAgent string, expiresAt time.Time) error

	DeleteSession(userID uint, deviceID string) error

	DeleteAllUserSessions(userID uint) error
	UpdateRefreshToken(sessionID uint, newHash string, expiresAt time.Time) error
	FindByToken(hash string) (*Session, error)

	GetUserSessions(userID uint) ([]*Session, error)
}
