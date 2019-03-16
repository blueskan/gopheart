package mssql

import (
	"database/sql"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type mssqlProvider struct {
	name             string
	connectionString string
	timeout          time.Duration
	interval         time.Duration
	downThreshold    int64
	upThreshold      int64
}

func NewMssqlProvider(
	name, connectionString string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &mssqlProvider{
		name:             name,
		connectionString: connectionString,
		timeout:          timeout,
		interval:         interval,
		downThreshold:    downThreshold,
		upThreshold:      upThreshold,
	}
}

func (mp mssqlProvider) GetName() string {
	return mp.name
}

func (mp mssqlProvider) GetInterval() time.Duration {
	return mp.interval
}

func (mp mssqlProvider) GetDownThreshold() int64 {
	return mp.downThreshold
}

func (mp mssqlProvider) GetUpThreshold() int64 {
	return mp.upThreshold
}

func (mp mssqlProvider) Heartbeat() bool {
	db, err := sql.Open("mssql", mp.connectionString)
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
