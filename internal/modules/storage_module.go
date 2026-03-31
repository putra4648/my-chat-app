package modules

import (
	"os"
	"strconv"

	"github.com/gofiber/storage/redis/v3"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func store() *redis.Storage {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		port = 6379
	}
	return redis.New(redis.Config{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     port,
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
}

func ConfigureStorage(registry types.ServiceRegistry) error {
	registry.Register(store, types.LifetimeSingleton)
	return nil
}
