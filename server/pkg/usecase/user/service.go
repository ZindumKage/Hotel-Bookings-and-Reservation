package user

import (
	"errors"
	"strings"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/security"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/middleware"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo domain.Repository
	rateLimiter middleware.RedisRateLimiter
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
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

	token := uuid.New().String()
	hashedToken := security.HashToken(token)
	s.repo.SaveVerificationToken(user.ID, hashedToken)

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed), // In production, hash the password before storing
		Role:     common.RoleUser,
		IsActive: true,
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

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


	allowed, err := s.rateLimiter.Allow("login:"+email, 10, time.Minute)
	if err != nil {
		return nil, err
	}
	if !allowed{
		return nil, domain.ErrRateLimited
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

func (s *Service) GetAllUsers() ([]*domain.User, error) {
	return s.repo.FindAll()
}
