package handler

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/config"
	
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

func GoogleCallback(c *fiber.Ctx) error {
	oauthCfg := config.OAuthConfig()

	// Get code & state from query parameters
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing code")

	}
	// verify state to prevent CSRF
	cookieState := c.Cookies("oauth_state")
	if state == "" || state != cookieState {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid state")
	}
	// Exchange code for token
	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to exchange token")
	}
	// Get Google user info
	client := oauthCfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user info")
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to decode user info")
	}
	//   Find or Create User
	var user models.User
err = config.DB.Where("email = ?", googleUser.Email).First(&user).Error

if err != nil {

	if errors.Is(err, gorm.ErrRecordNotFound) {

		user = models.User{
			Name:  googleUser.Name,
			Email: googleUser.Email,
			Role:  "USER",
		}

		if err := config.DB.Create(&user).Error; err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
		}

	} else {
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}
}

	//Generate Jwt Tokens
	accessToken, err := config.GenerateAccessToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Access token failed")
	}
	refreshToken, err := config.GenerateRefreshToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Refresh token failed")
	}

	// save  refresh token in Db

	if err := config.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save refresh token")
	}
	// set secure cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		MaxAge:   900,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 3, // 3 days
	})

	return c.Redirect("http://localhost:3000/")

}
