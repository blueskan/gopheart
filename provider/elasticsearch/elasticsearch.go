package elasticsearch

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type elasticsearchProvider struct {
	name          string
	url           string
	timeout       time.Duration
	interval      time.Duration
	downThreshold int64
	upThreshold   int64
}

func NewElasticsearchProvider(
	name, url string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &elasticsearchProvider{
		name:          name,
		url:           url,
		timeout:       timeout,
		interval:      interval,
		downThreshold: downThreshold,
		upThreshold:   upThreshold,
	}
}

func (es elasticsearchProvider) GetName() string {
	return es.name
}

func (es elasticsearchProvider) GetInterval() time.Duration {
	return es.interval
}

func (es elasticsearchProvider) GetDownThreshold() int64 {
	return es.downThreshold
}

func (es elasticsearchProvider) GetUpThreshold() int64 {
	return es.upThreshold
}

func (es elasticsearchProvider) Heartbeat() bool {
	timeout := time.Duration(es.timeout)

	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Get(es.url)
	if err != nil {
		return false
	}

	if res.StatusCode != 200 {
		return false
	}

	var body []byte
	_, err = res.Body.Read(body)

	if err != nil {
		return false
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)

	if err != nil {
		return false
	}

	status := bodyMap["status"].(string)

	if status != "green" {
		return false
	}

	return true
}
