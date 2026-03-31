package modules

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/template/html/v3"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

var _ types.ModuleFunc = ConfigureFiber

func ConfigureFiber(registry types.ServiceRegistry) error {
	engine := html.New("../views", ".html")
	engine.Reload(true)

	registration.RegisterInstance(registry, fiber.Config{
		AppName:     "my-chat-app",
		Immutable:   true,
		Views:       engine,
		ViewsLayout: "layouts/main",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			var pgErr *pgconn.PgError
			if errors.As(err, &e) {
				return c.Status(e.Code).Render("error", fiber.Map{
					"error":      err.Error(),
					"error_code": e.Code,
					"title":      "Error",
					"BodyClass":  "auth-page",
				})
			} else if errors.As(err, &pgErr) {
				log.Error(pgErr)
				return c.Status(fiber.StatusInternalServerError).Render("error", fiber.Map{
					"error":      "Internal Server Error",
					"error_code": fiber.ErrInternalServerError,
					"title":      "Error",
					"BodyClass":  "auth-page",
				})
			} else {
				code := fiber.StatusInternalServerError
				return c.Status(code).Render("error", fiber.Map{
					"error":      err.Error(),
					"error_code": code,
					"title":      "Error",
					"BodyClass":  "auth-page",
				})
			}

		},
	})

	registry.Register(newFiber, types.LifetimeSingleton)
	registry.RegisterModule(RegisterRouteHandlers)

	return nil
}

func newFiber(config fiber.Config) *fiber.App {
	return fiber.New(config)
}
