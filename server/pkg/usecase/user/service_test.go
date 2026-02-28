package user_test

import (
	"testing"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	usecase "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/user"
)

// Mock repository
type mockRepo struct {
	users map[string]*domain.User
}

func (m *mockRepo) Create(user *domain.User) error {
	user.ID = 1
	m.users[user.Email] = user
	return nil
}

func (m *mockRepo) FindByEmail(email string) (*domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockRepo) FindByID(id uint) (*domain.User, error) {
	return nil, nil
}

func TestRegisterUser(t *testing.T) {
	mock := &mockRepo{users: make(map[string]*domain.User)}
	service := usecase.NewService(mock) // <--- use the usecase package

	user, err := service.Register("Stanley", "stan@test.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != "stan@test.com" {
		t.Fatal("email mismatch")
	}

	t.Logf("Test passed! Created user: %+v", user)
}

func TestLoginUser(t *testing.T) {
	mock := &mockRepo{users: make(map[string]*domain.User)}
	service := usecase.NewService(mock) // <--- use the usecase package

	// First, register a user
	_, err := service.Register("Stanley", "stan@test.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	// Then, try to login
	user, err := service.Login("stan@test.com", "password")
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != "stan@test.com" {
		t.Fatal("email mismatch")
	}

	t.Logf("Test passed! Logged in user: %+v", user)
}		