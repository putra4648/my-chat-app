package modules

import (
	"context"
	"fmt"
	"log"
	"os"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func db() *pgxpool.Pool {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	config.MaxConns = 5
	config.MinIdleConns = 2

	log.Println("Registering gofrs/uuid support for all connections in the pool")
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	return pool
}

func ConfigureDatabase(registry types.ServiceRegistry) error {
	registry.Register(db, types.LifetimeSingleton)
	return nil
}
