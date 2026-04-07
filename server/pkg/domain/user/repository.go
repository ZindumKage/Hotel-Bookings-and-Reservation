package user

import (
	"context"
	"time"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context ,id uint) (*User, error)
	FindAll(page, limit int) ([]*User,int, error)
	Update(user *User) error
	WithTx(txFunc func(repo Repository) error) error

	SaveVerificationToken(userID uint, token string, expiresAt time.Time)error
	FindVerificationByToken(token string) (*EmailVerification, error)
	DeleteVerificationToken(id uint) error

	SaveRefreshToken(userID uint, token string, expiresAt time.Time)error
	FindRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(id uint)error

	ReplaceRefreshToken(id uint, token string, expiresAt time.Time) error

	SavePasswordResetToken(userID uint, token string, expiresAt time.Time)error
	FindPasswordResetByToken(token string) (*PasswordReset, error)
	DeletePasswordResetToken(id uint) error

	DeleteAllUserSessions(userID uint) error
	DeleteAllRefreshTokens(userID uint) error

}

