package hbase

import (
	"github.com/tsuna/gohbase"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type hBaseProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewHBaseProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &hBaseProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (hp hBaseProvider) GetName() string {
	return hp.name
}

func (hp hBaseProvider) GetInterval() time.Duration {
	return hp.interval
}

func (hp hBaseProvider) GetDownThreshold() int64 {
	return hp.downThreshold
}

func (hp hBaseProvider) GetUpThreshold() int64 {
	return hp.upThreshold
}

func (hp hBaseProvider) Heartbeat() bool {
	client := gohbase.NewClient(hp.connectionString)
	defer client.Close()

	if client == nil {
		return false
	}

	return true
}
