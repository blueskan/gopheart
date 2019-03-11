package url

import (
	"net/http"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type urlProvider struct {
	name          string
	url           string
	timeout       time.Duration
	interval      time.Duration
	downThreshold int
	upThreshold   int
}

func NewUrlProvider(
	name, url string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int,
) provider.Provider {
	return &urlProvider{
		name:          name,
		url:           url,
		timeout:       timeout,
		interval:      interval,
		downThreshold: downThreshold,
		upThreshold:   upThreshold,
	}
}

func (up urlProvider) GetName() string {
	return up.name
}

func (up urlProvider) GetInterval() time.Duration {
	return up.interval
}

func (up urlProvider) GetDownThreshold() int {
	return up.downThreshold
}

func (up urlProvider) GetUpThreshold() int {
	return up.upThreshold
}

// TODO provide other http verbs also extra headers and status code checks
func (up urlProvider) Heartbeat() bool {
	timeout := time.Duration(up.timeout)

	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(up.url)
	if err != nil {
		return false
	}

	if res.StatusCode != 200 {
		return false
	}

	return true
}
