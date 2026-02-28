package handler

import (
	"errors"
	"os"


	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/config"
	
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// RefreshTokenHandler issues a new access token using a valid refresh token
func RefreshTokenHandler(c *fiber.Ctx) error {
	// 1️⃣ get refresh token cookie
	
	rtCookie := c.Cookies("refresh_token")
	if rtCookie == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token")
	}
	hashedToken := config.HashToken(rtCookie)

	// 2️⃣ parse & validate JWT
	token, err := jwt.Parse(rtCookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
	}
	userID := uint(userIDFloat)

	// 3️⃣ verify token exists in DB
	var rt models.RefreshToken
	if err := config.DB.Where("user_id = ? AND token = ?", userID, hashedToken).First(&rt).Error; err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Refresh token not found")
	}

	// 4️⃣ issue new tokens
	newAccessToken, err := config.GenerateAccessToken(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	newRefreshToken, err := config.GenerateRefreshToken(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate refresh token")
	}

	// 5️⃣ rotate refresh token in DB
	if err := config.DB.Delete(&rt).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete old refresh token")
	}

	if err := config.SaveRefreshToken(userID, newRefreshToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save new refresh token")
	}

	// 6️⃣ set cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 15,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return c.JSON(fiber.Map{
		"message": "Tokens refreshed",
	})
}
