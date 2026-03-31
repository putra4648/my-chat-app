package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DBMiddleware(pool *pgxpool.Pool) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		c.Locals("db", pool)
		return c.Next()
	}
}
