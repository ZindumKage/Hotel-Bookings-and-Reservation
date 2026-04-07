package handlers

import (
	"os"
	"time"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/auth"
	session "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/user_session"

	"github.com/gofiber/fiber/v2"
)

func RefreshTokenHandler(sessionRepo session.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// 1. get refresh token
		rtCookie := c.Cookies("refresh_token")
		if rtCookie == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token")
		}

		hashedToken := auth.HashToken(rtCookie)

		// 2. find session
		sess, err := sessionRepo.FindByToken(hashedToken)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
		}

		// 3. check expiry
		if sess.ExpiresAt.Before(time.Now()) {
			return fiber.NewError(fiber.StatusUnauthorized, "Token expired")
		}

		// 4. generate new tokens
		newAccessToken, _ := auth.GenerateAccessToken(sess.UserID)
		newRefreshToken, _ := auth.GenerateRefreshToken()

		newHash := auth.HashToken(newRefreshToken)

		// 5. rotate session (IMPORTANT)
		err = sessionRepo.CreateUserSession(
			sess.UserID,
			sess.DeviceID,
			newHash,
			c.IP(),
			c.Get("User-Agent"),
			time.Now().Add(7*24*time.Hour),
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create session")
		}

		// delete old session
		err = sessionRepo.DeleteSession(sess.UserID, sess.DeviceID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete old session")
		}

		// 6. set cookies
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			HTTPOnly: true,
			Secure:   os.Getenv("ENV") == "production",
			MaxAge:   60 * 15,
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			HTTPOnly: true,
			Secure:   os.Getenv("ENV") == "production",
			MaxAge:   60 * 60 * 24 * 7,
		})

		return c.JSON(fiber.Map{
			"message": "Tokens refreshed",
		})
	}
}