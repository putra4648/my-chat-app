package main

import (
	"log"
	"os"
	"putra4648/my-chat-app/internal/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v3"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	engine := html.New("../views", ".html")
	engine.Reload(true)
	engine.Debug(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("home1", fiber.Map{
			"Title": "Home",
		}, "layouts/main")
	})

	handlers.WSRoute(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
