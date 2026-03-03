package user

import (

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

func (r *UserRepository) Create(user *user.User) error {
model := ToUserModel(user)
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
	return ToUserDomain(&model), nil
}

func (r *UserRepository) FindByID(id uint) (*user.User, error) {
	var model models.UserModel
	err := r.db.First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return ToUserDomain(&model), nil
}

func (r *UserRepository) Update(user *user.User) error {
	model := ToUserModel(user)
	return r.db.Save(model).Error
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
	users := make([]*user.User, 0,len(models))
	for _, m := range models {
		users = append(users, ToUserDomain(&m))
	}
	return users, nil
}