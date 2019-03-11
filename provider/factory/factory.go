package factory

import (
	"time"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/provider"
	"github.com/blueskan/gopheart/provider/mongodb"
	"github.com/blueskan/gopheart/provider/redis"
	"github.com/blueskan/gopheart/provider/url"
)

type providerFactory struct {
}

func NewProviderFactory() ProviderFactory {
	return &providerFactory{}
}

func (pf *providerFactory) CreateProvider(name string, config config.HealthCheck) provider.Provider {
	timeout, _ := time.ParseDuration(config.RetryPolicy.Timeout)
	interval, _ := time.ParseDuration(config.CheckInterval)

	switch config.Type {
	case "url":
		return url.NewUrlProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "redis":
		return redis.NewRedisProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "mongodb":
		return mongodb.NewMongoDbProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	default:
		panic("Cannot parse url type `" + config.Type + "`")
	}
}
