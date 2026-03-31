package modules

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

var _ types.ModuleFunc = ConfigureFiber

func ConfigureFiber(registry types.ServiceRegistry) error {
	engine := html.New("../views", ".html")

	engine.Debug(true)
	engine.Reload(true)

	registration.RegisterInstance(registry, fiber.Config{
		AppName:     "my-chat-app",
		Immutable:   true,
		Views:       engine,
		ViewsLayout: "layouts/main",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

			// Return status code with error message
			return c.Status(code).Render("error", fiber.Map{
				"error":      err.Error(),
				"error_code": code,
				"title":      "Error",
			})
		},
	})

	registry.Register(newFiber, types.LifetimeSingleton)
	registry.RegisterModule(RegisterRouteHandlers)

	return nil
}

func newFiber(config fiber.Config) *fiber.App {
	return fiber.New(config)
}
