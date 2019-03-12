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
	downThreshold    int
	upThreshold      int
}

func NewHBaseProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int,
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

func (hp hBaseProvider) GetDownThreshold() int {
	return hp.downThreshold
}

func (hp hBaseProvider) GetUpThreshold() int {
	return hp.upThreshold
}

func (hp hBaseProvider) Heartbeat() bool {
	client := gohbase.NewClient(hp.connectionString)
	defer client.Close()

	// Check this thing
	if client == nil {
		return false
	}

	return true
}
