package user

import (
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"
	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User ) error {
	model := ToUserModel(user)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var model models.UserModel

	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}

	return ToUserDomain(&model), nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var model models.UserModel

	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return ToUserDomain(&model), nil
}

func (r *UserRepository) Update(user *domain.User) error {
	model := ToUserModel(user)
	return r.db.Save(&model).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.UserModel{}, id).Error
}

func (r *UserRepository) FindAll(page, limit int) ([]*domain.User, int, error) {
	var users []models.UserModel
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&models.UserModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return ToUserDomains(users), int(total), nil
}