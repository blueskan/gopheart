package solr

import (
	"regexp"
	"strconv"
	"time"

	"github.com/blueskan/gopheart/provider"
	"github.com/rtt/Go-Solr"
)

type solrProvider struct {
	name          string
	host          string
	core          string
	port          string
	timeout       time.Duration
	interval      time.Duration
	downThreshold int64
	upThreshold   int64
}

func NewSolrProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	regex := regexp.MustCompile(`solr:\/\/(.*?)@(.*?):([0-9]*)`)
	res := regex.FindStringSubmatch(connectionString)

	host := res[1]
	core := res[2]
	port := res[3]

	return &solrProvider{
		name:          name,
		host:          host,
		core:          core,
		port:          port,
		timeout:       timeout,
		interval:      interval,
		downThreshold: downThreshold,
		upThreshold:   upThreshold,
	}
}

func (sp solrProvider) GetName() string {
	return sp.name
}

func (sp solrProvider) GetInterval() time.Duration {
	return sp.interval
}

func (sp solrProvider) GetDownThreshold() int64 {
	return sp.downThreshold
}

func (sp solrProvider) GetUpThreshold() int64 {
	return sp.upThreshold
}

func (sp solrProvider) Heartbeat() bool {
	port, err := strconv.Atoi(sp.port)

	if err != nil {
		return false
	}

	_, err = solr.Init(sp.host, port, sp.core)

	if err != nil {
		return false
	}

	return true
}
