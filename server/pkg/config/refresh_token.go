package config

import (
	"crypto/sha256"
	"encoding/hex"

)

func HashToken(token string) string{
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
func SaveRefreshToken(userID uint, token string) error {

	hashed := HashToken(token)

	// delete old tokens
	if err := DB.Where("user_id = ?", userID).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return err
	}

	refresh := models.RefreshToken{
		UserID: userID,
		Token:  hashed,
	}

	return DB.Create(&refresh).Error
}
