package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Versioning() fiber.Handler {
	return func(c *fiber.Ctx) error {
		version := c.Params("version")
		if version != "v1" {
			return c.Status(404).SendString("Version not found")
		}
		return c.Next()
	}
}