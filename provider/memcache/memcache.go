package memcache

import (
	"strings"
	"time"

	"github.com/blueskan/gopheart/provider"
	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheProvider struct {
	name             string
	connectionString string
	addresses        []string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewMemcacheProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &memcacheProvider{
		name:             name,
		connectionString: connectionString,
		addresses:        strings.Split(connectionString, ","),
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (mp memcacheProvider) GetName() string {
	return mp.name
}

func (mp memcacheProvider) GetInterval() time.Duration {
	return mp.interval
}

func (mp memcacheProvider) GetDownThreshold() int64 {
	return mp.downThreshold
}

func (mp memcacheProvider) GetUpThreshold() int64 {
	return mp.upThreshold
}

func (mp memcacheProvider) Heartbeat() bool {
	mc := memcache.New(mp.addresses...)
	err := mc.Set(&memcache.Item{Key: "gopheart_health_check", Value: []byte("1")})

	if err != nil {
		return false
	}

	it, err := mc.Get("gopheart_health_check")

	if err != nil {
		return false
	}

	if string(it.Value) != "1" {
		return false
	}

	return true
}
