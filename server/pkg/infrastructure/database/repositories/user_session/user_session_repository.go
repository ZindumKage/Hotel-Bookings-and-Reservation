package repository

import (
	"time"

	session "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user_session"
	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/infrastructure/database/models"

	"gorm.io/gorm"
	"crypto/rand"
	"encoding/hex"

)

type UserSessionRepository struct {
	db *gorm.DB
}

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

func (r *UserSessionRepository) CreateUserSession(
	userID uint,
	deviceID string,
	refreshTokenHash string,
	ip string,
	userAgent string,
	expiresAt time.Time,
) error {

	sessionModel := models.UserSessionModel{
		UserID:           userID,
		DeviceID:         deviceID,
		RefreshTokenHash: refreshTokenHash,
		IPAddress:        ip,
		SessionToken:    generateSessionToken(), // You can implement this function to create a unique session token
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
		ExpiresAt:        expiresAt,
	}

	return r.db.Create(&sessionModel).Error
}

func (r *UserSessionRepository) DeleteSession(userID uint, deviceID string) error {
	return r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&models.UserSessionModel{}).Error
}

func (r *UserSessionRepository) DeleteAllUserSessions(userID uint) error {

	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.UserSessionModel{}).
		Error
}

func (r *UserSessionRepository) DeleteBySessionToken(token string) error {
	return r.db.
		Where("session_token = ?", token).
		Delete(&models.UserSessionModel{}).
		Error
}

func (r *UserSessionRepository) FindByToken(hash string) (*session.Session, error) {

	var model models.UserSessionModel

	err := r.db.
		Where("refresh_token_hash = ?", hash).
		First(&model).Error

	if err != nil {
		return nil, err
	}

	return &session.Session{
		ID:              model.ID,
		UserID:          model.UserID,
		DeviceID:        model.DeviceID,
		RefreshTokenHash: model.RefreshTokenHash,
		SessionToken:    model.SessionToken,
		IPAddress:       model.IPAddress,
		UserAgent:       model.UserAgent,
		ExpiresAt:       model.ExpiresAt,
		CreatedAt:       model.CreatedAt,
	}, nil
}

func (r *UserSessionRepository) GetUserSessions(userID uint) ([]*session.Session, error) {

	var modelsList []models.UserSessionModel

	err := r.db.
		Where("user_id = ?", userID).
		Find(&modelsList).Error

	if err != nil {
		return nil, err
	}

	var sessions []*session.Session

	for _, m := range modelsList {

		sessions = append(sessions, &session.Session{
			ID:               m.ID,
			UserID:           m.UserID,
			DeviceID:         m.DeviceID,
			RefreshTokenHash: m.RefreshTokenHash,
			IPAddress:        m.IPAddress,
			UserAgent:        m.UserAgent,
			ExpiresAt:        m.ExpiresAt,
			CreatedAt:        m.CreatedAt,
		})
	}

	return sessions, nil
}

func (r *UserSessionRepository) UpdateRefreshToken(
	sessionID uint,
	newHash string,
	expiresAt time.Time,
) error {
	return r.db.Model(&models.UserSessionModel{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"refresh_token_hash": newHash,
			"expires_at":         expiresAt,
		}).Error
}