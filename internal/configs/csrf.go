package configs

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func CsrfConfig(sessionStore *session.Store) fiber.Handler {
	return csrf.New(csrf.Config{
		Extractor: extractors.Chain(
			extractors.FromHeader("X-Csrf-Token"),
			extractors.FromForm("_csrf"),
		),
		Session: sessionStore,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			log.Errorf("CSRF Error: %v Request: %v From: %v\n", err, c.OriginalURL(), c.IP())
			switch c.Accepts("html", "json") {
			case "json":
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "403 Forbidden"})
			case "html":
				return c.Status(fiber.StatusForbidden).Render("error", fiber.Map{
					"title":      "Error",
					"error":      "403 Forbidden",
					"error_code": "403",
				})
			default:
				return c.Status(fiber.StatusForbidden).SendString("403 Forbidden")
			}
		},
	})
}
