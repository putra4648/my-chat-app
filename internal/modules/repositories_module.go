package modules

import (
	"putra4648/my-chat-app/internal/repositories"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureRepositories(registry types.ServiceRegistry) error {
	// Register repository services with the resolved database connection pool.
	// add more repositories here
	registration.RegisterTransient(registry, repositories.NewUserRepository)
	return nil
}
