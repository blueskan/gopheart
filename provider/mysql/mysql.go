package mysql

import (
	"database/sql"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type mysqlProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewMysqlProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &mysqlProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (mp mysqlProvider) GetName() string {
	return mp.name
}

func (mp mysqlProvider) GetInterval() time.Duration {
	return mp.interval
}

func (mp mysqlProvider) GetDownThreshold() int64 {
	return mp.downThreshold
}

func (mp mysqlProvider) GetUpThreshold() int64 {
	return mp.upThreshold
}

func (mp mysqlProvider) Heartbeat() bool {
	db, err := sql.Open("mysql", mp.connectionString)
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
