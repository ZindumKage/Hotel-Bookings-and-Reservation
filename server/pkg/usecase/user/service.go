package user

import (
	"errors"
	
	"strings"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/security"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo domain.Repository
	rateLimiter domain.RedisRateLimiter
}

func NewService(repo domain.Repository, rl domain.RedisRateLimiter) *Service {
	return &Service{repo: repo, rateLimiter: rl}
}
func (s *Service) Register(name, email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	existingUser, err := s.repo.FindByEmail(email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:            name,
		Email:           email,
		Password:        string(hashedPassword),
		Role:            common.RoleUser,
		IsActive:        true,
		IsEmailVerified: false,
	}

	err = s.repo.WithTx(func(txRepo domain.Repository) error {

		if err := txRepo.Create(user); err != nil {
			return err
		}

		token := uuid.New().String()
		hashedToken := security.HashToken(token)
		expiresAt := time.Now().Add(10 * time.Minute)

		if err := txRepo.SaveVerificationToken(user.ID, hashedToken, expiresAt); err != nil {
			return err
		}

		// send email after commit in real system

		return nil
	})

	if err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *Service) Login(email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	
	allowed, err := s.rateLimiter.Allow("login:"+email, 10, time.Minute)
	if err != nil {
		return nil, err
	}
	if !allowed{
		return nil, domain.ErrRateLimited
	}
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsEmailVerified {
		return nil, domain.ErrEmailNotVerified
	}

	if user.AccountLockedUntil != nil && user.AccountLockedUntil.After(time.Now()) {
		return nil, domain.ErrAccountLocked
	}



	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		user.FailedLoginAttempts++ 

		if user.FailedLoginAttempts >= 5 {
			lockTime := time.Now().Add(15 * time.Minute)
			user.AccountLockedUntil = &lockTime
			user.FailedLoginAttempts = 0
		}
		s.repo.Update(user)
		return nil, domain.ErrInvalidCredentials
	}
	user.FailedLoginAttempts = 0
	user.AccountLockedUntil = nil
	s.repo.Update(user)

	user.Password = ""

	return user, nil
}

func (s *Service) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetAllUsers(page, limit int) ([]*domain.User,int, error) {
	users, total, err := s.repo.FindAll(page, limit)
	if err != nil {
		return nil, 0 , err 
	}
	return users, total, nil 
}


func (s *Service) RequestPasswordReset(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		// Do NOT reveal whether user exists
		return nil
	}

	token := uuid.New().String()
	hashedToken := security.HashToken(token)

	expiresAt := time.Now().Add(10 * time.Minute)

	if err := s.repo.SavePasswordResetToken(user.ID, hashedToken, expiresAt); err != nil {
		return err
	}
	
	// TODO: send email with raw token
	// example link:
	// https://yourdomain.com/reset-password?token=abc123

	return nil
}

func (s *Service) ResetPassword(token string, newPassword string) error {
	hashedToken := security.HashToken(token)

	record, err := s.repo.FindPasswordResetByToken(hashedToken)
	if err != nil {
		return domain.ErrInvalidOrExpiredToken
	}

	if record.ExpiresAt.Before(time.Now()) {
		return domain.ErrInvalidOrExpiredToken}

	user, err := s.repo.FindByID(record.UserID)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.FailedLoginAttempts = 0
	user.AccountLockedUntil = nil

	if err := s.repo.Update(user); err != nil {
		return err
	}

	if err := s.repo.DeletePasswordResetToken(record.ID); err != nil {
		return err
	}

	return nil
}