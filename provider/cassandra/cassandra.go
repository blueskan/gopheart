package cassandra

import (
	"github.com/gocql/gocql"
	"strings"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type cassandraProvider struct {
	name             string
	connectionString string
	addresses        []string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewCassandraProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &cassandraProvider{
		name:             name,
		connectionString: connectionString,
		addresses:        strings.Split(connectionString, ","),
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (cp cassandraProvider) GetName() string {
	return cp.name
}

func (cp cassandraProvider) GetInterval() time.Duration {
	return cp.interval
}

func (cp cassandraProvider) GetDownThreshold() int64 {
	return cp.downThreshold
}

func (cp cassandraProvider) GetUpThreshold() int64 {
	return cp.upThreshold
}

func (cp cassandraProvider) Heartbeat() bool {
	cluster := gocql.NewCluster(cp.addresses...)
	cluster.Keyspace = "gopheart"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	defer session.Close()

	if err != nil {
		return false
	}

	return true
}
