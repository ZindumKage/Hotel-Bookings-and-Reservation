package user

import (
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
)



func ToUserModel(u *user.User) *models.UserModel{
	if u == nil {
		return nil
	}
	return &models.UserModel{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		Password: u.Password,
		Role: string(u.Role),
		IsActive: u.IsActive,
	}
}

func ToUserDomain(m *models.UserModel) *user.User {
	if m == nil {
		return nil
	}
	return &user.User{
		ID: m.ID,
		Name: m.Name,
		Email: m.Email,
		Password: m.Password,
		Role: common.Role(m.Role),
		IsActive: m.IsActive,
	}
}

func ToUserDomains(models []models.UserModel) []*user.User {
	var result []*user.User
	for i := range models {
		result = append(result, ToUserDomain(&models[i]))
	}
	return result
}