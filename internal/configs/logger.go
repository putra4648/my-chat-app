package configs

import (
	middlewareLogger "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func LoggerConfig(logger *zap.Logger) fiber.Handler {
	return middlewareLogger.New(
		middlewareLogger.Config{
			Logger: logger,
		},
	)
}
