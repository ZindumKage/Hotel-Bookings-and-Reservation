package user

import "time"

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	FindAll(page, limit int) ([]*User,int, error)
	Update(user *User) error
	WithTx(txFunc func(repo Repository) error) error

	SaveVerificationToken(userID uint, token string, expiresAt time.Time)error
	FindVerificationByToken(token string) (*EmailVerification, error)
	DeleteVerificationToken(id uint) error

	SavePasswordResetToken(userID uint, token string, expiresAt time.Time)error
	FindPasswordResetByToken(token string) (*PasswordReset, error)
	DeletePasswordResetToken(id uint) error
}

