package factory

import (
	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/provider"
)

type ProviderFactory interface {
	CreateProvider(name string, config config.HealthCheck) provider.Provider
}
