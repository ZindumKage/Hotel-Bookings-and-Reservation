package middleware

import (
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(c *fiber.Ctx) error {
	// get access token from cookie
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing access token")
	}

	// Parse and validate JWT
	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired access token")

	}

	// validate claims

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "access" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in token")
	}
	userID := uint(userIDFloat)

	// Attach user ID to context for handlers to use
	c.Locals("user_id", userID)

	return c.Next()
}
