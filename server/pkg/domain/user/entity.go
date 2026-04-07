package user

import (
	"time"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	"errors"
)


var (
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrAccountLocked 	   = errors.New("account is temporarily locked")
	ErrEmailNotVerified	   = errors.New("email not verified")
	ErrRateLimited         = errors.New("too many login attempts, try again later")
	ErrInvalidOrExpiredToken = errors.New("token is expired or invalid")
)

type User struct {
	ID        uint
	Name      string
	Email     string
	Password  string
	Role      common.Role
	IsActive  bool
	IsEmailVerified bool

	FailedLoginAttempts int
	AccountLockedUntil *time.Time
	
	LastLoginIP         string
	LastUserAgent       string 
}

type EmailVerification struct {
	ID uint
	UserID uint
	Token string
	ExpiresAt time.Time
	CreatedAt time.Time
}
type PasswordReset struct {
	ID uint
	UserID uint
	Token string
	ExpiresAt time.Time
	CreatedAt time.Time
}
type AuthPayload struct {
	User         *User
	AccessToken  string
	RefreshToken string
	DeviceID 	string
}

type RefreshToken struct {
	ID uint
	UserID uint
	Token string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UserSession struct {
	ID               uint
	UserID           uint
	DeviceID         string
	RefreshTokenHash string
	IPAddress        string
	UserAgent        string
	ExpiresAt        time.Time
	CreatedAt        time.Time
}