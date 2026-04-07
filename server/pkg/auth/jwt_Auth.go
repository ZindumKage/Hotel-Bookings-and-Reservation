package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type TokenBlacklist interface {
	IsBlacklisted(token string) bool
}

const ContextUserID = "user_id"

func JWTMiddleware(blacklist TokenBlacklist, tokenType string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token missing",
			})
		}

		


		if blacklist != nil && blacklist.IsBlacklisted(HashToken(tokenString)) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token revoked",
			})
		}
	
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

	
		if claims.Type != tokenType {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token type",
			})
		}

	
		c.Locals(ContextUserID, claims.UserID)
		c.Locals("claims", claims)

		return c.Next()
	}
}