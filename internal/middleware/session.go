package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

func SessionMiddleware(store *redis.Storage) func(c fiber.Ctx) error {
	return session.New(session.Config{
		CookieSameSite:    "Lax", // CSRF protection
		CookieSessionOnly: true,
		IdleTimeout:       30 * time.Minute, // Session timeout
		Storage:           store,
	})
}
