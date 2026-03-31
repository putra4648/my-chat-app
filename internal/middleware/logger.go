package middleware

import (
	middleware "github.com/gofiber/contrib/v3/zap"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) func(c fiber.Ctx) error {
	return middleware.New(
		middleware.Config{
			Logger: logger,
		},
	)
}
