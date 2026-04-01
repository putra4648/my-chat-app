package configs

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CorsConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "https://i.pravatar.cc", "https://www.gravatar.com"},
		AllowHeaders: []string{"Content-Type", "Origin", "Accept", "X-CSRF-Token", "Authorization", "X-Requested-With"},
	})
}
