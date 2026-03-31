package main

import (
	"context"
	"putra4648/my-chat-app/internal"
	"putra4648/my-chat-app/internal/modules"

	_ "github.com/joho/godotenv/autoload"
	"github.com/matzefriedrich/parsley/pkg/bootstrap"
)

func main() {

	ctx := context.Background()
	bootstrap.RunParsleyApplication(ctx, internal.NewApp,
		modules.ConfigureLogger,
		modules.ConfigureDatabase,
		modules.ConfigureStorage,
		modules.ConfigureRepositories,
		modules.ConfigureServices,
		modules.ConfigureFiber,
	)


}
