package modules

import (
	"github.com/matzefriedrich/parsley/pkg/types"
	"go.uber.org/zap"
)

func newLogger() *zap.Logger {
	l, _ := zap.NewProduction()
	return l
}

func ConfigureLogger(registry types.ServiceRegistry) error {
	registry.Register(newLogger, types.LifetimeSingleton)
	return nil
}


