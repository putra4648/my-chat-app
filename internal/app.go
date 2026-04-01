package internal

import (
	"context"
	"os"
	"putra4648/my-chat-app/internal/handlers"
	"putra4648/my-chat-app/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/storage/redis/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
	"go.uber.org/zap"
)

type parsleyApplication struct {
	app    *fiber.App
	logger *zap.Logger
}

var _ bootstrap.Application = &parsleyApplication{}

// NewApp Creates the main application service instance. This constructor function gets invoked by Parsley; add parameters for all required services.
func NewApp(app *fiber.App, routeHandlers []handlers.RouteHandler, pool *pgxpool.Pool, logger *zap.Logger, store *redis.Storage) bootstrap.Application {

	// Middlewares
	app.Use(middleware.LoggerMiddleware(logger))
	app.Use(middleware.HelmetMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.SessionMiddleware(store))
	app.Use(middleware.CSRFMiddleware())
	app.Use(middleware.DBMiddleware(pool))
	app.Use(middleware.AuthMiddleware())

	// Register RouteHandler services with the resolved Fiber instance.
	for _, routeHandler := range routeHandlers {
		routeHandler.Register(app)
	}

	return &parsleyApplication{
		app: app,
	}
}

// Run The entrypoint for the Parsley application.
func (a *parsleyApplication) Run(_ context.Context) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	return a.app.Listen(":" + port)
}
