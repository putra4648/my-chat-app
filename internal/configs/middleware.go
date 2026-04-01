package configs

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/storage/redis/v3"
	"go.uber.org/zap"
)

func SetupMiddleware(app *fiber.App, logger *zap.Logger, store *redis.Storage) {
	// 1. Recover middleware
	app.Use(recover.New())

	// 2. Logger middleware
	app.Use(LoggerConfig(logger))

	// 3. CORS middleware
	app.Use(CorsConfig())

	// 4. Helmet middleware
	app.Use(helmet.New())

	// 5. Session middleware (Sets global SessionStore in the package)
	SetupSession(app, store)

	// 6. CSRF middleware (Uses the SessionStore from SetupSession)
	app.Use(CsrfConfig(SessionStore))
}
