package repositories

import (
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/common"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	"gorm.io/gorm"
)




type UserRepository struct {
	 db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.UserModel) error {
	model := models.UserModel{
		Name: user.Name,
		Email: user.Email,
		Password: user.Password,
		Role: string(user.Role),
		IsActive: user.IsActive,
	}
	err := r.db.Create(&model).Error
	if err != nil {
		return  err
}
 user.ID = model.ID
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var model models.UserModel
	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &user.User{
		ID: model.ID,
		Name: model.Name,
		Email: model.Email,
		Password: model.Password,
		Role: common.Role(model.Role),
		IsActive: model.IsActive,
	}, nil
}

func (r *UserRepository) FindByID(id uint) (*user.User, error) {
	var model models.UserModel
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &user.User{
		ID: model.ID,
		Name: model.Name,
		Email: model.Email,
		Password: model.Password,
		Role: common.Role(model.Role),
		IsActive: model.IsActive,
	}, nil
}

func (r *UserRepository) Update(user *user.User) error {
	var model models.UserModel
	err := r.db.First(&model, user.ID).Error
	if err != nil {
		return err
	}
	model.Name = user.Name
	model.Email = user.Email
	model.Password = user.Password
	model.Role = string(user.Role)
	model.IsActive = user.IsActive
	return r.db.Save(&model).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserModel{}, id).Error
}	

func (r *UserRepository) FindAll() ([]*user.User, error) {
	var models []models.UserModel
	err := r.db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	users := make([]*user.User, len(models))
	for i, m := range models {
		users[i] = &user.User{
			ID: m.ID,
			Name: m.Name,
			Email: m.Email,
			Password: m.Password,
			Role: common.Role(m.Role),
			IsActive: m.IsActive,
		}
	}
	return users, nil
}