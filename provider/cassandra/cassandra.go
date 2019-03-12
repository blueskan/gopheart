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
	downThreshold    int
	upThreshold      int
}

func NewCassandraProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int,
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

func (cp cassandraProvider) GetDownThreshold() int {
	return cp.downThreshold
}

func (cp cassandraProvider) GetUpThreshold() int {
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
