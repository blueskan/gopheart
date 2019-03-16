package postgresql

import (
	"database/sql"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type postgresqlProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewPostgresqlProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &postgresqlProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (pp postgresqlProvider) GetName() string {
	return pp.name
}

func (pp postgresqlProvider) GetInterval() time.Duration {
	return pp.interval
}

func (pp postgresqlProvider) GetDownThreshold() int64 {
	return pp.downThreshold
}

func (pp postgresqlProvider) GetUpThreshold() int64 {
	return pp.upThreshold
}

func (mp postgresqlProvider) Heartbeat() bool {
	db, err := sql.Open("postgres", mp.connectionString)
	if err != nil {
		return false
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		return false
	}

	return true
}
