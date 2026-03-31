package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func CSRFMiddleware() func(c fiber.Ctx) error {
	return csrf.New(csrf.Config{
		Extractor: extractors.Chain(
			extractors.FromHeader("X-Csrf-Token"),
			extractors.FromForm("_csrf"),
		),
		Session: session.NewStore(),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Log the error without crashing the server
			log.Errorf("CSRF Error: %v Request: %v From: %v\n", err, c.OriginalURL(), c.IP())

			// check accepted content types
			switch c.Accepts("html", "json") {
			case "json":
				// Return a 403 Forbidden response for JSON requests
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "403 Forbidden",
				})
			case "html":
				// Return a 403 Forbidden response for HTML requests
				return c.Render("error", fiber.Map{
					"title":      "Error",
					"error":      "403 Forbidden",
					"error_code": "403",
				})
			default:
				// Return a 403 Forbidden response for all other requests
				return c.Status(fiber.StatusForbidden).SendString("403 Forbidden")
			}
		},
	})
}
