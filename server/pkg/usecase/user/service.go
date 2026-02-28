package user

import (
	"errors"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
)



type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}		
}

func (s *Service) Register(name, email, password string) (*domain.User, error) {
	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	} 

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	user := &domain.User{
		Name: name,
		Email: email,
		Password: string(hashed), // In production, hash the password before storing
		Role: common.RoleUser,
		IsActive: true,
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	return user, nil
}

func (s *Service) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) GetAllUsers() ([]*domain.User, error) {
	return s.repo.FindAll()	
}