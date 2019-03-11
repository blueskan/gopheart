package main

import (
	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/provider"
	"github.com/blueskan/gopheart/provider/factory"
)

func main() {
	config := config.Config{}
	config.FromYaml("./config.yaml")

	providerFactory := factory.NewProviderFactory()
	providers := make([]provider.Provider, 0)

	// Implement factory pattern
	for key, value := range config.HealthChecks {
		providers = append(providers, providerFactory.CreateProvider(key, value))
	}

	scheduler := provider.NewScheduler(providers)
	scheduler.Schedule()
}
