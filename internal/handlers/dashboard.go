package handlers

import "github.com/gofiber/fiber/v3"

type dashboardHandler struct {
}

func NewDashboardHandler() RouteHandler {
	return &dashboardHandler{}
}

func (h *dashboardHandler) Register(app *fiber.App) {
	app.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{
			"Title":     "Dashboard",
			"BodyClass": "chat-page",
		})
	})
}
