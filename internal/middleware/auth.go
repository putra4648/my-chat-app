package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func AuthMiddleware() func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		sess := session.FromContext(c)

		// skip login and register route
		if c.Path() == "/login" || c.Path() == "/register" {
			return c.Next()
		}

		if sess == nil || sess.Get("authenticated") != true {
			return c.Redirect().To("/login")
		}

		return c.Next()
	}
}
