package modules

import (
	"putra4648/my-chat-app/internal/services"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

func ConfigureServices(registry types.ServiceRegistry) error {
	registration.RegisterTransient(registry, services.NewUserService)
	registration.RegisterTransient(registry, services.NewChatService)
	return nil
}
