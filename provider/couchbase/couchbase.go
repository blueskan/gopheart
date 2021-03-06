package couchbase

import (
	"github.com/couchbase/go-couchbase"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type couchbaseProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewCouchbaseProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &couchbaseProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (cp couchbaseProvider) GetName() string {
	return cp.name
}

func (cp couchbaseProvider) GetInterval() time.Duration {
	return cp.interval
}

func (cp couchbaseProvider) GetDownThreshold() int64 {
	return cp.downThreshold
}

func (cp couchbaseProvider) GetUpThreshold() int64 {
	return cp.upThreshold
}

func (cp couchbaseProvider) Heartbeat() bool {
	_, err := couchbase.Connect(cp.connectionString)

	if err != nil {
		return false
	}

	return true
}
