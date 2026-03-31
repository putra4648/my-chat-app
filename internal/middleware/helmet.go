package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
)

// change this to use in production
func HelmetMiddleware() func(c fiber.Ctx) error {
	return helmet.New()
}
