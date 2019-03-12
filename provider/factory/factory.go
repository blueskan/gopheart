package factory

import (
	"github.com/blueskan/gopheart/provider/cassandra"
	"github.com/blueskan/gopheart/provider/couchbase"
	"github.com/blueskan/gopheart/provider/mssql"
	"github.com/blueskan/gopheart/provider/mysql"
	"github.com/blueskan/gopheart/provider/postgresql"
	"time"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/provider"
	"github.com/blueskan/gopheart/provider/memcache"
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
	case "memcache":
		return memcache.NewMemcacheProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "mysql":
		return mysql.NewMysqlProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "postgresql":
		return postgresql.NewPostgresqlProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "mssql":
		return mssql.NewMssqlProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "cassandra":
		return cassandra.NewCassandraProvider(
			name,
			config.Source,
			timeout,
			interval,
			config.RetryPolicy.DownThreshold,
			config.RetryPolicy.UpThreshold,
		)
	case "couchbase":
		return couchbase.NewCouchbaseProvider(
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
