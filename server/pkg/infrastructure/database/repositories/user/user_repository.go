package repositories

import (
	"context"
	"time"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	model := ToUserModel(user)

	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	user.ID = model.ID
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {

	var model models.UserModel

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return ToUserDomain(&model), nil
}
func (r *UserRepository) FindByID(ctx context.Context ,id uint) (*domain.User, error) {
	var model models.UserModel

	err := r.db.WithContext(ctx).First(&model, id).Error 
	
	if err != nil {
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

func (r *UserRepository) WithTx(txFunc func(repo domain.Repository) error) error {

	return r.db.Transaction(func(tx *gorm.DB) error {

		repo := &UserRepository{db: tx}

		return txFunc(repo)
	})
}

//////////////////////////////
//// REFRESH TOKENS
//////////////////////////////

func (r *UserRepository) SaveRefreshToken(userID uint, token string, expiresAt time.Time) error {

	model := models.RefreshTokenModel{
		UserID:    userID,
		Token:     token, // hashed token
		ExpiresAt: expiresAt,
	}

	return r.db.Create(&model).Error
}

func (r *UserRepository) FindRefreshToken(token string) (*domain.RefreshToken, error) {

	var model models.RefreshTokenModel

	if err := r.db.
		Where("token = ?", token).
		First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.RefreshToken{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *UserRepository) ReplaceRefreshToken(id uint, token string, expiresAt time.Time) error {

	return r.db.Model(&models.RefreshTokenModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"token":      token,
			"expires_at": expiresAt,
		}).Error
}

func (r *UserRepository) DeleteRefreshToken(id uint) error {
	return r.db.Delete(&models.RefreshTokenModel{}, id).Error
}

func (r *UserRepository) DeleteAllRefreshTokens(userID uint) error {

	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.RefreshTokenModel{}).
		Error
}

//////////////////////////////
//// EMAIL VERIFICATION
//////////////////////////////

func (r *UserRepository) SaveVerificationToken(userID uint, token string, expiresAt time.Time) error {

	model := models.EmailVerificationModel{
		UserID:    userID,
		Token:     token, // hashed token
		ExpiresAt: expiresAt,
	}

	return r.db.Create(&model).Error
}

func (r *UserRepository) FindVerificationByToken(token string) (*domain.EmailVerification, error) {

	var model models.EmailVerificationModel

	if err := r.db.
		Where("token = ?", token).
		First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.EmailVerification{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *UserRepository) DeleteVerificationToken(userID uint) error {

	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.EmailVerificationModel{}).
		Error
}

//////////////////////////////
//// PASSWORD RESET
//////////////////////////////

func (r *UserRepository) SavePasswordResetToken(userID uint, token string, expiresAt time.Time) error {

	model := models.PasswordResetModel{
		UserID:    userID,
		Token:     token, // hashed token
		ExpiresAt: expiresAt,
	}

	return r.db.Create(&model).Error
}

func (r *UserRepository) FindPasswordResetByToken(token string) (*domain.PasswordReset, error) {

	var model models.PasswordResetModel

	if err := r.db.
		Where("token = ?", token).
		First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.PasswordReset{
		ID:        model.ID,
		UserID:    model.UserID,
		Token:     model.Token,
		ExpiresAt: model.ExpiresAt,
		CreatedAt: model.CreatedAt,
	}, nil
}

func (r *UserRepository) DeletePasswordResetToken(id uint) error {

	return r.db.
		Delete(&models.PasswordResetModel{}, id).
		Error
}

//////////////////////////////
//// USER SESSIONS
//////////////////////////////

func (r *UserRepository) DeleteAllUserSessions(userID uint) error {

	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.UserSessionModel{}).
		Error
}