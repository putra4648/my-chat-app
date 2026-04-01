package configs

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var SessionStore *session.Store

func SetupSession(app *fiber.App, store *redis.Storage) {
	SessionStore = session.NewStore(session.Config{
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		IdleTimeout:       30 * time.Minute,
		Storage:           store,
	})

	app.Use(session.New(session.Config{
		Store: SessionStore,
	}))
}
