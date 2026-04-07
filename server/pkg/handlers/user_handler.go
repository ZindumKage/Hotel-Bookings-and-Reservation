

package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
	})
}