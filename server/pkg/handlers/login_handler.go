package handlers

import (
	"context"
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/config"
	userservice "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/usecase/user"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

type AuthHandler struct {
	userService *userservice.Service
}

func NewAuthHandler(userService *userservice.Service) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {

	oauthCfg := config.OAuthConfig()

	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "missing code")
	}

	cookieState := c.Cookies("oauth_state")
	if state == "" || state != cookieState {
		return fiber.NewError(fiber.StatusBadRequest, "invalid state")
	}

	token, err := oauthCfg.Exchange(context.Background(), code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "token exchange failed")
	}

	client := oauthCfg.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to fetch user info")
	}
	defer resp.Body.Close()

	var googleUser GoogleUser

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "decode failed")
	}

	if !googleUser.VerifiedEmail {
		return fiber.NewError(fiber.StatusBadRequest, "email not verified by google")
	}

	ip := c.IP()
	userAgent := c.Get("User-Agent")

	deviceID := c.Cookies("device_id")

if deviceID == "" {
	deviceID = uuid.New().String()

	c.Cookie(&fiber.Cookie{
		Name:     "device_id",
		Value:    deviceID,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 365, // 1 year
	})
}

	payload, err := h.userService.OAuthLogin(
		c.Context(),
		googleUser.Name,
		googleUser.Email,
		ip,
		userAgent,
		deviceID,
	)

	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    payload.AccessToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		MaxAge:   900,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    payload.RefreshToken,
		HTTPOnly: true,
		Secure:   os.Getenv("ENV") == "production",
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
	})
	

	return c.Redirect("http://localhost:3000/")
}