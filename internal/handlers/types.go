package handlers

import "github.com/gofiber/fiber/v3"

type RouteHandler interface {
	Register(app *fiber.App)
}
